package main

import (
	testhelper "orth/tests/test_helper"
	"strings"
	"testing"
)

func TestRule110(t *testing.T) {
	testhelper.PrepareComp("./repo/rule110.orth")
	expected := testhelper.LoadExpected("TestRule110")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestRule110")
		t.FailNow()
	}
}

func TestBigStrings(t *testing.T) {
	testhelper.PrepareComp("./repo/big_strings.orth")
	expected := testhelper.LoadExpected("TestBigStrings")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestBigStrings")
		t.FailNow()
	}
}

func TestBitWise(t *testing.T) {
	testhelper.PrepareComp("./repo/bitwise.orth")
	expected := testhelper.LoadExpected("TestBitWise")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestBitWise")
		t.FailNow()
	}
}

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

func TestLoops(t *testing.T) {
	testhelper.PrepareComp("./repo/loops.orth")
	expected := testhelper.LoadExpected("TestLoops")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestLoops")
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

func TestProc(t *testing.T) {
	testhelper.PrepareComp("./repo/procs.orth")
	expected := testhelper.LoadExpected("Procs")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "Procs")
		t.FailNow()
	}
}

func TestInvalidMemUsage(t *testing.T) {
	errors := testhelper.PrepareComp("./repo/mem_invalid_usage.orth")
	expected := testhelper.LoadExpected("TestInvalidMemUsage")

	programErros := strings.Join(testhelper.ErrSliceToStringSlice(errors), "\n")

	if programErros != expected {
		testhelper.DumpOutput(programErros, "TestInvalidMemUsage")
		t.FailNow()
	}
}
