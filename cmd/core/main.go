package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"t/cmd/core/debug"
	"t/cmd/core/embedded"
	"t/cmd/core/lexer"
	"t/cmd/pkg/helpers/functions"
	orthtypes "t/cmd/pkg/types"
)

func init() {
	if *debug.Help {
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
	if !*debug.Help && !*debug.Simulate && (*debug.Compile == "") && !*debug.CompileRun && !*debug.DumpVMCode {
		fmt.Println("Error, must select a run option.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	// args is
	// [0] -> file to be processed

	strProgram := lexer.LoadProgramFromFile(flag.Args()[0])
	program := embedded.CrossReferenceBlocks(embedded.ParseTokenAsOperation(strProgram))

	switch {
	case *debug.DumpVMCode:
		mapped := debug.ToStringIntruction()
		for _, v := range program.Operations {
			fmt.Printf("action %q\t type %q\t operand %q\t refblock %d\n", mapped[v.Instruction], v.Operand.VarType, v.Operand.Operand, v.RefBlock)
		}
	case *debug.Simulate:
		embedded.Simulate(program)
	case *debug.CompileRun:
		panic("compile & run not implemented")
	case *debug.Compile != "":
		fmt.Println("[WARN] Be aware that compilation may not work as expected due the fact of me being lazy :)")
		embedded.Compile(program, functions.CheckAsmType(*debug.Compile))
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
