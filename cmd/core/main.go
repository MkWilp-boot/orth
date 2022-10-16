package main

import (
	"flag"
	"fmt"
	"orth/cmd/core/embedded"
	"orth/cmd/core/lexer"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers/functions"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"strings"
)

func init() {
	flag.Parse()

	if *orth_debug.Help {
		flag.PrintDefaults()
	}
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: <orth> <file_path>")
		os.Exit(1)
	}
	if !strings.HasSuffix(flag.Args()[0], orthtypes.FileType) {
		fmt.Println("=================================================================================================")
		fmt.Printf("WARNING! The following file %q is not of type %q, the content may not be well formatted\n", flag.Args()[0], orthtypes.FileType)
		fmt.Println("=================================================================================================")
	}
	if !*orth_debug.Help && !*orth_debug.Simulate && (*orth_debug.Compile == "") && !*orth_debug.DumpVMCode {
		fmt.Println("Error, must select a run option.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	strProgram := lexer.LoadProgramFromFile(flag.Args()[0])
	tokenProgram := lexer.LexFile(strProgram)
	program := embedded.CrossReferenceBlocks(embedded.ParseTokenAsOperation(tokenProgram))

	switch {
	case *orth_debug.DumpVMCode:
		mapped := orth_debug.ToStringIntruction()
		for _, v := range program.Operations {
			fmt.Printf("action %q\t type %q\t operand %q\t refblock %d\n", mapped[v.Instruction], v.Operand.VarType, v.Operand.Operand, v.RefBlock)
		}
	case *orth_debug.Simulate:
		embedded.Simulate(program)
	case *orth_debug.Compile != "":
		orth_debug.LogStep(fmt.Sprintf("[INFO] Compilation started. Selected assembly is %q.\n", *orth_debug.Compile))
		embedded.Compile(program, functions.CheckAsmType(*orth_debug.Compile))
		orth_debug.LogStep("[INFO] Finished compilation.")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
