package main

import (
	testhelper "orth/tests/test_helper"
	"testing"
)

func TestDefineDirective(t *testing.T) {
	testhelper.PrepareComp("./repo/rule110.orth")
	expected := testhelper.LoadExpected("TestRule110")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestRule110")
		t.FailNow()
	}
}
