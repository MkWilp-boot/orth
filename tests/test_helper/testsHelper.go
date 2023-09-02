package testhelper

import (
	"bytes"
	"fmt"
	"orth/cmd/core/embedded"
	"orth/cmd/core/lexer"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
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

func WarningSliceToStringSlice(warns []orthtypes.CompilerMessage) []string {
	sWarns := make([]string, len(warns))
	for i, warn := range warns {
		sWarns[i] = warn.Message
	}
	return sWarns
}

func PrepareComp(fileName string) ([]error, []orthtypes.CompilerMessage) {
	strProgram := lexer.LoadProgramFromFile(fileName)
	lexedFiles := lexer.LexFile(strProgram)

	parsedOperations := make(chan orthtypes.Pair[orthtypes.Operation, error])

	program := orthtypes.Program{
		Operations: make([]orthtypes.Operation, 0),
		Warnings:   make([]orthtypes.CompilerMessage, 0),
		Error:      make([]error, 0),
	}

	go embedded.ParseTokenAsOperation(lexedFiles, parsedOperations)

	analyzerOperations := make([]orthtypes.Operation, 0)
	for parsedOperation := range parsedOperations {
		if parsedOperation.Right != nil {
			program.Error = append(program.Error, parsedOperation.Right)
			break
		}
		analyzerOperations = append(analyzerOperations, parsedOperation.Left)
	}

	if len(program.Error) != 0 {
		return program.Error, program.Warnings
	}

	optimizedOperation, warnings := embedded.AnalyzeAndOptimizeOperations(analyzerOperations)
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
