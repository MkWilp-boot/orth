package embedded

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers"
	orth_types "orth/cmd/pkg/types"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const MASM_MAX_8BIT_CHAR_PER_LINE float64 = 20.0

// Compile compiles a program into assembly
func Compile(program orth_types.Program, assemblyType string) {
	orth_debug.LogStep("[INFO] Started compilation workflow")

	if assemblyType != "masm" {
		panic("[TEMP]: the current supported assembly is MASM")
	}

	finalAsm := fmt.Sprintf("%s.asm", *orth_debug.ObjectName)

	output, err := os.Create(finalAsm)
	if err != nil {
		panic(err)
	}

	compileMasm(program, output)

	if !*orth_debug.NoLink {
		compileCmd := exec.Command("ml64.exe", finalAsm, "/nologo", "/Zi", "/W3", "/link", "/entry:main")

		orth_debug.LogStep("[CMD] Running ML64")
		var stdout bytes.Buffer

		compileCmd.Stdout = &stdout

		if err = compileCmd.Run(); err != nil {
			fmt.Println(stdout.String())
			os.Exit(1)
		}
	}
	orth_debug.LogStep("[CMD] Finished running ML64")

	if *orth_debug.UnclearFiles || *orth_debug.NoLink {
		orth_debug.LogStep("[CMD] UCLR or NL flag active, files won't be deleted")
		return
	}
	embedded_helpers.CleanUp()
}

func compileMasm(program orth_types.Program, output *os.File) {
	orth_debug.LogStep("[CMD] Writing assembly")
	defer output.Close()

	// basic header stuff
	writer := bufio.NewWriter(output)
	writer.WriteString("include C:\\masm64\\include64\\masm64rt.inc\n")

	// data segment (pre-defined)
	writer.WriteString(".DATA\n")
	for i := 0; i < 32; i++ {
		writer.WriteString(fmt.Sprintf("	proc_arg_%d QWORD 0\n", i))
	}
	for i := 0; i < 32; i++ {
		writer.WriteString(fmt.Sprintf("	proc_ret_%d QWORD 0\n", i))
	}
	writer.Flush()

	writer.WriteString("\n.DATA ; MultScoped variables\n")
	for _, variable := range program.Variables {
		asmVar := embedded_helpers.BuildVarDataSeg(variable)
		writer.WriteString(fmt.Sprintf("	%s\n", asmVar))
	}
	writer.Flush()

	writer.WriteString("\n.DATA ; MultScoped constants\n")
	for _, variable := range program.Constants {
		asmVar := embedded_helpers.BuildVarDataSeg(variable)
		writer.WriteString(fmt.Sprintf("	%s\n", asmVar))
	}
	writer.WriteString("\n")

	writer.WriteString("	nArgc QWORD 0\n")
	writer.WriteString("	lError QWORD 0\n")

	// data segment (undefined)
	writer.WriteString(".DATA?\n")
	writer.WriteString("	mem  BYTE 640000 dup(?)\n")
	writer.WriteString("	trash QWORD ?\n")

	// code segment
	writer.WriteString(".CODE\n")

	writer.WriteString("; no return label\n")
	writer.WriteString("last_error_propagation:\n")
	writer.WriteString("	mrm lError, LastError$()\n")
	writer.WriteString("	invoke StdOut, lError\n")
	writer.WriteString("	invoke ExitProcess, 1\n")

	writer.WriteString("clear_proc_params PROC\n")
	for i := 0; i < 32; i++ {
		writer.WriteString(fmt.Sprintf("	mov proc_arg_%d, 0\n", i))
	}
	writer.WriteString("	ret\n")
	writer.WriteString("clear_proc_params ENDP\n")

	writer.WriteString("clear_proc_returns PROC\n")
	for i := 0; i < 32; i++ {
		writer.WriteString(fmt.Sprintf("	mov proc_ret_%d, 0\n", i))
	}
	writer.WriteString("	ret\n")
	writer.WriteString("clear_proc_returns ENDP\n")

	writer.WriteString("; RCX string buffer ptr\n")
	writer.WriteString("string_length proc\n")
	writer.WriteString("	mov rax, rcx\n")
	writer.WriteString(".L1:\n")
	writer.WriteString("	mov bl, BYTE PTR[rcx]\n")
	writer.WriteString("	cmp bl, 0\n")
	writer.WriteString("	jz .L2\n")
	writer.WriteString("	inc rcx\n")
	writer.WriteString("	jmp .L1\n")
	writer.WriteString(".L2:\n")
	writer.WriteString("	sub rcx, rax\n")
	writer.WriteString("	xchg rax, rcx\n")
	writer.WriteString("	inc rax\n")
	writer.WriteString("	ret\n")
	writer.WriteString("string_length endp\n")

	writer.WriteString("; RCX: pointer pointing to where to start slicing\n")
	writer.WriteString("; RDX: amount of chars to slice\n")
	writer.WriteString("p_dump_mem proc\n")
	writer.WriteString("	local buffer[1024]: byte\n")
	writer.WriteString("	push rbx\n")
	writer.WriteString("	push rax\n")
	writer.WriteString("	push r8\n")
	writer.WriteString("	xor r8, r8\n")
	writer.WriteString("	lea rax, buffer\n")
	writer.WriteString(".begin:\n")
	writer.WriteString("	xor rbx, rbx\n")
	writer.WriteString("	mov bl, BYTE PTR [rcx+r8]\n")
	writer.WriteString("	mov [rax+r8], bl\n")
	writer.WriteString("	inc r8\n")
	writer.WriteString("	cmp rdx, r8\n")
	writer.WriteString("	jne .begin\n")
	writer.WriteString(".end:\n")
	writer.WriteString("	mov BYTE PTR [rax+r8], 0\n")
	writer.WriteString("	invoke StdOut, rax\n")
	writer.WriteString("	pop r8\n")
	writer.WriteString("	pop rax\n")
	writer.WriteString("	pop rbx\n")
	writer.WriteString("	ret\n")
	writer.WriteString("p_dump_mem endp\n")
	writer.WriteString("put_char proc\n")
	writer.WriteString("	LOCAL hHandle   :QWORD\n")
	writer.WriteString("	LOCAL pChar     :QWORD\n")
	writer.WriteString("	LOCAL pBuff     :QWORD\n\n")
	writer.WriteString("	mov     pChar, rcx\n")
	writer.WriteString("	invoke  GetStdHandle, STD_OUTPUT_HANDLE\n")
	writer.WriteString("	cmp     rax, INVALID_HANDLE_VALUE\n")
	writer.WriteString("	je      last_error_propagation	; error handler defined on another file\n")
	writer.WriteString("	mov     hHandle, rax\n")
	writer.WriteString("	mov     pBuff, alloc(2)			; Allocate two bytes, one for the char and the null terminator.\n")
	writer.WriteString("	push	rsi\n")
	writer.WriteString("	mov     rdx, pBuff				; Load the address of pBuff into rdx.\n")
	writer.WriteString("	mov     rsi, pChar\n")
	writer.WriteString("	push	rax\n")
	writer.WriteString("	mov		al, [rsi]\n")
	writer.WriteString("	mov		[rdx], al\n")
	writer.WriteString("	pop		rax\n")
	writer.WriteString("	pop		rsi\n")
	writer.WriteString("	mov     BYTE PTR [rdx+1], 0  ; Null-terminate the buffer.\n")
	writer.WriteString("	invoke  WriteFile, hHandle, rdx, 1, 0, 0\n")
	writer.WriteString("	mfree   pBuff  ; Free the allocated memory.\n")
	writer.WriteString("	ret\n")
	writer.WriteString("put_char endp\n")
	writer.Flush()

	var immediateStringCount int
	immediateStrings := make(map[orth_types.Operand]int)

	lastProcMain := false

	for ip := 0; ip < len(program.Operations); ip++ {
		op := program.Operations[ip]
		if op.Instruction == orth_types.Skip {
			continue
		}
		switch op.Instruction {
		case orth_types.InstructionPush:
			writer.WriteString("; push\n")
			writer.WriteString("	push " + op.Operator.Operand + "\n")
		case orth_types.InstructionPushStr:
			strNum, ok := immediateStrings[op.Operator]
			if !ok {
				immediateStrings[op.Operator] = immediateStringCount
				strNum = immediateStringCount
			}
			writer.WriteString("; push string\n")
			writer.WriteString("	mov rax, offset str_" + fmt.Sprint(strNum) + "\n")
			writer.WriteString("	push rax\n")
			immediateStringCount++
		case orth_types.InstructionMem:
			writer.WriteString("; push offset mem\n")
			writer.WriteString("	mov rax, offset mem\n")
			writer.WriteString("	push rax\n")
		case orth_types.FunctionPutChar:
			writer.WriteString("; put_char\n")
			writer.WriteString("	pop rcx\n")
			writer.WriteString("	invoke put_char\n")
		case orth_types.FunctionAlloc:
			writer.WriteString("; alloc\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	mov rbx, alloc(rax)\n")
			writer.WriteString("	mov rax, rbx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rax\n")
		case orth_types.FunctionFree:
			writer.WriteString("; free\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	mfree rax\n")
		case orth_types.FunctionSetNumber:
			writer.WriteString("; set_number\n")
			writer.WriteString("	pop rax ; address\n")
			writer.WriteString("	pop rbx ; value\n")
			writer.WriteString("	mov [rax], rbx\n")
		case orth_types.FunctionSetString:
			writer.WriteString("; set_string\n")
			writer.WriteString("	pop rdi ; destination\n")
			writer.WriteString("	pop rsi ; source\n")
			writer.WriteString("	mov rcx, rsi ; source\n")
			writer.WriteString("	invoke string_length\n")
			writer.WriteString("	mov rcx, rax\n")
			writer.WriteString("	rep movsb\n")
		case orth_types.InstructionDeref:
			writer.WriteString("; deref\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	mov rbx, [rax]\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	mov rbx, rax\n")
		case orth_types.InstructionLoad:
			writer.WriteString("; load\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	xor rbx, rbx\n")
			writer.WriteString("	mov bl, BYTE PTR [rax]\n")
			writer.WriteString("	push rbx\n")
		case orth_types.InstructionStore:
			writer.WriteString("; store\n")
			writer.WriteString("	pop rbx ; value to store\n")
			writer.WriteString("	pop rax ; address of mem\n")
			writer.WriteString("	mov BYTE PTR [rax], bl\n")
			writer.WriteString("	xor rax, rax\n")
		case orth_types.FunctionDumpMem:
			writer.WriteString("; dump_mem\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke p_dump_mem, rbx, rax\n")
		case orth_types.InstructionSum:
			writer.WriteString("; Sum\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	add rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orth_types.InstructionGt:
			writer.WriteString("; GT\n")
			writer.WriteString("	mov rdx, 1\n")
			writer.WriteString("	mov rcx, 0\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	cmovg rcx, rdx\n")
			writer.WriteString("	push rcx\n")
		case orth_types.InstructionLt:
			writer.WriteString("; LT\n")
			writer.WriteString("	mov rdx, 1\n")
			writer.WriteString("	mov rcx, 0\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	cmovl rcx, rdx\n")
			writer.WriteString("	push rcx\n")
		case orth_types.InstructionEqual:
			writer.WriteString("; Equal\n")
			writer.WriteString("	mov rdx, 1\n")
			writer.WriteString("	mov rcx, 0\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	cmove rcx, rdx\n")
			writer.WriteString("	push rcx\n")
		case orth_types.InstructionIf:
			writer.WriteString("; If\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	test rax, rax\n")

			indexToJump := op.PrioritizeAddress()
			if indexToJump == -1 {
				log.Fatal("if indexToJump is -1")
			}

			writer.WriteString(fmt.Sprintf("	jz .L%d\n", indexToJump))
		case orth_types.InstructionElse:
			writer.WriteString(fmt.Sprintf(".L%d:\n", ip))
			writer.WriteString("; Else\n")
		case orth_types.InstructionProc:
			writer.WriteString("; Proc\n")
			writer.WriteString(op.Operator.Operand + " proc\n")

			procLocalVariables := make([]struct{ Initializer, Decl, Type string }, len(op.Context.Declarations))

			for ctxDeclI, ctxDecl := range op.Context.Declarations {
				scopeVariable := program.Operations[ctxDecl.Index]
				programDecl := scopeVariable.Links["variable_value"]

				varType := embedded_helpers.VarTypeToLocalAsmType(programDecl.Operator)
				varName := scopeVariable.Operator.Operand

				procLocalVariables[ctxDeclI] = struct {
					Initializer, Decl, Type string
				}{
					Type:        varType,
					Decl:        fmt.Sprintf("	LOCAL %s :%s\n", varName, varType),
					Initializer: fmt.Sprintf("	mov %s, %s\n", varName, programDecl.Operator.Operand),
				}
			}

			sort.Slice(procLocalVariables, func(i, j int) bool {
				left := embedded_helpers.AsmVariablePriority[procLocalVariables[i].Type]
				right := embedded_helpers.AsmVariablePriority[procLocalVariables[j].Type]
				return left > right
			})

			for _, variable := range procLocalVariables {
				writer.WriteString(variable.Decl)
			}
			for _, variable := range procLocalVariables {
				writer.WriteString(variable.Initializer)
			}

			lastProcMain = op.Operator.Operand == "main"
		case orth_types.InstructionWith:
			procParamsCount := 0
			for k := range op.Links {
				if !strings.HasPrefix(k, "proc_param_") {
					continue
				}
				procParamsCount++
			}

			writer.WriteString("; Params\n")

			if lastProcMain && procParamsCount > 0 {
				fmt.Println("[WARN] `with` instruction detected with more than 0 parameters for proc main, if you are trying to get command line arguments, proceed with `with cli` instead")
			}

			if lastProcMain && op.Operator.Operand == "cli" {
				writer.WriteString("; ArgC & ArgV\n")
				writer.WriteString("	invoke GetCommandLineW\n")
				writer.WriteString("	invoke CommandLineToArgvW, rax, addr nArgc\n")
				writer.WriteString("	push rax	; rax = pointer to argv\n")
				writer.WriteString("	mov  rax, nArgc\n")
				writer.WriteString("	push rax\n")
				writer.WriteString("	xor rax, rax\n")
			} else {
				for i := procParamsCount - 1; i >= 0; i-- {
					writer.WriteString(fmt.Sprintf("	push proc_arg_%d\n", i))
				}
			}
		case orth_types.InstructionEnd:
			writer.WriteString(fmt.Sprintf(".L%d:\n", ip))
			procAddress, procFound := op.Addresses[orth_types.InstructionProc]
			whileAddress, whileFound := op.Addresses[orth_types.InstructionWhile]

			closingMainProc := false

			for k, v := range op.Addresses {
				closingMainProc = !closingMainProc && (k == orth_types.InstructionProc && program.Operations[v].Operator.Operand == "main")
			}

			if procFound {
				writer.WriteString(fmt.Sprintf("; End for %s\n", orth_types.InstructionToStr(orth_types.InstructionProc)))

				outOpertaion := orth_types.Operation{}
				for _, operation := range program.Operations[procAddress:] {
					if operation.Instruction != orth_types.InstructionOut {
						continue
					}
					outOpertaion = operation
					break
				}

				outAmount := 0
				for k := range outOpertaion.Links {
					if !strings.HasPrefix(k, "proc_out_param_") {
						continue
					}
					outAmount++
				}

				if outAmount > 0 {
					for i := outAmount - 1; i >= 0; i-- {
						writer.WriteString(fmt.Sprintf("	pop proc_ret_%d\n", i))
					}
				}
				writer.WriteString("	invoke clear_proc_params\n")
				if closingMainProc {
					writer.WriteString("	invoke ExitProcess, 0\n")
				}
				writer.WriteString("	ret\n")
				writer.WriteString(fmt.Sprint(program.Operations[procAddress].Operator.Operand, " ", "endp\n"))
			} else if whileFound {
				writer.WriteString(fmt.Sprintf("; End for %s\n", orth_types.InstructionToStr(orth_types.InstructionWhile)))
				writer.WriteString(fmt.Sprintf("; Jump to %s\n", orth_types.InstructionToStr(orth_types.InstructionWhile)))
				writer.WriteString(fmt.Sprintf("	jmp .L%d\n", whileAddress))
				// post-instruction label
				writer.WriteString(fmt.Sprintf(".LA%d:\n", ip))
			}
		case orth_types.InstructionCall:
			writer.WriteString("; invoke\n")

			var callingProcedureArgumentsCount int
			var callingProcedureOutParamsCount int
			var callingProcedureIndex int
			for i, operation := range program.Operations {
				if operation.Operator.Operand == op.Operator.Operand && operation.Instruction == orth_types.InstructionProc {
					callingProcedureIndex = i
					break
				}
			}

			for _, operation := range program.Operations[callingProcedureIndex:] {
				if operation.Instruction == orth_types.InstructionWith {
					for k := range operation.Links {
						if !strings.HasPrefix(k, "proc_param_") {
							continue
						}
						callingProcedureArgumentsCount++
					}
				}
				if operation.Instruction == orth_types.InstructionOut {
					for k := range operation.Links {
						if !strings.HasPrefix(k, "proc_out_param_") {
							continue
						}
						callingProcedureOutParamsCount++
					}
				}
			}

			for i := 0; i < callingProcedureArgumentsCount; i++ {
				writer.WriteString(fmt.Sprintf("	pop proc_arg_%d\n", i))
			}
			writer.WriteString(fmt.Sprintf("	invoke %s\n", op.Operator.Operand))

			for i := 0; i < callingProcedureOutParamsCount; i++ {
				writer.WriteString(fmt.Sprintf("	push proc_ret_%d\n", i))
			}

			writer.WriteString("	invoke clear_proc_returns\n")
		case orth_types.InstructionDup:
			writer.WriteString("; Dup\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rax\n")
		case orth_types.InstructionTwoDup:
			writer.WriteString("; 2Dup\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	push rax\n")
		case orth_types.InstructionOver:
			writer.WriteString("; Over\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rbx\n")
		case orth_types.InstructionWhile:
			writer.WriteString(fmt.Sprintf(".L%d:\n", ip))
			writer.WriteString("; While\n")
		case orth_types.InstructionDo:
			writer.WriteString("; Do\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	test rax, rax\n")
			endAddress, ok := op.Addresses[orth_types.InstructionEnd]
			if !ok {
				log.Fatalln("do wihtout end")
			}
			writer.WriteString(fmt.Sprintf("	jz .LA%d\n", endAddress))
		case orth_types.InstructionDrop:
			writer.WriteString("; Drop\n")
			writer.WriteString("	pop trash\n")
		case orth_types.InstructionExit:
			writer.WriteString("; Exit program\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke ExitProcess, rax\n")
		case orth_types.FunctionPutU64:
			writer.WriteString("; DumpUI64\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	conout str$(rax)\n")
		case orth_types.InstructionHold:
			// priority for local variables, since Hold instruction can't point to more than one symbol
			if holdingVariable, ok := op.Links["hold_local"]; ok {
				writer.WriteString("; Hold local\n")
				writer.WriteString("	lea rax, " + holdingVariable.Operator.Operand + "\n")
				writer.WriteString("	push rax\n")
			} else {
				holdingVariable := op.Links["hold_mult"]
				writer.WriteString("; Hold MultScoped\n")
				writer.WriteString("	mov rax, offset " + embedded_helpers.MangleVarName(holdingVariable) + "\n")
				writer.WriteString("	push rax\n")
			}
		case orth_types.FunctionPutString:
			writer.WriteString("; Print string\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke StdOut, rax\n")
		case orth_types.InstructionMult:
			writer.WriteString("; Mult\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	imul rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orth_types.InstructionMinus:
			writer.WriteString("; Sub\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	sub rbx, rax\n")
			writer.WriteString("	push rbx\n")
		case orth_types.InstructionDiv:
			writer.WriteString("; Div\n")
			writer.WriteString("	xor rdx, rdx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	div rbx\n")
			writer.WriteString("	push rax\n")
		case orth_types.InstructionMod:
			writer.WriteString("; Mod\n")
			writer.WriteString("	xor rdx, rdx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	div rbx\n")
			writer.WriteString("	push rdx\n")
		case orth_types.InstructionSwap:
			writer.WriteString("; Swap\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rbx\n")
		case orth_types.InstructionLShift:
			writer.WriteString("; shift left\n")
			writer.WriteString("	pop rcx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	shl rbx, cl\n")
			writer.WriteString("	push rbx\n")
		case orth_types.InstructionRShift:
			writer.WriteString("; shift right\n")
			writer.WriteString("	pop rcx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	shr rbx, cl\n")
			writer.WriteString("	push rbx\n")
		case orth_types.InstructionLAnd:
			writer.WriteString("; bitwise and\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	and rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orth_types.InstructionLOr:
			writer.WriteString("; bitwise or\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	or rax, rbx\n")
			writer.WriteString("	push rax\n")
		}
		writer.Flush()
	}
	writer.WriteString(".DATA ; immediate strings\n")
	for v, i := range immediateStrings {
		length := float64(len(v.Operand))

		// checks if the string is larger than this weird masm exclusive constant
		if length > MASM_MAX_8BIT_CHAR_PER_LINE {
			// gets the amount of slices the string must have afte helpers.Chunks
			size := int(math.Ceil(length / MASM_MAX_8BIT_CHAR_PER_LINE))

			// chunk the string into slices of MASM_MAX_8BIT_CHAR_PER_LINE size
			chunks := helpers.Chunks(v.Operand, int(MASM_MAX_8BIT_CHAR_PER_LINE))

			// writes the string label definition
			writer.WriteString(fmt.Sprintf("	str_%d \\\n", i))
			for i, c := range chunks {
				var endWithNullByte bool

				// if it's the last element, must end in a '0' byte
				if i == size-1 {
					endWithNullByte = true
				}
				// writes the bytes
				writer.WriteString(fmt.Sprintf("\t\tdb %s\n", embedded_helpers.StringToByteRep(c, endWithNullByte)))
			}
			continue
		}
		writer.WriteString(fmt.Sprintf("	str_%d db %s \n", i, embedded_helpers.VarValueToAsmSyntax(v, true)))
	}
	writer.WriteString("end ; code segment\n")
	writer.Flush()
	orth_debug.LogStep("[CMD] Finished writing assembly")
}
