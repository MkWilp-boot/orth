package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"t/cmd/core/debug"
	orthtypes "t/cmd/pkg/types"
)

// StackPop pops the last item from the stack
func StackPop(root *[]orthtypes.Operand) orthtypes.Operand {
	if len(*root) < 1 {
		panic("Stack underflow error!")
	}
	stack := *root

	ret := stack[len(stack)-1]
	*root = stack[:len(stack)-1]
	return ret
}

func DissectRange(o1 orthtypes.Operand) (orthtypes.Operand, orthtypes.Operand) {
	if o1.VarType != orthtypes.RNGABL {
		panic(fmt.Errorf(debug.InvalidTypeForInstruction, o1.VarType, "DissectRange"))
	}
	nums := strings.Split(o1.Operand, "|")
	start, _ := strconv.Atoi(nums[0])
	end, _ := strconv.Atoi(nums[1])

	return orthtypes.Operand{
			VarType: orthtypes.PrimitiveI32,
			Operand: fmt.Sprint(start),
		}, orthtypes.Operand{
			VarType: orthtypes.PrimitiveI32,
			Operand: fmt.Sprint(end),
		}
}

// BasedOnType executes a 'act' and appends it's result to the 'root' (or stack)
func BasedOnType(
	root *[]orthtypes.Operand,
	superType string,
	act func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand,
	operands ...orthtypes.Operand) {

	stack := *root
	originalType := operands[0].GrabRootType()
	switch originalType {
	case orthtypes.INTS:
		stack = append(stack, act(superType, operands[0], operands[1]))
	case orthtypes.FLOATS:
		stack = append(stack, act(superType, operands[0], operands[1]))
	case orthtypes.STRING:
		stack = append(stack, act(superType, operands[0], operands[1]))
	case orthtypes.RNT:
		stack = append(stack, act(superType, operands[0], operands[1]))
	default:
		panic(fmt.Errorf("INVALID TYPE OF %q\n", originalType))
	}

	*root = stack
}
