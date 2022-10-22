package main

import (
	testhelper "orth/tests/test_helper"
	"strings"
	"testing"
)

func TestCompilationErrorMessages(t *testing.T) {
	errors := testhelper.PrepareComp("./repo/compilation_error.orth")
	expected := testhelper.LoadExpected("TestCompilationErrorMessages")

	programErros := strings.Join(testhelper.ErrSliceToStringSlice(errors), "\n")
	if programErros != expected {
		testhelper.DumpOutput(programErros, "TestCompilationErrorMessages")
		t.FailNow()
	}
}

func TestVarMangle(t *testing.T) {
	errors := testhelper.PrepareComp("./repo/var_mangle.orth")
	expected := testhelper.LoadExpected("TestVarMangle")

	programErros := strings.Join(testhelper.ErrSliceToStringSlice(errors), "\n")

	if programErros != expected {
		testhelper.DumpOutput(programErros, "TestVarMangle")
		t.FailNow()
	}
}

func TestIntegerArithmetics(t *testing.T) {
	testhelper.PrepareComp("./repo/integer_arithmetics.orth")
	expected := testhelper.LoadExpected("TestIntegerArithmetics")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestIntegerArithmetics")
		t.FailNow()
	}
}

func TestMem(t *testing.T) {
	testhelper.PrepareComp("./repo/simple_mem.orth")
	expected := testhelper.LoadExpected("TestMem")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestMem")
		t.FailNow()
	}
}
