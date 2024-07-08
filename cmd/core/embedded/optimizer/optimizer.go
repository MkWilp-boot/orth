package optimizer

import (
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"strconv"
)

func AnalyzeAndOptimizeOperations(operations []orth_types.Operation) ([]orth_types.Operation, []orth_types.CompilerMessage) {
	stack := make([]orth_types.Operation, 0)
	warnings := make([]orth_types.CompilerMessage, 0)

	for _, operation := range operations {
		switch operation.Instruction {
		case orth_types.InstructionMult:
			fallthrough
		case orth_types.InstructionMod:
			fallthrough
		case orth_types.InstructionDiv:
			fallthrough
		case orth_types.InstructionMinus:
			fallthrough
		case orth_types.InstructionSum:
			if stack[len(stack)-1].Instruction == orth_types.InstructionPush && stack[len(stack)-2].Instruction == orth_types.InstructionPush {
				p1 := embedded_helpers.PopLast(&stack)
				p2 := embedded_helpers.PopLast(&stack)

				if p1.Operator.SymbolName != p2.Operator.SymbolName {
					msg := orth_debug.BuildMessage(
						orth_debug.ORTH_WARN_01,
						orth_types.InstructionToStr(operation.Instruction),
						p1.Operator.SymbolName,
						p2.Operator.SymbolName,
					)
					warnings = append(warnings, orth_types.CompilerMessage{
						Type:    orth_types.Commom,
						Message: msg,
					})
				}

				if p1.IsNumeric() && p2.IsNumeric() {
					operand := operation.Operator.Operand
					if p1.IsInt() && p2.IsInt() {
						param1, _ := strconv.Atoi(p1.Operator.Operand)
						param2, _ := strconv.Atoi(p2.Operator.Operand)

						if op, ok := embedded_helpers.ProduceOperator(param1, param2, operation.Instruction); ok {
							operand = op
						}

					} else if p1.IsFloat() && p2.IsFloat() {
						p1BitSize := 64
						p2BitSize := 64
						if p1.IsFloat32() {
							p1BitSize = 32
						}
						if p2.IsFloat32() {
							p2BitSize = 32
						}
						param1, _ := strconv.ParseFloat(p1.Operator.Operand, p1BitSize)
						param2, _ := strconv.ParseFloat(p2.Operator.Operand, p2BitSize)

						if op, ok := embedded_helpers.ProduceOperator(param1, param2, operation.Instruction); ok {
							operand = op
						}
					}

					stack = append(stack, orth_types.Operation{
						Instruction: orth_types.InstructionPush,
						Operator: orth_types.Operand{
							SymbolName: orth_types.StdINT,
							Operand:    operand,
						},
						Context:   operation.Context,
						Addresses: operation.Addresses,
					})
					continue
				}
			} else if stack[len(stack)-1].Instruction == orth_types.InstructionPushStr && stack[len(stack)-2].Instruction == orth_types.InstructionPushStr {
				p1 := embedded_helpers.PopLast(&stack)
				p2 := embedded_helpers.PopLast(&stack)
				stack = append(stack, orth_types.Operation{
					Instruction: orth_types.InstructionPushStr,
					Operator: orth_types.Operand{
						SymbolName: orth_types.StdSTR,
						Operand:    p2.Operator.Operand + p1.Operator.Operand, // concat
					},
					Context:   operation.Context,
					Addresses: operation.Addresses,
				})
				continue
			}

		}
		stack = append(stack, operation)
	}

	return stack, warnings
}
