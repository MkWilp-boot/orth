package debug

import "flag"

var (
	Simulate   = flag.Bool("s", false, "Simulates a program")
	Compile    = flag.String("c", "", "-c[masm|nasm|fasm]")
	CompileRun = flag.Bool("r", false, "Runs a compiled program")
	DumpVMCode = flag.Bool("d", false, "Print to stdout the generated orth VM code")
	Help       = flag.Bool("h", false, "Describes useful thing about the compiler")
)

func ToStringIntruction() map[int]string {
	m := make(map[int]string)
	m[0] = "Push"
	m[1] = "Sum"
	m[2] = "Minus"
	m[3] = "Mult"
	m[4] = "Div"
	m[5] = "If"
	m[6] = "Else"
	m[7] = "End"
	m[8] = "Equal"
	m[8] = "Lt"
	m[10] = "Gt"
	m[11] = "NotEqual"
	m[12] = "Dup"
	m[13] = "Dump"
	m[14] = "Print"
	m[15] = "Do"
	m[16] = "Drop"
	m[17] = "While"
	m[18] = "Swap"
	m[19] = "Mod"
	m[20] = "Mem"
	m[21] = "Store"
	m[22] = "Load"
	m[23] = "LoadStay"
	m[24] = "Func"
	m[25] = "Call"
	return m
}

func init() {
	flag.Parse()
}
