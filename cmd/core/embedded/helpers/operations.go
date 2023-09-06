package embedded_helpers

import (
	"errors"
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

func GetVariableContext(varName string, context *orthtypes.Context) (string, error) {
	if context == nil {
		return "", errors.New("variable is undeclared: " + varName)
	}
	for _, declaration := range context.Declarations {
		if declaration == varName {
			return context.Name, nil
		}
	}
	return GetVariableContext(varName, context.Parent)
}

func LinkVariableToValue(operation orthtypes.Operation, analyzerOperations *[]orthtypes.Operation, program *orthtypes.Program) orthtypes.Operation {
	if operation.Instruction == orthtypes.Var {
		// set to skip so the value won't be on the final asm
		(*analyzerOperations)[len(*analyzerOperations)-1].Instruction = orthtypes.Skip
		// set to len - 1 because the last element will always be the var value
		operation.RefBlock = len(*analyzerOperations) - 1
		program.Variables = append(program.Variables, operation)
	} else if operation.Instruction == orthtypes.Const {
		// set to skip so the value won't be on the final asm
		(*analyzerOperations)[len(*analyzerOperations)-1].Instruction = orthtypes.Skip
		// set to len - 1 because the last element will always be the var value
		operation.RefBlock = len(*analyzerOperations) - 1
		program.Constants = append(program.Constants, operation)
	} else {
		if operation.Instruction == orthtypes.Hold {
			operationType, err := OperationIsVariableLike(operation, program)
			if err != nil {
				program.Error = append(program.Error, err)
				return orthtypes.Operation{}
			}
			operation.Operator.SymbolName = operationType
		}
	}
	return operation
}

func OperationIsVariableLike(operation orthtypes.Operation, program *orthtypes.Program) (string, error) {
	for _, v := range program.Variables {
		if v.Operator.Operand == operation.Operator.Operand {
			return orthtypes.PrimitiveVar, nil
		}
	}
	if operation.Operator.SymbolName != orthtypes.PrimitiveVar {
		for _, v := range program.Constants {
			if v.Operator.Operand == operation.Operator.Operand {
				return orthtypes.PrimitiveConst, nil
			}
		}
	}
	err := errors.New("could not define if holding istruction is set to a variable or const")
	return "", err
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
