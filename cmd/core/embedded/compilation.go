package embedded

import (
	"bufio"
	"fmt"
	"os"
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
}

func compileMasm(program orthtypes.Program, output *os.File) {
	defer output.Close()

	// basic header stuff
	writer := bufio.NewWriter(output)
	writer.WriteString("option casemap :none\n")
	writer.WriteString("include C:\\masm32\\include\\masm32rt.inc\n")

	// data segment (pre-defined)
	writer.WriteString(".DATA\n")
	writer.WriteString("good db \"good\", 10, 0\n")
	writer.WriteString("bad db \"bad\", 10, 0\n")

	// data segment (undefined)
	writer.WriteString(".DATA?\n")
	writer.WriteString("trash dd ?\n")

	// code segment
	writer.WriteString(".CODE\n")
	writer.WriteString("start:\n")
	for ip := 0; ip < len(program.Operations); ip++ {
		writer.WriteString(fmt.Sprintf("addr_%d:\n", ip))
		op := program.Operations[ip]
		switch op.Instruction {
		case orthtypes.Push:
			writer.WriteString("; push\n")
			writer.WriteString("push " + op.Operand.Operand + "\n")
		case orthtypes.Sum:
			writer.WriteString("; Sum\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("pop ebx\n")
			writer.WriteString("add eax, ebx\n")
			writer.WriteString("push eax\n")
		case orthtypes.Gt:
			writer.WriteString("; GT\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("pop ebx\n")
			writer.WriteString("cmp eax, ebx\n")
			writer.WriteString(".if(eax > ebx)\n")
			writer.WriteString("push 1\n")
			writer.WriteString(".else\n")
			writer.WriteString("push 0\n")
			writer.WriteString(".endif\n")
		case orthtypes.Lt:
			writer.WriteString("; LT\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("pop ebx\n")
			writer.WriteString("cmp eax, ebx\n")
			writer.WriteString(".if(eax < ebx)\n")
			writer.WriteString("push 1\n")
			writer.WriteString(".else\n")
			writer.WriteString("push 0\n")
			writer.WriteString(".endif\n")
		case orthtypes.Equal:
			writer.WriteString("; Equal\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("pop ebx\n")
			writer.WriteString("cmp eax, ebx\n")
			writer.WriteString(".if(eax == ebx)\n")
			writer.WriteString("push 1\n")
			writer.WriteString(".else\n")
			writer.WriteString("push 0\n")
			writer.WriteString(".endif\n")
		case orthtypes.If:
			writer.WriteString("; If\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("test eax, eax\n")
			writer.WriteString(fmt.Sprintf("jz addr_%d\n", op.RefBlock))
		case orthtypes.Else:
			writer.WriteString("; Else\n")
			writer.WriteString(fmt.Sprintf("jmp addr_%d\n", op.RefBlock))
		case orthtypes.End:
			writer.WriteString("; End\n")
			writer.WriteString(fmt.Sprintf("jmp addr_%d\n", op.RefBlock))
		case orthtypes.Dup:
			writer.WriteString("; Dup\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("push eax\n")
			writer.WriteString("push eax\n")
		case orthtypes.While:
			writer.WriteString("; While\n")
		case orthtypes.Do:
			writer.WriteString("; Do\n")
			writer.WriteString("pop eax\n")
			writer.WriteString("test eax, eax\n")
			writer.WriteString(fmt.Sprintf("jz addr_%d\n", op.RefBlock))
		case orthtypes.Drop:
			writer.WriteString("; Drop\n")
			writer.WriteString("pop trash\n")
		}
	}
	writer.WriteString(fmt.Sprintf("addr_%d:\n", len(program.Operations)))
	writer.WriteString("; end program\n")
	writer.WriteString("invoke ExitProcess, 0\n")
	writer.WriteString("end start\n")
	writer.Flush()
}
