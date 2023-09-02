package embedded_helpers

import (
	"fmt"
	"os"

	orthtypes "orth/cmd/pkg/types"

	"golang.org/x/exp/constraints"
)

const MainScope = "_global"

func PopLast[T comparable](root *[]T) T {
	stack := *root
	ret := stack[len(stack)-1]
	*root = stack[:len(stack)-1]
	return ret
}

func ProduceOperator[TOperand constraints.Float | constraints.Integer](param1, param2 TOperand, instruction int) (string, bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	operand := ""
	if instruction == orthtypes.Mult {
		operand = fmt.Sprint(param1 * param2)
	} else if instruction == orthtypes.Sum {
		operand = fmt.Sprint(param1 + param2)
	} else if instruction == orthtypes.Mod {
		var param1Inter interface{} = param1
		switch param1Inter.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			operand = fmt.Sprint(int64(param1) % int64(param2))
		default:
			panic("modulo operation is only supported for integer types.")
		}
	} else if instruction == orthtypes.Div {
		operand = fmt.Sprint(param1 / param2)
	} else if instruction == orthtypes.Minus {
		operand = fmt.Sprint(param1 - param2)
	}

	return operand, operand != ""
}
