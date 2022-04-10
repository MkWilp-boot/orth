package functions

import (
	"t/cmd/pkg/helpers"
	orthtypes "t/cmd/pkg/types"
)

var Functions map[string]func(*[]orthtypes.Operand)

func init() {
	Functions = make(map[string]func(*[]orthtypes.Operand))
	Functions["to_string"] = func(stack *[]orthtypes.Operand) {
		o1 := helpers.StackPop(&(*stack))
		res := ToString(o1)
		*stack = append(*stack, res)
	}
}
