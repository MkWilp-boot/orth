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
	"orth/cmd/pkg/simulation"
	orth_types "orth/cmd/pkg/types"
	"os"
	"strings"
)

func init() {
	flag.Parse()
	sourceCodePath := flag.Args()[0]

	if *orth_debug.Help {
		flag.PrintDefaults()
	}
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: <orth> <file_path>")
		os.Exit(1)
	}
	if !strings.HasSuffix(sourceCodePath, orth_types.FileType) {
		fmt.Printf("[ERROR] The selected file %q is not of type %q\n", sourceCodePath, orth_types.FileType)
		os.Exit(1)
	}
	if !*orth_debug.Help && (*orth_debug.Compile == "") {
		fmt.Println("Error, must select a run option.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	sourceCodePath := flag.Args()[0]
	strProgram := lexer.LoadProgramFromFile(sourceCodePath)
	lexedFiles := lexer.LexFile(strProgram)

	parsedOperations := make(chan orth_types.Pair[orth_types.Operation, error])

	program := orth_types.Program{
		Operations: make([]orth_types.Operation, 0),
	}

	go embedded.ParseTokenAsOperation(lexedFiles, parsedOperations)

	analyzerOperations := make([]orth_types.Operation, 0)
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
	if *orth_debug.Sim {
		simulation.SimulateStack(&program)
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *orth_debug.Compile != "":
		orth_debug.LogStep(fmt.Sprintf("[INFO] Compilation started. Selected assembly is %q", *orth_debug.Compile))
		asmTarget, err := functions.CheckAsmType(*orth_debug.Compile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		embedded.Compile(program, asmTarget)
		orth_debug.LogStep("[INFO] Finished compilation.")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
