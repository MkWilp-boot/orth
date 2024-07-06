package orthtypes

import (
	"fmt"
	"strings"
)

const (
	MAX_PROC_PARAM_COUNT  = 32
	MAX_PROC_OUTPUT_COUNT = 32
)

type Instruction uint16

const (
	Push Instruction = iota + 1
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
	SetNumber
	SetString
	Alloc
	Free
	PutChar
	TotalOps
)

var instructionNames map[Instruction]string

func init() {
	instructionNames = map[Instruction]string{
		Push:      "Push",
		PushStr:   "PushStr",
		Sum:       "Sum",
		Minus:     "Minus",
		Mult:      "Mult",
		Div:       "Div",
		If:        "If",
		Else:      "Else",
		End:       "End",
		Equal:     "Equal",
		Lt:        "Lt",
		Gt:        "Gt",
		NotEqual:  "NotEqual",
		Dup:       "Dup",
		TwoDup:    "TwoDup",
		PutU64:    "PutU64",
		PutString: "PutString",
		Do:        "Do",
		Drop:      "Drop",
		While:     "While",
		Swap:      "Swap",
		Mod:       "Mod",
		Mem:       "Mem",
		Store:     "Store",
		Load:      "Load",
		LoadStay:  "LoadStay",
		Func:      "Func",
		Call:      "Call",
		OType:     "OType",
		Const:     "Const",
		Var:       "Var",
		Gvar:      "Gvar",
		Hold:      "Hold",
		Skip:      "Skip",
		Nop:       "Nop",
		Proc:      "Proc",
		In:        "In",
		Invoke:    "Invoke",
		DumpMem:   "DumpMem",
		LShift:    "LShift",
		RShift:    "RShift",
		LAnd:      "LAnd",
		LOr:       "LOr",
		Over:      "Over",
		Exit:      "Exit",
		With:      "With",
		Out:       "Out",
		Deref:     "Deref",
		SetNumber: "SetNumber",
		SetString: "SetString",
		Alloc:     "Alloc",
		Free:      "Free",
		PutChar:   "PutChar",
	}

	if len(instructionNames) != int(TotalOps)-1 {
		panic("[DEV] Missing instruction on name map")
	}
}

func PPrintOperation(op Operation) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", InstructionToStr(op.Instruction)))
	builder.WriteString(fmt.Sprintf("	operand: %s | symbolName%q\n", op.Operator.Operand, op.Operator.SymbolName))
	for k, v := range op.Links {
		builder.WriteString(fmt.Sprintf("	link_name: %q | link_type: %q | link_value: %q\n", k, v.Operator.SymbolName, v.Operator.Operand))
	}
	for k, v := range op.Addresses {
		builder.WriteString(fmt.Sprintf("\n** %s: %d", InstructionToStr(k), v))
	}
	builder.WriteString("****************************************************\n")
	return builder.String()
}

func InstructionToStr(inst Instruction) string {
	if inst >= TotalOps {
		return ""
	}
	return instructionNames[inst]
}

// Program is the main struct for a transpiled
// orth code into machine code
type Program struct {
	Warnings   []CompilerMessage
	Error      []error
	Variables  []Operation
	Constants  []Operation
	Operations []Operation
}

func (p *Program) Filter(predicate func(op Operation, i int) bool) []Pair[int, Operation] {
	ops := make([]Pair[int, Operation], 0)
	for i, op := range p.Operations {
		if predicate(op, i) {
			ops = append(ops, Pair[int, Operation]{
				Left:  i,
				Right: op,
			})
		}
	}
	return ops
}

type WarnDegree uint8

const (
	Minor WarnDegree = iota
	Commom
	High
)

type CompilerMessage struct {
	Type    WarnDegree
	Message string
}
