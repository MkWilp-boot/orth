package main

import (
	testhelper "orth/tests/test_helper"
	"testing"
)

func TestDefineDirective(t *testing.T) {
	testhelper.PrepareComp("./repo/TestDefineDirective.orth")
	expected := testhelper.LoadExpected("TestDefineDirective")

	programOutput := testhelper.ExecOutput()

	if programOutput != expected {
		testhelper.DumpOutput(programOutput, "TestDefineDirective")
		t.FailNow()
	}
}
