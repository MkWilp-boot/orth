package testhelper

import (
	"bytes"
	"fmt"
	"orth/cmd/core/embedded"
	"orth/cmd/core/lexer"
	"orth/cmd/core/orth_debug"
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
	strProgram := lexer.LoadProgramFromFile(fileName)
	lexedFiles := lexer.LexFile(strProgram)
	program, tokenErrors := embedded.ParseTokenAsOperation(lexedFiles)

	if len(tokenErrors) == 0 {
		program = embedded.CrossReferenceBlocks(program)
		embedded.Compile(program, *orth_debug.Compile)
	}

	return tokenErrors
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
