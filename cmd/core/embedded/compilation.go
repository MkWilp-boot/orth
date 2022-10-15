package embedded

import (
	"bufio"
	"fmt"
	"log"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"os/exec"
)

var varsAndValues chan orthtypes.Pair[orthtypes.Operand, orthtypes.Operand]
var rmvFiles = []string{"mllink$.lnk", "output.ilk", "output.obj", "output.pdb"}

func init() {
	varsAndValues = make(chan orthtypes.Pair[orthtypes.Operand, orthtypes.Operand])
}

// Compile compiles a program into assembly
func Compile(program orthtypes.Program, assemblyType string) {
	log.Println("[INFO] Started compilation workflow")

	go embedded_helpers.RetrieveProgramInfo(program, varsAndValues, embedded_helpers.GetVarsAndValues)

	if assemblyType != "masm" {
		panic("[TEMP]: the current supported assembly is MASM")
	}
	output, err := os.Create("../../output.asm")
	if err != nil {
		panic(err)
	}

	compileMasm(program, output)

	compileCmd := exec.Command("ml64.exe", "../../output.asm", "/nologo", "/Zi", "/W3", "/link", "/entry:main")

	log.Println("[CMD] Running ML64")
	if err = compileCmd.Run(); err != nil {
		panic(err)
	}
	log.Println("[CMD] Finished running ML64")

	for _, file := range rmvFiles {
		log.Println("[CMD] Deleting extra files")
		if err = os.Remove(file); err != nil {
			panic(err)
		}
	}
}

func compileMasm(program orthtypes.Program, output *os.File) {
	log.Println("[CMD] Writing assembly")
	defer output.Close()

	// basic header stuff
	writer := bufio.NewWriter(output)
	writer.WriteString("include C:\\masm64\\include64\\masm64rt.inc\n")

	// data segment (pre-defined)
	writer.WriteString(".DATA\n")
	for pair := range varsAndValues {
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
		case orthtypes.Syscall5:
			writer.WriteString("; syscall with 3 params\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke GetStdHandle, rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	pop r9\n")
			writer.WriteString("	pop r10\n")
			writer.WriteString("	pop r11\n")
			writer.WriteString("	invoke WriteConsole, rax, rbx, r9, r10, r11\n")
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
			writer.WriteString("	lea rax, " + op.Operand.Operand + "\n")
			writer.WriteString("	push rax\n")
		case orthtypes.PutString:
			writer.WriteString("; Print string\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke StdOut, rax\n")
		}
	}
	writer.WriteString("end\n")
	writer.Flush()
	log.Println("[CMD] Finished writing assembly")
}
