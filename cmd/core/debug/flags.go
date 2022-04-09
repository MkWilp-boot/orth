package debug

import "flag"

var (
	DumpVMCode = flag.Bool("d", false, "Prints to stdout the generated orth VM code")
	Help       = flag.Bool("h", false, "Describes useful thing about the compiler")
)

func init() {
	flag.Parse()
}
