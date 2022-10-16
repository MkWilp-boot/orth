package main

import (
	testhelper "orth/tests/test_helper"
	"testing"
)

func TestArithmetics(t *testing.T) {
	testhelper.PrepareComp("./repo/arithmetics.orth")
	expected := testhelper.LoadExpected("TestArithmetics")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestArithmetics")
		t.FailNow()
	}
}
