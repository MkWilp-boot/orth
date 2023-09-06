package main

import (
	"flag"
	"fmt"
	"orth/cmd/core/embedded"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/embedded/optimizer"
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

	parsedOperations := make(chan orthtypes.Pair[orthtypes.Operation, error])

	program := orthtypes.Program{
		Operations: make([]orthtypes.Operation, 0),
	}

	go embedded.ParseTokenAsOperation(lexedFiles, parsedOperations)

	analyzerOperations := make([]orthtypes.Operation, 0)
	for parsedOperation := range parsedOperations {
		if parsedOperation.Right != nil {
			program.Error = append(program.Error, parsedOperation.Right)
			break
		}
		parsedOperation.Left = embedded_helpers.LinkVariableToValue(parsedOperation.Left, &analyzerOperations, &program)
		analyzerOperations = append(analyzerOperations, parsedOperation.Left)
	}

	if len(program.Error) != 0 {
		for _, err := range program.Error {
			fmt.Fprint(os.Stderr, err)
		}
		os.Exit(1)
	}

	optimizedOperation, warnings := optimizer.AnalyzeAndOptimizeOperations(analyzerOperations)
	program.Warnings = append(program.Warnings, warnings...)
	program.Operations = append(program.Operations, optimizedOperation...)

	for _, warning := range program.Warnings {
		fmt.Println(warning.Message)
	}

	program, err := embedded.CrossReferenceBlocks(program)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

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
