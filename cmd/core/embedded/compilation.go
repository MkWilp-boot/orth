package embedded

import (
	"bufio"
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
	writer.WriteString("	fmt db \"%s\", 10, 0\n")
	writer.WriteString("	num_str db 32 dup (?)\n")

	// data segment (undefined)
	writer.WriteString(".DATA?\n")

	// code segment
	writer.WriteString(".CODE\n")
	writer.WriteString("start:\n")
	for i, operation := range program.Operations {
		switch operation.Instruction {
		case orthtypes.Push:
			writer.WriteString("	; push \n")
			writer.WriteString("	push " + operation.Operand.Operand + "\n")
		case orthtypes.Sum:
			writer.WriteString("	; sum \n")
			writer.WriteString("	pop eax\n")
			writer.WriteString("	pop ebx\n")
			writer.WriteString("	add eax, ebx\n")
			writer.WriteString("	push eax\n")
		case orthtypes.Minus:
			writer.WriteString("    ; minus \n")
			writer.WriteString("    pop eax\n")
			writer.WriteString("    pop ebx\n")
			writer.WriteString("    sub ebx, eax\n")
			writer.WriteString("    push ebx\n")
		case orthtypes.Equal:
			writer.WriteString("    ; equal \n")
			writer.WriteString("    pop eax\n")
			writer.WriteString("    pop ebx\n")
			writer.WriteString("    .if(eax == ebx)\n")
			writer.WriteString("    	push 1\n")
			writer.WriteString("    .else\n")
			writer.WriteString("    	push 0\n")
			writer.WriteString("    .endif\n")
		case orthtypes.Gt:
			writer.WriteString("    ; GT \n")
			writer.WriteString("    pop eax\n")
			writer.WriteString("    pop ebx\n")
			writer.WriteString("    .if(eax > ebx)\n")
			writer.WriteString("    	push 1\n")
			writer.WriteString("    .else\n")
			writer.WriteString("    	push 0\n")
			writer.WriteString("    .endif\n")
		case orthtypes.Lt:
			writer.WriteString("    ; LT \n")
			writer.WriteString("    pop eax\n")
			writer.WriteString("    pop ebx\n")
			writer.WriteString("    .if(eax < ebx)\n")
			writer.WriteString("    	push 1\n")
			writer.WriteString("    .else\n")
			writer.WriteString("    	push 0\n")
			writer.WriteString("    .endif\n")
		case orthtypes.If:
			writer.WriteString("	pop eax\n")
			writer.WriteString("	.if(eax)\n")
		case orthtypes.Else:
			writer.WriteString("	.else\n")
		case orthtypes.End:
			for _, v := range program.Operations[:i] {
				if v.RefBlock == operation.RefBlock-1 {
					if v.Instruction == orthtypes.Else || v.Instruction == orthtypes.If {
						writer.WriteString("	.endif\n")
					} else if v.Instruction == orthtypes.Do {
						writer.WriteString("	.endw\n")
					}
				}
			}
		case orthtypes.Print:
			writer.WriteString("    pop eax\n")
			writer.WriteString("	invoke crt__itoa, eax, OFFSET num_str, 10\n")
			writer.WriteString("	invoke crt_printf, OFFSET fmt, OFFSET num_str\n")
		}
	}
	writer.WriteString("; end program\n")
	writer.WriteString("	invoke ExitProcess, 0\n")
	writer.WriteString("end start\n")
	writer.Flush()
}
