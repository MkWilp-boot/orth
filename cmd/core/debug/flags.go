package debug

import "flag"

var (
	Simulate   = flag.Bool("s", false, "Simulates a program")
	Compile    = flag.String("c", "", "-c[masm|nasm|fasm]")
	CompileRun = flag.Bool("r", false, "Runs a compiled program")
	DumpVMCode = flag.Bool("d", false, "Print to stdout the generated orth VM code")
	Help       = flag.Bool("h", false, "Describes useful thing about the compiler")
)

func init() {
	flag.Parse()
}
