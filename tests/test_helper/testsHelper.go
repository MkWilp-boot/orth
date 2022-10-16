package testhelper

import (
	"bytes"
	"fmt"
	"orth/cmd/core/embedded"
	"orth/cmd/core/lexer"
	"orth/cmd/core/orth_debug"
	"os"
	"os/exec"
)

func PrepareComp(fileName string) {
	program := embedded.CrossReferenceBlocks(
		embedded.ParseTokenAsOperation(
			lexer.LexFile(
				lexer.LoadProgramFromFile(fileName))))

	embedded.Compile(program, *orth_debug.Compile)
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
	expected, _ := os.ReadFile(fmt.Sprintf("./expected/%s.txt", fileName))
	return string(expected)
}

func DumpOutput(out, fileName string) {
	dumpFile, _ := os.Create(fmt.Sprintf("./dumps/%s.txt", fileName))
	dumpFile.WriteString(out)
	dumpFile.Close()
}
