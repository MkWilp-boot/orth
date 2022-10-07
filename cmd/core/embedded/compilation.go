package embedded

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	orthtypes "t/cmd/pkg/types"
)

// Compile compiles a program into assembly
func Compile(program orthtypes.Program, assemblyType string) {
	if assemblyType != "masm" {
		panic("[TEMP]: the current supported assembly is MASM")
	}
	output, err := os.Create("../../output.asm")
	if err != nil {
		panic(err)
	}

	compileMasm(program, output)

	compileCmd := exec.Command("ml64.exe", "../../output.asm", "/nologo", "/Zi", "/W3", "/link", "/entry:main")

	if err = compileCmd.Run(); err != nil {
		panic(err)
	}
}

func compileMasm(program orthtypes.Program, output *os.File) {
	defer output.Close()

	// basic header stuff
	writer := bufio.NewWriter(output)
	writer.WriteString("include C:\\masm64\\include64\\masm64rt.inc\n")

	// data segment (pre-defined)
	writer.WriteString(".DATA\n")

	// data segment (undefined)
	writer.WriteString(".DATA?\n")
	writer.WriteString("trash dd ?\n")

	// code segment
	writer.WriteString(".CODE\n")
	writer.WriteString("dump_ui64 PROC\n")
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
	writer.WriteString("dump_ui64 ENDP\n")

	writer.WriteString("main PROC\n")
	for ip := 0; ip < len(program.Operations); ip++ {
		writer.WriteString(fmt.Sprintf("addr_%d:\n", ip))
		op := program.Operations[ip]
		switch op.Instruction {
		case orthtypes.Push:
			writer.WriteString("; push\n")
			writer.WriteString("	push " + op.Operand.Operand + "\n")
		case orthtypes.Sum:
			writer.WriteString("; Sum\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	add rax, rbx\n")
			writer.WriteString("	push rax\n")
		case orthtypes.Gt:
			writer.WriteString("; GT\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	.if(rax } rbx)\n")
			writer.WriteString("		push 1\n")
			writer.WriteString("	.else\n")
			writer.WriteString("		push 0\n")
			writer.WriteString("	.endif\n")
		case orthtypes.Lt:
			writer.WriteString("; LT\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	.if(rax { rbx)\n")
			writer.WriteString("		push 1\n")
			writer.WriteString("	.else\n")
			writer.WriteString("		push 0\n")
			writer.WriteString("	.endif\n")
		case orthtypes.Equal:
			writer.WriteString("; Equal\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	pop rbx\n")
			writer.WriteString("	cmp rax, rbx\n")
			writer.WriteString("	.if(rax == rbx)\n")
			writer.WriteString("		push 1\n")
			writer.WriteString("	.else\n")
			writer.WriteString("		push 0\n")
			writer.WriteString("	.endif\n")
		case orthtypes.If:
			writer.WriteString("; If\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	test rax, rax\n")
			writer.WriteString(fmt.Sprintf("	jz addr_%d\n", op.RefBlock))
		case orthtypes.Else:
			writer.WriteString("; Else\n")
			writer.WriteString(fmt.Sprintf("	jmp addr_%d\n", op.RefBlock))
		case orthtypes.End:
			writer.WriteString("; End\n")
			writer.WriteString(fmt.Sprintf("	jmp addr_%d\n", op.RefBlock))
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
		case orthtypes.DumpUI64:
			writer.WriteString("; DumpUI64\n")
			writer.WriteString("	pop rax\n")
			writer.WriteString("	invoke dump_ui64, rax\n")
		}
	}
	writer.WriteString(fmt.Sprintf("addr_%d:\n", len(program.Operations)))
	writer.WriteString("	ret\n")
	writer.WriteString("main ENDP\n")
	writer.WriteString("end\n")
	writer.Flush()
}
