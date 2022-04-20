package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"t/cmd/core/debug"
	"t/cmd/core/lexer"
	"t/cmd/core/embedded"
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
}

func main() {
	// args is
	// [0] -> file to be processed

	strProgram := lexer.LoadProgramFromFile(flag.Args()[0])
	program := embedded.CrossReferenceBlocks(embedded.ParseTokenAsOperation(strProgram))

	switch {
	case *debug.DumpVMCode:
		for _, v := range program.Operations {
			fmt.Printf("%#v\n", v)
		}
	case *debug.Simulate:
		embedded.Simulate(program)
	case *debug.CompileRun:
		panic("compile & run not implemented")
	case *debug.Compile:
		panic("compile not implemented")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
