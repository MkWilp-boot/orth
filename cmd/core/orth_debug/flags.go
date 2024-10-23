package orth_debug

import (
	"flag"
	"log"
)

var (
	ObjectName   = flag.String("o", "output", "-o=final_executable.exe")
	Compile      = flag.String("com", "", "-com[masm|nasm|fasm]")
	Help         = flag.Bool("help", false, "Describes useful thing about the compiler")
	Log          = flag.Bool("log", false, "Enable log for each step")
	NoLink       = flag.Bool("nl", false, "Generates the assembly whitout linking")
	UnclearFiles = flag.Bool("uclr", false, "do not remove the generated output files")
	I            = flag.String("I", "", "appends paths for includes separeted by ','")
	Sim          = flag.Bool("sim", false, "simulate program's stack")
)

func LogStep(message string) {
	if !*Log {
		return
	}

	log.Println(message)
}
