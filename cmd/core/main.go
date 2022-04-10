package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"t/cmd/core/debug"
	"t/cmd/core/lexer"
	"t/cmd/core/simulation"
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
	strProgram := lexer.LoadProgramFromFile(flag.Args()[0])
	program := simulation.CrossReferenceBlocks(simulation.ParseTokenAsOperation(strProgram))

	if *debug.DumpVMCode {
		for _, v := range program.Operations {
			fmt.Printf("%#v\n", v)
		}
	} else {
		simulation.Run(program)
	}
}
