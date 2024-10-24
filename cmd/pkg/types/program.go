package orth_types

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
	InstructionIn
	InstructionInvoke
	FunctionDumpMem
	InstructionLShift
	InstructionRShift
	InstructionLAnd
	InstructionLOr
	InstructionOver
	InstructionExit
	InstructionWith
	InstructionOut
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
		Skip:                "Skip",
		InstructionPush:     "Push",
		InstructionPushStr:  "PushStr",
		InstructionSum:      "Sum",
		InstructionMinus:    "Minus",
		InstructionMult:     "Mult",
		InstructionDiv:      "Div",
		InstructionIf:       "If",
		InstructionElse:     "Else",
		InstructionEnd:      "End",
		InstructionEqual:    "Equal",
		InstructionLt:       "Lt",
		InstructionGt:       "Gt",
		InstructionNotEqual: "NotEqual",
		InstructionDup:      "Dup",
		InstructionTwoDup:   "TwoDup",
		InstructionDo:       "Do",
		InstructionDrop:     "Drop",
		InstructionWhile:    "While",
		InstructionSwap:     "Swap",
		InstructionMod:      "Mod",
		InstructionMem:      "Mem",
		InstructionStore:    "Store",
		InstructionLoad:     "Load",
		InstructionLoadStay: "LoadStay",
		InstructionFunc:     "Func",
		InstructionCall:     "Call",
		InstructionType:     "Type",
		InstructionConst:    "Const",
		InstructionVar:      "Var",
		InstructionGvar:     "Gvar",
		InstructionHold:     "Hold",
		InstructionNop:      "Nop",
		InstructionProc:     "Proc",
		InstructionIn:       "In",
		InstructionInvoke:   "Invoke",
		InstructionLShift:   "LShift",
		InstructionRShift:   "RShift",
		InstructionLAnd:     "LAnd",
		InstructionLOr:      "LOr",
		InstructionOver:     "Over",
		InstructionExit:     "Exit",
		InstructionParam:    "Param",
		InstructionWith:     ":",
		InstructionOut:      "--",
		InstructionDeref:    "Deref",
		FunctionPutU64:      "PutU64",
		FunctionPutString:   "PutString",
		FunctionDumpMem:     "DumpMem",
		FunctionSetNumber:   "SetNumber",
		FunctionSetString:   "SetString",
		FunctionAlloc:       "Alloc",
		FunctionFree:        "Free",
		FunctionPutChar:     "PutChar",
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
	Error      []error
	Variables  []Operation
	Constants  []Operation
	Operations []Operation
}

type ProcedureSchema struct {
	InParamsAmount, OutParamsAmount []Operation
}

func (p *Program) FindProc(operation Operation) (ProcedureSchema, error) {
	callingProcedureArguments := make([]Operation, 0)
	callingProcedureOutParams := make([]Operation, 0)

	for callingProcedureIndex, op := range p.Operations {
		if op.Operator.Operand == operation.Operator.Operand && op.Instruction == InstructionProc {
			for _, operation := range p.Operations[callingProcedureIndex:] {
				if operation.Instruction == InstructionWith {
					for k, v := range operation.Links {
						if !strings.HasPrefix(k, "proc_param_") {
							continue
						}
						callingProcedureArguments = append(callingProcedureArguments, v)
					}
				}
				if operation.Instruction == InstructionOut {
					for k, v := range operation.Links {
						if !strings.HasPrefix(k, "proc_out_param_") {
							continue
						}
						callingProcedureOutParams = append(callingProcedureOutParams, v)
					}
				}
			}

			return ProcedureSchema{
				InParamsAmount:  callingProcedureArguments,
				OutParamsAmount: callingProcedureOutParams,
			}, nil
		}
	}
	return ProcedureSchema{}, fmt.Errorf("proc named %q was not found", operation.Operator.Operand)
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
