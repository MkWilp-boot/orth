package embedded

import (
	"bufio"
	"bytes"
	"fmt"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"os/exec"
)

// Compile compiles a program into assembly
func Compile(program orthtypes.Program, assemblyType string) {
	orth_debug.LogStep("[INFO] Started compilation workflow")

	outOfOrder := orthtypes.OutOfOrder{
		Vars: make(chan orthtypes.Pair[orthtypes.Operation, orthtypes.Operand]),
	}

	go embedded_helpers.RetrieveProgramInfo(program, outOfOrder, embedded_helpers.GetVarsAndValues)

	if assemblyType != "masm" {
		panic("[TEMP]: the current supported assembly is MASM")
	}

	finalAsm := fmt.Sprintf("%s.asm", *orth_debug.ObjectName)

	output, err := os.Create(finalAsm)
	if err != nil {
		panic(err)
	}

	compileMasm(program, outOfOrder, output)

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

func compileMasm(program orthtypes.Program, outOfOrder orthtypes.OutOfOrder, output *os.File) {
	orth_debug.LogStep("[CMD] Writing assembly")
	defer output.Close()

	// basic header stuff
	writer := bufio.NewWriter(output)
	writer.WriteString("include C:\\masm64\\include64\\masm64rt.inc\n")

	// data segment (pre-defined)
	writer.WriteString(".DATA\n")
	for pair := range outOfOrder.Vars {
		writer.WriteString("\t" + embedded_helpers.BuildVarDataSeg(pair) + "\n")
	}

	// data segment (undefined)
	writer.WriteString(".DATA?\n")
	writer.WriteString("	mem  BYTE 640000 dup(?)\n")
	writer.WriteString("	trash QWORD ?\n")

	// code segment
	writer.WriteString(".CODE\n")
	writer.WriteString("p_dump_ui64 PROC\n")
	writer.WriteString("	local buf[10]: BYTE\n")
	writer.WriteString("	push	rbx\n")
	writer.WriteString("	mov		rax, rcx\n")
	writer.WriteString("	lea		rcx, buf\n")
	writer.WriteString("    mov     rbx, 10\n")
	writer.WriteString("@@:\n")
	writer.WriteString("    xor     rdx, rdx\n")
	writer.WriteString("    div     rbx\n")
	writer.WriteString("    add     rdx, 48\n")
	writer.WriteString("    mov     BYTE PTR [rcx], dl\n")
	writer.WriteString("	dec		rcx\n")
	writer.WriteString("    test    rax, rax\n")
	writer.WriteString("    jnz     @b\n")
	writer.WriteString("	inc		rcx\n")
	writer.WriteString("	xor		rax, rax\n")
	writer.WriteString("    invoke	StdOut, rcx\n")
	writer.WriteString("	pop		rbx\n")
	writer.WriteString("    ret\n")
	writer.WriteString("p_dump_ui64 ENDP\n")

	writer.WriteString("; RCX: pointer pointing to where to start slicing\n")
	writer.WriteString("; RDX: amount of chars to slice\n")
	writer.WriteString("p_dump_mem proc\n")
	writer.WriteString("	local buffer[1024]: byte\n")
	writer.WriteString("	push rbx\n")
	writer.WriteString("	push rax\n")
	writer.WriteString("	push r8\n")
	writer.WriteString("	xor r8, r8\n")
	writer.WriteString("	lea rax, buffer\n")
	writer.WriteString("_begin:\n")
	writer.WriteString("	xor rbx, rbx\n")
	writer.WriteString("	mov bl, BYTE PTR [rcx+r8]\n")
	writer.WriteString("	mov [rax+r8], bl\n")
	writer.WriteString("	inc r8\n")
	writer.WriteString("	cmp rdx, r8\n")
	writer.WriteString("	jne _begin\n")
	writer.WriteString("_end:\n")
	writer.WriteString("	mov BYTE PTR [rax+r8], 0\n")
	writer.WriteString("	invoke StdOut, rax\n")
	writer.WriteString("	pop r8\n")
	writer.WriteString("	pop rax\n")
	writer.WriteString("	pop rbx\n")
	writer.WriteString("	ret\n")
	writer.WriteString("p_dump_mem endp\n")

	var immediateStringCount int
	immediateStrings := make([]orthtypes.Operand, 0)

	for ip := 0; ip < len(program.Operations); ip++ {
		writer.WriteString(fmt.Sprintf("addr_%d:\n", ip))

		op := program.Operations[ip]
		if op.Instruction == orthtypes.Skip {
			continue
		}
		// ignore vars so they are located on the data segment
		switch op.Instruction {
		case orthtypes.Push:
			writer.WriteString("; push\n")
			writer.WriteString("	push " + op.Operand.Operand + "\n")
		case orthtypes.PushStr:
			writer.WriteString("; push string\n")
			writer.WriteString("	mov rax, offset str_" + fmt.Sprint(immediateStringCount) + "\n")
			writer.WriteString("	push rax\n")
			immediateStrings = append(immediateStrings, op.Operand)
			immediateStringCount++
		case orthtypes.Mem:
			writer.WriteString("; push offset mem\n")
			writer.WriteString("	mov rax, offset mem\n")
			writer.WriteString("	push rax\n")
		case orthtypes.Load:
			writer.WriteString("; load\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	xor rbx, rbx\n")
			writer.WriteString("	mov bl, BYTE PTR [rax]\n")
			writer.WriteString("	push rbx\n")
		case orthtypes.Store:
			writer.WriteString("; store\n")
			writer.WriteString("	pop rbx ; value to store\n")
			writer.WriteString("	pop rax ; address of mem\n")
			writer.WriteString("	mov BYTE PTR [rax], bl\n")
			writer.WriteString("	xor rax, rax\n")
		case orthtypes.DumpMem:
			writer.WriteString("; dump_mem\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke p_dump_mem, rbx, rax\n")
		case orthtypes.Sum:
			writer.WriteString("; Sum\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	add rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orthtypes.Gt:
			writer.WriteString("; GT\n")
			writer.WriteString("	mov rdx, 1\n")
			writer.WriteString("	mov rcx, 0\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	cmovg rcx, rdx\n")
			writer.WriteString("	push rcx\n")
		case orthtypes.Lt:
			writer.WriteString("; LT\n")
			writer.WriteString("	mov rdx, 1\n")
			writer.WriteString("	mov rcx, 0\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	cmovl rcx, rdx\n")
			writer.WriteString("	push rcx\n")
		case orthtypes.Equal:
			writer.WriteString("; Equal\n")
			writer.WriteString("	mov rdx, 1\n")
			writer.WriteString("	mov rcx, 0\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	cmove rcx, rdx\n")
			writer.WriteString("	push rcx\n")
		case orthtypes.If:
			writer.WriteString("; If\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	test rax, rax\n")
			writer.WriteString(fmt.Sprintf("	jz addr_%d\n", op.RefBlock))
		case orthtypes.Else:
			writer.WriteString("; Else\n")
			writer.WriteString(fmt.Sprintf("	jmp addr_%d\n", op.RefBlock))
		case orthtypes.Proc:
			writer.WriteString("; Proc\n")
			writer.WriteString(op.Operand.Operand + " proc\n")
		case orthtypes.End:
			if program.Operations[op.RefBlock].Instruction == orthtypes.Proc {
				writer.WriteString("; Endp\n")
				writer.WriteString("	ret\n")
				writer.WriteString(program.Operations[op.RefBlock].Operand.Operand + " endp\n")
				continue
			}
			writer.WriteString("; End\n")
			writer.WriteString(fmt.Sprintf("	jmp addr_%d\n", op.RefBlock))
		case orthtypes.Invoke:
			writer.WriteString("; invoke\n")
			writer.WriteString("	invoke " + op.Operand.Operand + "\n")
		case orthtypes.Dup:
			writer.WriteString("; Dup\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rax\n")
		case orthtypes.TwoDup:
			writer.WriteString("; 2Dup\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	push rax\n")
		case orthtypes.Over:
			writer.WriteString("; Over\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rbx\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rbx\n")
		case orthtypes.While:
			writer.WriteString("; While\n")
		case orthtypes.Do:
			writer.WriteString("; Do\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	test rax, rax\n")
			writer.WriteString(fmt.Sprintf("	jz addr_%d\n", op.RefBlock))
		case orthtypes.Drop:
			writer.WriteString("; Drop\n")
			writer.WriteString("	pop trash\n")
		case orthtypes.PutU64:
			writer.WriteString("; DumpUI64\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke p_dump_ui64, rax\n")
		case orthtypes.Hold:
			writer.WriteString("; Hold var\n")
			writer.WriteString("	mov rax, offset " + embedded_helpers.MangleVarName(op) + "\n")
			writer.WriteString("	push rax\n")
		case orthtypes.PutString:
			writer.WriteString("; Print string\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke StdOut, rax\n")
		case orthtypes.Mult:
			writer.WriteString("; Mult\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	imul rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orthtypes.Minus:
			writer.WriteString("; Sub\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	sub rbx, rax\n")
			writer.WriteString("	push rbx\n")
		case orthtypes.Div:
			writer.WriteString("; Div\n")
			writer.WriteString("	xor rdx, rdx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	div rbx\n")
			writer.WriteString("	push rax\n")
		case orthtypes.Mod:
			writer.WriteString("; Mod\n")
			writer.WriteString("	xor rdx, rdx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	div rbx\n")
			writer.WriteString("	push rdx\n")
		case orthtypes.Swap:
			writer.WriteString("; Swap\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	push rax\n")
			writer.WriteString("	push rbx\n")
		case orthtypes.LShift:
			writer.WriteString("; shift left\n")
			writer.WriteString("	pop rcx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	shl rbx, cl\n")
			writer.WriteString("	push rbx\n")
		case orthtypes.RShift:
			writer.WriteString("; shift right\n")
			writer.WriteString("	pop rcx\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	shr rbx, cl\n")
			writer.WriteString("	push rbx\n")
		case orthtypes.LAnd:
			writer.WriteString("; bitwise and\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	and rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orthtypes.LOr:
			writer.WriteString("; bitwise or\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	or rax, rbx\n")
			writer.WriteString("	push rax\n")
		}
	}
	writer.WriteString(".DATA ;immediate strings\n")
	for i, s := range immediateStrings {
		writer.WriteString(fmt.Sprintf("	str_%d db %s \n", i, embedded_helpers.VarValueToAsmSyntax(s)))
	}
	writer.WriteString("end\n")
	writer.Flush()
	orth_debug.LogStep("[CMD] Finished writing assembly")
}
