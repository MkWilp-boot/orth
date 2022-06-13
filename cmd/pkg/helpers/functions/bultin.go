package functions

import (
	"fmt"
	"t/cmd/core/debug"
	"t/cmd/pkg/helpers"
	orthtypes "t/cmd/pkg/types"
)

var Functions map[string]func(stack *[]orthtypes.Operand) (bool, func(mem *[]orthtypes.Operand))

func init() {
	Functions = make(map[string]func(*[]orthtypes.Operand) (bool, func(mem *[]orthtypes.Operand)))

	Functions["to_string"] = func(stack *[]orthtypes.Operand) (bool, func(mem *[]orthtypes.Operand)) {
		o1 := helpers.StackPop(&*stack)
		res := ToString(o1)
		*stack = append(*stack, res)
		return false, nil
	}

	Functions["length_of"] = func(stack *[]orthtypes.Operand) (bool, func(mem *[]orthtypes.Operand)) {
		o1 := helpers.StackPop(&*stack)
		switch o1.VarType {
		case orthtypes.PrimitiveSTR:
			*stack = append(*stack, orthtypes.Operand{
				VarType: orthtypes.PrimitiveI32,
				Operand: fmt.Sprint(len(o1.Operand)),
			})
		default:
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, o1.VarType, "Functions[length_of]"))
		}
		return false, nil
	}

	Functions["make_array"] = func(stack *[]orthtypes.Operand) (bool, func(mem *[]orthtypes.Operand)) {
		if len(*stack) < 2 {
			panic(debug.StackUnderFlow)
		}

		capacity := helpers.ToInt(helpers.StackPop(&*stack))
		typ := helpers.StackPop(&*stack)

		return true, func(originalMem *[]orthtypes.Operand) {
			// why go?
			memCopy := make([]orthtypes.Operand, len(*originalMem), cap(*originalMem))
			copy(memCopy, *originalMem)

			var start int
			var foundPlace bool

			// look for a place where all members can fit by direct index
			for i := range memCopy {
				fitInAmount := 0
				slice := memCopy[i:(i + capacity)]

				for _, xx := range slice {
					if xx.Operand == "" {
						fitInAmount++
					}
				}

				if fitInAmount == capacity {
					foundPlace = true
					start = i
					break
				}
			}

			if foundPlace {
				for i := start; i < capacity; i++ {
					memCopy[i] = orthtypes.Operand{
						VarType: typ.Operand,
					}
				}
				*originalMem = memCopy
				*stack = append(*stack, orthtypes.Operand{
					VarType: orthtypes.ADDR,
					Operand: fmt.Sprint(capacity),
				}, orthtypes.Operand{
					VarType: orthtypes.ADDR,
					Operand: fmt.Sprint(start),
				})
			}
		}
	}
}
