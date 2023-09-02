package optimizer

import (
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
	"strconv"
)

func AnalyzeAndOptimizeOperations(operations []orthtypes.Operation) ([]orthtypes.Operation, []orthtypes.CompilerMessage) {
	stack := make([]orthtypes.Operation, 0)
	warnings := make([]orthtypes.CompilerMessage, 0)

	for _, operation := range operations {
		switch operation.Instruction {
		case orthtypes.Mult:
			fallthrough
		case orthtypes.Mod:
			fallthrough
		case orthtypes.Div:
			fallthrough
		case orthtypes.Minus:
			fallthrough
		case orthtypes.Sum:
			if stack[len(stack)-1].Instruction == orthtypes.Push && stack[len(stack)-2].Instruction == orthtypes.Push {
				p1 := embedded_helpers.PopLast(&stack)
				p2 := embedded_helpers.PopLast(&stack)

				if p1.Operator.VarType != p2.Operator.VarType {
					msg := orth_debug.BuildMessage(
						orth_debug.ORTH_WARN_01,
						orthtypes.InstructionToStr(operation.Instruction),
						p1.Operator.VarType,
						p2.Operator.VarType,
					)
					warnings = append(warnings, orthtypes.CompilerMessage{
						Type:    orthtypes.Commom,
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

					stack = append(stack, orthtypes.Operation{
						Instruction: orthtypes.Push,
						Operator: orthtypes.Operand{
							VarType: orthtypes.PrimitiveInt,
							Operand: operand,
						},
						Context:  operation.Context,
						RefBlock: operation.RefBlock,
					})
					continue
				}
			} else if stack[len(stack)-1].Instruction == orthtypes.PushStr && stack[len(stack)-2].Instruction == orthtypes.PushStr {
				p1 := embedded_helpers.PopLast(&stack)
				p2 := embedded_helpers.PopLast(&stack)
				stack = append(stack, orthtypes.Operation{
					Instruction: orthtypes.PushStr,
					Operator: orthtypes.Operand{
						VarType: orthtypes.PrimitiveSTR,
						Operand: p2.Operator.Operand + p1.Operator.Operand, // concat
					},
					Context:  operation.Context,
					RefBlock: operation.RefBlock,
				})
				continue
			}

		}
		stack = append(stack, operation)
	}

	return stack, warnings
}
