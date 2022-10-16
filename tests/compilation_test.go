package main

import (
	testhelper "orth/tests/test_helper"
	"testing"
)

func TestIntegerArithmetics(t *testing.T) {
	testhelper.PrepareComp("./repo/integer_arithmetics.orth")
	expected := testhelper.LoadExpected("TestIntegerArithmetics")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestIntegerArithmetics")
		t.FailNow()
	}
}
