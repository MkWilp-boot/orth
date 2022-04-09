package simulation

import (
	"fmt"
	"strings"
	"t/cmd/core/debug"
	"t/cmd/pkg/helpers"
	"t/cmd/pkg/helpers/functions"
	orthtypes "t/cmd/pkg/types"
)

const memCap = 10

func Simulate(program orthtypes.Program) {
	stack := make([]orthtypes.Operand, 0)
	mem := make([]orthtypes.Operand, memCap)

	ip := 0
	for ip < len(program.Operations) {
		stackItem := program.Operations[ip]

		if !stackItem.IsValidType() {
			fmt.Println("====================================")
			fmt.Printf("Error At instruction number %d ->'%#v'", ip, stackItem)
			fmt.Printf("The argument %s has a invalid type of %q\n", stackItem.Operand, stackItem.Operand.VarType)
			fmt.Println("====================================")
			panic(debug.DefaultRuntimeException)
		}

		switch stackItem.Instruction {
		case orthtypes.Push:
			stack = append(stack, stackItem.Operand)
			ip++
		case orthtypes.Sum:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)

			fun := functions.SumBasedOnType(superType)

			helpers.BasedOnType(&stack, superType, fun, o1, o2)
			ip++
		case orthtypes.Minus:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.SubBasedOnType(superType)

			helpers.BasedOnType(&stack, superType, fun, o2, o1)
			ip++
		case orthtypes.Div:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.DivBasedOnType(superType)

			helpers.BasedOnType(&stack, superType, fun, o2, o1)
			ip++
		case orthtypes.Mult:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.MultBasedOnType(superType)

			helpers.BasedOnType(&stack, superType, fun, o2, o1)
			ip++
		case orthtypes.Equal:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.EqualBasedOnType(superType)

			stack = append(stack, fun(superType, o1, o2))
			ip++
		case orthtypes.NotEqual:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.NotEqualBasedOnType(superType)

			stack = append(stack, fun(superType, o1, o2))
			ip++
		case orthtypes.Lt:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.LowerBasedOnType(superType)

			stack = append(stack, fun(superType, o1, o2))
			ip++
		case orthtypes.Gt:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.GreaterBasedOnType(superType)

			stack = append(stack, fun(superType, o1, o2))
			ip++
		case orthtypes.If:
			o1 := helpers.StackPop(&stack)

			if o1.VarType != "b" {
				panic(debug.InvalidBoolType)
			}

			if o1.Operand == orthtypes.StdTrue {
				ip++ // next intruction
			} else {
				ip = stackItem.RefBlock
			}
		case orthtypes.Else:
			ip = stackItem.RefBlock
		case orthtypes.End:
			ip = stackItem.RefBlock
		case orthtypes.Dump:
			o1 := helpers.StackPop(&stack)
			fmt.Printf("%#v\n", o1)
			ip++
		case orthtypes.Print:
			o1 := helpers.StackPop(&stack)

			if strings.HasSuffix(o1.Operand, "\\n") {
				fmt.Printf("%v\n", o1.Operand[:len(o1.Operand)-2])
			} else {
				fmt.Printf("%v", o1.Operand)
			}
			ip++
		case orthtypes.Dup:
			o1 := helpers.StackPop(&stack)
			stack = append(stack, o1, o1)
			ip++
		case orthtypes.Do:
			o1 := helpers.StackPop(&stack)
			if o1.VarType != orthtypes.PrimitiveBOOL {
				panic(debug.InvalidBoolType)
			}

			if o1.Operand == orthtypes.StdTrue {
				ip++
			} else {
				ip = stackItem.RefBlock
			}
		case orthtypes.While:
			ip++
		case orthtypes.Drop:
			helpers.StackPop(&stack) // just pops the value
			ip++
		case orthtypes.Swap:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			stack = append(stack, o1, o2)

			ip++
		case orthtypes.Mod:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)

			fun := functions.ModBasedOnType(superType)

			helpers.BasedOnType(&stack, superType, fun, o2, o1)
			ip++
		case orthtypes.Mem:
			ip++
		case orthtypes.Store:
			value := helpers.StackPop(&stack)
			address := helpers.StackPop(&stack)

			if !helpers.IsInt(address.VarType) {
				panic(fmt.Errorf(debug.InvalidTypeForIndex, orthtypes.PrimitiveInt))
			}

			mem[helpers.ToInt(address)] = value
			ip++
		case orthtypes.Load:
			address := helpers.StackPop(&stack)
			stack = append(stack, mem[helpers.ToInt(address)])
			mem[helpers.ToInt(address)] = orthtypes.Operand{}
			ip++
		default:
			panic(fmt.Errorf(debug.InvalidInstruction, stackItem.Instruction))
		}
	}
}
