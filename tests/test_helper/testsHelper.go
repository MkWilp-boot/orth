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

func PrepareComp(fileName string) []error {
	errs := make([]error, 0)
	strProgram := lexer.LoadProgramFromFile(fileName)
	lexedFiles := lexer.LexFile(strProgram)

	parseTokenResult := make(chan orthtypes.Pair[orthtypes.Program, error])
	go embedded.ParseTokenAsOperation(lexedFiles, parseTokenResult)

	result := <-parseTokenResult

	if result.VarValue != nil {
		errs = append(errs, result.VarValue)
		return errs
	}

	crossRefererenceResult := make(chan orthtypes.Pair[orthtypes.Program, error])
	go embedded.CrossReferenceBlocks(result.VarName, crossRefererenceResult)

	result = <-crossRefererenceResult

	if result.VarValue != nil {
		errs = append(errs, result.VarValue)
		return errs
	}
	program := result.VarName
	embedded.Compile(program, *orth_debug.Compile)

	return errs
}

func ExecOutput() (programOutput string) {
	execOutputExe := exec.Command(`.\output.exe`)
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
