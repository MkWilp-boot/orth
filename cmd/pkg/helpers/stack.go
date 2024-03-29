package helpers

import (
	"fmt"
	orthtypes "orth/cmd/pkg/types"
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
		panic(fmt.Errorf("invalid type of %q", originalType))
	}

	*root = stack
}
