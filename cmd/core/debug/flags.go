package debug

import "flag"

var (
	Simulate   = flag.Bool("s", false, "Simulates a program")
	Compile    = flag.Bool("c", false, "Compiles a program")
	CompileRun = flag.Bool("cr", false, "Compile and runs a program")
	DumpVMCode = flag.Bool("d", false, "Print to stdout the generated orth VM code")
	Help       = flag.Bool("h", false, "Describes useful thing about the compiler")
)

func init() {
	flag.Parse()
}
