package orthtypes

// Program is the main struct for a transpiled
// orth code into machine code
type Program struct {
	Id         uint
	Operations []Operation
}

const (
	Push int = iota
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
	Dump
	Print
	Do
	Drop
	While
	Swap
	Mod
	Mem
	Store
	Load
	TotalOps
)
