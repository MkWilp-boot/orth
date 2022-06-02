package functions

import (
	"fmt"
	"t/cmd/core/debug"
	"t/cmd/pkg/helpers"
	orthtypes "t/cmd/pkg/types"
)

var Functions map[string]func(stack *[]orthtypes.Operand)

func init() {
	Functions = make(map[string]func(*[]orthtypes.Operand))

	Functions["to_string"] = func(stack *[]orthtypes.Operand) {
		o1 := helpers.StackPop(&*stack)
		res := ToString(o1)
		*stack = append(*stack, res)
	}

	Functions["length_of"] = func(stack *[]orthtypes.Operand) {
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
	}
}
