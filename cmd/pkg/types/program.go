package orthtypes

// Program is the main struct for a transpiled
// orth code into machine code
type Program struct {
	Id         int
	Operations []Operation
}

func (p *Program) Filter(predicate func(op Operation, i int) bool) []Pair[int, Operation] {
	ops := make([]Pair[int, Operation], 0)
	for i, op := range p.Operations {
		if predicate(op, i) {
			ops = append(ops, Pair[int, Operation]{
				VarName:  i,
				VarValue: op,
			})
		}
	}
	return ops
}

const (
	MAX_PROC_PARAM_COUNT  = 32
	MAX_PROC_OUTPUT_COUNT = 32
)

const (
	Push int = iota
	PushStr
	Sum
	Minus
	Mult
	Div
	If
	Else
	End
	Equal
	Lt
	Gt
	NotEqual
	Dup
	TwoDup
	PutU64
	PutString
	Do
	Drop
	While
	Swap
	Mod
	Mem
	Store
	Load
	LoadStay
	Func
	Call
	OType
	Const
	Var
	Gvar
	Hold
	Skip
	Nop
	Proc
	In
	Invoke
	DumpMem
	LShift
	RShift
	LAnd
	LOr
	Over
	Exit
	With
	Out
	Deref
	TotalOps
)
