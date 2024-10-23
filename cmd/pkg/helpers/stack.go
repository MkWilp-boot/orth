package helpers

import (
	"fmt"
	orth_types "orth/cmd/pkg/types"
)

const MEM_MAX_CAP uint = 64000

// StackPop pops the last item from the stack
func StackPop(root *[]orth_types.Operand) orth_types.Operand {
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
	root *[]orth_types.Operand,
	superType string,
	act func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand,
	operands ...orth_types.Operand) {

	stack := *root
	originalType := operands[0].GrabRootType()
	switch originalType {
	case orth_types.INTS:
		stack = append(stack, act(superType, operands[0], operands[1]))
	case orth_types.FLOATS:
		stack = append(stack, act(superType, operands[0], operands[1]))
	case orth_types.STRING:
		stack = append(stack, act(superType, operands[0], operands[1]))
	case orth_types.RNT:
		stack = append(stack, act(superType, operands[0], operands[1]))
	default:
		panic(fmt.Errorf("invalid type of %q", originalType))
	}

	*root = stack
}
