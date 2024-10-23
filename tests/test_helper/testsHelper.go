package testhelper

import (
	"bytes"
	"fmt"
	"orth/cmd/core/embedded"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/embedded/optimizer"
	"orth/cmd/core/lexer"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"os"
	"os/exec"
	"regexp"
)

func ErrSliceToStringSlice(errs []error) []string {
	sErrors := make([]string, len(errs), cap(errs))
	for i, err := range errs {
		sErrors[i] = err.Error()
	}
	return sErrors
}

func WarningSliceToStringSlice(warns []orth_types.CompilerMessage) []string {
	sWarns := make([]string, len(warns))
	for i, warn := range warns {
		sWarns[i] = warn.Message
	}
	return sWarns
}

func PrepareComp(fileName string) ([]error, []orth_types.CompilerMessage) {
	strProgram := lexer.LoadProgramFromFile(fileName)
	lexedFiles := lexer.LexFile(strProgram)

	parsedOperations := make(chan orth_types.Pair[orth_types.Operation, error])

	program := orth_types.Program{
		Operations: make([]orth_types.Operation, 0),
		Warnings:   make([]orth_types.CompilerMessage, 0),
		Error:      make([]error, 0),
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
		return program.Error, program.Warnings
	}

	optimizedOperation, warnings := optimizer.AnalyzeAndOptimizeOperations(analyzerOperations)
	program.Warnings = append(program.Warnings, warnings...)
	program.Operations = append(program.Operations, optimizedOperation...)

	program, err := embedded.CrossReferenceBlocks(program)
	if err != nil {
		program.Error = append(program.Error, err)
	}

	if len(program.Error) != 0 {
		return program.Error, program.Warnings
	}

	embedded.Compile(program, *orth_debug.Compile)

	return program.Error, program.Warnings
}

func ExecOutput() (programOutput string) {
	execOutputExe := exec.Command(`.\output.exe`)
	var out bytes.Buffer
	execOutputExe.Stdout = &out

	execOutputExe.Run()
	programOutput = out.String()
	return
}

func ExecWithArgs(args ...string) (programOutput string) {
	execOutputExe := exec.Command(`.\output.exe`, args...)
	var out bytes.Buffer
	execOutputExe.Stdout = &out

	execOutputExe.Run()
	programOutput = out.String()
	return
}

func LoadExpected(fileName string) string {
	rgx := regexp.MustCompile(`\r`)
	expected, _ := os.ReadFile(fmt.Sprintf("./expected/%s.txt", fileName))

	return rgx.ReplaceAllString(string(expected), "")
}

func DumpOutput(out, fileName string) {
	dumpFile, _ := os.Create(fmt.Sprintf("./dumps/%s.txt", fileName))
	dumpFile.WriteString(out)
	dumpFile.Close()
}
