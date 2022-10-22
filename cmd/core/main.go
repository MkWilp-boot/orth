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
	if !*orth_debug.Help && (*orth_debug.Compile == "") {
		fmt.Println("Error, must select a run option.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	strProgram := lexer.LoadProgramFromFile(flag.Args()[0])
	lexedFiles := lexer.LexFile(strProgram)
	program, tokenErrors := embedded.ParseTokenAsOperation(lexedFiles)
	if tokenErrors != nil {
		for _, err := range tokenErrors {
			fmt.Printf("error: %v\n", err)
		}
		return
	}

	program = embedded.CrossReferenceBlocks(program)

	switch {
	case *orth_debug.Compile != "":
		orth_debug.LogStep(fmt.Sprintf("[INFO] Compilation started. Selected assembly is %q", *orth_debug.Compile))
		embedded.Compile(program, functions.CheckAsmType(*orth_debug.Compile))
		orth_debug.LogStep("[INFO] Finished compilation.")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
