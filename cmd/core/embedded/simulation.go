package embedded

import (
	"fmt"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers"
	"orth/cmd/pkg/helpers/functions"
	orthtypes "orth/cmd/pkg/types"
	"strings"
)

const memCap = 64000

// Simulate runs a program built by previous steps
func Simulate(program orthtypes.Program) {
	stack := make([]orthtypes.Operand, 0, memCap)
	mem := make([]orthtypes.Operand, memCap)

	vars := make(map[string]orthtypes.Operand)

	ip := 0
	for ip < len(program.Operations) {
		stackItem := program.Operations[ip]

		if !stackItem.IsValidType() {
			fmt.Println("====================================")
			fmt.Printf("Error At instruction number %d ->'%#v'", ip, stackItem)
			fmt.Printf("The argument %s has a invalid type of %q\n", stackItem.Operator, stackItem.Operator.SymbolName)
			fmt.Println("====================================")
			panic(orth_debug.DefaultRuntimeException)
		}
		switch stackItem.Instruction {
		case orthtypes.Push:
			if len(stack) >= memCap {
				panic("stack overflow!")
			}
			stack = append(stack, stackItem.Operator)
			ip++
		case orthtypes.Sum:
			o1 := helpers.StackPop(&stack)
			o2 := helpers.StackPop(&stack)

			if o2.SymbolName == orthtypes.RNGABL {
				start, end := functions.DissectRangeAsInt(o2)
				if !helpers.IsInt(o1.SymbolName) {
					panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, o1.SymbolName, "sum rangeable"))
				}
				sum := start + helpers.ToInt(o1)

				teste := orthtypes.Operand{
					SymbolName: orthtypes.RNGABL,
					Operand:    fmt.Sprintf("%d|%d", sum, end),
				}

				stack = append(stack, teste)
				ip++
				continue
			}

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)

			fun, err := functions.SumBasedOnType(superType)
			if err != nil {
				fmt.Println(err)
				functions.DumpStack(stack)
			}

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
			var o1 orthtypes.Operand
			var o2 orthtypes.Operand

			o1 = helpers.StackPop(&stack)

			if o1.SymbolName == orthtypes.RNGABL {
				o1, o2 = functions.DissectRange(o1)
			} else {
				o2 = helpers.StackPop(&stack)
			}

			superType := functions.GetSupersetType(o1, o2)

			fun := functions.LowerBasedOnType(superType)
			stack = append(stack, fun(superType, o1, o2))
			ip++
		case orthtypes.Gt:
			var o1 orthtypes.Operand
			var o2 orthtypes.Operand

			o1 = helpers.StackPop(&stack)

			if o1.SymbolName == orthtypes.RNGABL {
				o1, o2 = functions.DissectRange(o1)
			} else {
				o2 = helpers.StackPop(&stack)
			}

			helpers.SameBaseType(o1, o2)

			superType := functions.GetSupersetType(o1, o2)
			fun := functions.GreaterBasedOnType(superType)

			stack = append(stack, fun(superType, o1, o2))
			ip++
		case orthtypes.If:
			o1 := helpers.StackPop(&stack)

			if o1.SymbolName != "b" {
				panic(orth_debug.InvalidBoolType)
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
		case orthtypes.PutU64:
			o1 := helpers.StackPop(&stack)
			if !helpers.IsInt(o1.SymbolName) {
				panic(fmt.Sprintf(orth_debug.InvalidTypeForInstruction, o1.SymbolName, "DumpUI64"))
			}

			fmt.Printf("%d\n", helpers.ToInt(o1))
			ip++
		case orthtypes.PutString:
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
			if o1.SymbolName != orthtypes.PrimitiveBOOL {
				panic(orth_debug.InvalidBoolType)
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

			if address.SymbolName != orthtypes.ADDR {
				panic(fmt.Errorf(orth_debug.InvalidTypeForIndex, orthtypes.ADDR))
			}

			mem[helpers.ToInt(address)] = value
			ip++
		case orthtypes.Load:
			address := helpers.StackPop(&stack)
			stack = append(stack, mem[helpers.ToInt(address)])
			mem[helpers.ToInt(address)] = orthtypes.Operand{}
			ip++
		case orthtypes.OType:
			//fmt.Println(stack)
			ip++
		case orthtypes.Call:
			fn := functions.Functions[stackItem.Operator.Operand]
			fn(&stack, &mem, vars)
			ip++
		case orthtypes.LoadStay:
			address := helpers.StackPop(&stack)
			stack = append(stack, mem[helpers.ToInt(address)])
			ip++
		case orthtypes.Var:
			//			vName					Value
			vars[stackItem.Operator.Operand] = helpers.StackPop(&stack)
			stack = append(stack, orthtypes.Operand{
				SymbolName: orthtypes.PrimitiveConst,
				Operand:    stackItem.Operator.Operand,
			})
			ip++
		case orthtypes.Hold:
			vName := stackItem.Operator.Operand
			v, ok := vars[vName]
			if !ok {
				panic(fmt.Errorf(orth_debug.VariableUndefined, vName))
			}
			stack = append(stack, v)
			ip++
		default:
			panic(fmt.Errorf(orth_debug.InvalidInstruction, stackItem.Instruction))
		}
	}
}
