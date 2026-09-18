package orth_types

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	MAX_PROC_PARAM_COUNT  = 32
	MAX_PROC_OUTPUT_COUNT = 32
)

type Instruction uint16

const (
	InstructionInvalid Instruction = iota
	InstructionPush
	InstructionPushStr
	InstructionSum
	InstructionMinus
	InstructionMult
	InstructionDiv
	InstructionIf
	InstructionElse
	InstructionEnd
	InstructionEqual
	InstructionLt
	InstructionGt
	InstructionNotEqual
	InstructionDup
	InstructionTwoDup
	FunctionPutU64
	FunctionPutString
	InstructionDo
	InstructionDrop
	InstructionWhile
	InstructionSwap
	InstructionMod
	InstructionMem
	InstructionStore
	InstructionLoad
	InstructionLoadStay
	InstructionFunc
	InstructionCall
	InstructionType
	InstructionConst
	InstructionVar
	InstructionGvar
	InstructionHold
	InstructionNop
	InstructionProc
	InstructionParam
	InstructionReturnType
	InstructionIn
	InstructionInvoke
	FunctionDumpMem
	InstructionLShift
	InstructionRShift
	InstructionLAnd
	InstructionLOr
	InstructionOver
	InstructionExit
	InstructionDeref
	FunctionSetNumber
	FunctionSetString
	FunctionAlloc
	FunctionFree
	FunctionPutChar
	Skip
	TotalOps
)

var instructionNames map[Instruction]string

func init() {
	instructionNames = map[Instruction]string{
		InstructionPush:       "InstructionPush",
		InstructionPushStr:    "InstructionPushStr",
		InstructionSum:        "InstructionSum",
		InstructionMinus:      "InstructionMinus",
		InstructionMult:       "InstructionMult",
		InstructionDiv:        "InstructionDiv",
		InstructionIf:         "InstructionIf",
		InstructionElse:       "InstructionElse",
		InstructionEnd:        "InstructionEnd",
		InstructionEqual:      "InstructionEqual",
		InstructionLt:         "InstructionLt",
		InstructionGt:         "InstructionGt",
		InstructionNotEqual:   "InstructionNotEqual",
		InstructionDup:        "InstructionDup",
		InstructionTwoDup:     "InstructionTwoDup",
		InstructionDo:         "InstructionDo",
		InstructionDrop:       "InstructionDrop",
		InstructionWhile:      "InstructionWhile",
		InstructionSwap:       "InstructionSwap",
		InstructionMod:        "InstructionMod",
		InstructionMem:        "InstructionMem",
		InstructionStore:      "InstructionStore",
		InstructionLoad:       "InstructionLoad",
		InstructionLoadStay:   "InstructionLoadStay",
		InstructionFunc:       "InstructionFunc",
		InstructionCall:       "InstructionCall",
		InstructionType:       "InstructionType",
		InstructionConst:      "InstructionConst",
		InstructionVar:        "InstructionVar",
		InstructionGvar:       "InstructionGvar",
		InstructionHold:       "InstructionHold",
		InstructionNop:        "InstructionNop",
		InstructionProc:       "InstructionProc",
		InstructionIn:         "InstructionIn",
		InstructionInvoke:     "InstructionInvoke",
		InstructionLShift:     "InstructionLShift",
		InstructionRShift:     "InstructionRShift",
		InstructionLAnd:       "InstructionLAnd",
		InstructionLOr:        "InstructionLOr",
		InstructionOver:       "InstructionOver",
		InstructionExit:       "InstructionExit",
		InstructionParam:      "InstructionParam",
		InstructionReturnType: "InstructionReturnType",
		InstructionDeref:      "InstructionDeref",
		FunctionPutU64:        "FunctionPutU64",
		FunctionPutString:     "FunctionPutString",
		FunctionDumpMem:       "FunctionDumpMem",
		FunctionSetNumber:     "FunctionSetNumber",
		FunctionSetString:     "FunctionSetString",
		FunctionAlloc:         "FunctionAlloc",
		FunctionFree:          "FunctionFree",
		FunctionPutChar:       "FunctionPutChar",
		Skip:                  "Skip",
	}
	for ins := InstructionInvalid + 1; ins < TotalOps-1; ins++ {
		_, found := instructionNames[ins]
		if !found {
			panic(fmt.Sprintf("[DEV] Missing instruction %d on name map", ins))
		}
	}
	if len(instructionNames) != int(TotalOps)-1 {
		panic("[DEV] Missing instruction on name map")
	}
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
	Procedures []Operation
	Error      []error
	Variables  []Operation
	Constants  []Operation
	Operations []Operation
}

func (p *Program) FindProcByOperand(operand string) (*Operation, error) {
	if operand == "" {
		return nil, errors.New("invalid operand")
	}
	if len(p.Procedures) == 0 {
		return nil, fmt.Errorf("proc %q not found", operand)
	}
	i, found := slices.BinarySearchFunc(p.Procedures, operand, func(a Operation, b string) int {
		return strings.Compare(a.Operator.Operand, b)
	})
	if !found {
		return nil, errors.New("proc does not exist")
	}
	return &p.Procedures[i], nil
}

func (p *Program) FindProc(operand Operation) (*Operation, error) {
	if operand.Operator.Operand == "" {
		return nil, errors.New("invalid operand")
	}
	if len(p.Procedures) == 0 {
		return nil, fmt.Errorf("proc %q not found", operand.Operator.Operand)
	}
	i, found := slices.BinarySearchFunc(p.Procedures, operand, func(a Operation, b Operation) int {
		return strings.Compare(a.Operator.Operand, b.Operator.Operand)
	})
	if !found {
		return nil, errors.New("proc does not exist")
	}
	return &p.Procedures[i], nil
}

func (p *Program) AppendProc(proc Operation) {
	p.Procedures = append(p.Procedures, proc)
	slices.SortFunc(p.Procedures, func(a, b Operation) int {
		return strings.Compare(a.Operator.Operand, b.Operator.Operand)
	})
}

func PPrintOperation(op Operation) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", InstructionToStr(op.Instruction)))
	builder.WriteString(fmt.Sprintf("	operand: %s | symbolName%q\n", op.Operator.Operand, op.Operator.SymbolName))
	for k, v := range op.Links {
		builder.WriteString(fmt.Sprintf("	link_name: %q | link_type: %q | link_value: %q\n", k, v.Operator.SymbolName, v.Operator.Operand))
	}
	for k, v := range op.Addresses {
		builder.WriteString(fmt.Sprintf("\nAddr %s: %d\n", InstructionToStr(k), v))
	}
	builder.WriteString("****************************************************\n")
	return builder.String()
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
