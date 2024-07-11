package embedded_helpers

import (
	"errors"
	"fmt"
	"log"
	"os"

	orth_types "orth/cmd/pkg/types"

	"golang.org/x/exp/constraints"
)

const MainScope = "_global"

type RefStackItem struct {
	AbsPosition uint
	Instruction orth_types.Instruction
}

func PopLast[T any](stack *[]T) T {
	var defaultValue T
	if len(*stack) == 0 {
		return defaultValue
	}
	item := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	return item
}

func HandleOperationDo(stack *[]RefStackItem, program *orth_types.Program, operationIndex uint) {
	lastStackItem := PopLast(stack)
	switch lastStackItem.Instruction {
	case orth_types.InstructionWhile:
		program.Operations[operationIndex].Addresses[orth_types.InstructionWhile] = int(lastStackItem.AbsPosition)
	default:
		log.Fatalln("non logical block found before 'DO' operation, syntax error")
	}
}

func HandleOperationEnd(stack *[]RefStackItem, program *orth_types.Program, currentOperationIndex uint) {
	lastStackItem := PopLast(stack)
	switch lastStackItem.Instruction {
	case orth_types.InstructionIf:
		ifOperation := lastStackItem

		program.Operations[ifOperation.AbsPosition].Addresses[orth_types.InstructionEnd] = int(currentOperationIndex)
		program.Operations[currentOperationIndex].Addresses[orth_types.InstructionIf] = int(ifOperation.AbsPosition)
	case orth_types.InstructionElse:
		elseOperation := lastStackItem

		program.Operations[elseOperation.AbsPosition].Addresses[orth_types.InstructionEnd] = int(currentOperationIndex)
		program.Operations[currentOperationIndex].Addresses[orth_types.InstructionElse] = int(elseOperation.AbsPosition)
	case orth_types.InstructionProc:
		procOperation := lastStackItem
		program.Operations[currentOperationIndex].Addresses[orth_types.InstructionProc] = int(procOperation.AbsPosition)
	case orth_types.InstructionDo:
		doOperation := lastStackItem
		whileAddress := program.Operations[doOperation.AbsPosition].Addresses[orth_types.InstructionWhile]

		program.Operations[currentOperationIndex].Addresses[orth_types.InstructionWhile] = int(whileAddress)
		program.Operations[doOperation.AbsPosition].Addresses[orth_types.InstructionEnd] = int(currentOperationIndex)
	}
}

func HandleOperationElse(stack *[]RefStackItem, program *orth_types.Program, operationIndex uint) {
	lastStackItem := PopLast(stack)
	switch lastStackItem.Instruction {
	case orth_types.InstructionIf:
		ifOperation := lastStackItem
		program.Operations[ifOperation.AbsPosition].Addresses[orth_types.InstructionElse] = int(operationIndex)
	}
}

func GetVariableContext(variable orth_types.ContextDeclaration, context *orth_types.Context) (string, error) {
	if context == nil {
		return "", fmt.Errorf("undefined variable at abs location: %d", variable.Index)
	}
	for _, declaration := range context.Declarations {
		if declaration.Name == variable.Name {
			return context.Name, nil
		}
	}
	return GetVariableContext(variable, context.Parent)
}

func LinkVariableToValue(operation orth_types.Operation, analyzerOperations *[]orth_types.Operation, program *orth_types.Program) orth_types.Operation {
	// if operation.Instruction == orth_types.Var {
	// 	// set to skip so the value won't be on the final asm
	// 	(*analyzerOperations)[len(*analyzerOperations)-1].Instruction = orth_types.Skip
	// 	// set to len - 1 because the last element will always be the var value
	// 	operation.RefBlock = len(*analyzerOperations) - 1
	// 	program.Variables = append(program.Variables, operation)
	// } else if operation.Instruction == orth_types.Const {
	// 	// set to skip so the value won't be on the final asm
	// 	(*analyzerOperations)[len(*analyzerOperations)-1].Instruction = orth_types.Skip
	// 	// set to len - 1 because the last element will always be the var value
	// 	operation.RefBlock = len(*analyzerOperations) - 1
	// 	program.Constants = append(program.Constants, operation)
	// } else {
	// 	if operation.Instruction == orth_types.Hold {
	// 		operationType, err := OperationIsVariableLike(operation, program)
	// 		if err != nil {
	// 			program.Error = append(program.Error, err)
	// 			return orth_types.Operation{}
	// 		}
	// 		operation.Operator.SymbolName = operationType
	// 	}
	// }
	return operation
}

func OperationIsVariableLike(operation orth_types.Operation, program *orth_types.Program) (string, error) {
	for _, v := range program.Variables {
		if v.Operator.Operand == operation.Operator.Operand {
			return orth_types.StdVar, nil
		}
	}
	if operation.Operator.SymbolName != orth_types.StdVar {
		for _, v := range program.Constants {
			if v.Operator.Operand == operation.Operator.Operand {
				return orth_types.StdConst, nil
			}
		}
	}
	err := errors.New("could not define if holding istruction is set to a variable or const")
	return "", err
}

func ProduceOperator[TOperand constraints.Float | constraints.Integer](param1, param2 TOperand, instruction orth_types.Instruction) (string, bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	operand := ""
	if instruction == orth_types.InstructionMult {
		operand = fmt.Sprint(param1 * param2)
	} else if instruction == orth_types.InstructionSum {
		operand = fmt.Sprint(param1 + param2)
	} else if instruction == orth_types.InstructionMod {
		var param1Inter interface{} = param1
		switch param1Inter.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			operand = fmt.Sprint(int64(param1) % int64(param2))
		default:
			panic("modulo operation is only supported for integer types.")
		}
	} else if instruction == orth_types.InstructionDiv {
		operand = fmt.Sprint(param1 / param2)
	} else if instruction == orth_types.InstructionMinus {
		operand = fmt.Sprint(param1 - param2)
	}

	return operand, operand != ""
}
