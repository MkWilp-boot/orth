package simulation

import (
	"fmt"
	"orth/cmd/pkg/helpers"
	"orth/cmd/pkg/helpers/functions"
	orth_types "orth/cmd/pkg/types"
	"os"
)

type doubleOperandsOperationtionGroup struct {
	Integer func(superType string, n1, n2 orth_types.Operand) orth_types.Operand
	Float   func(superType string, n1, n2 orth_types.Operand) orth_types.Operand
}

type stack struct {
	items []orth_types.Operation
	ptr   int
}

func (s *stack) push(itens ...orth_types.Operation) {
	for _, item := range itens {
		s.ptr++
		(s.items)[s.ptr] = item
	}
}

func (s *stack) rmv(numItensToRemove int) {
	var defaultValue orth_types.Operation
	for i := 0; i < numItensToRemove; i++ {
		if s.ptr < 0 {
			panic("stack underflow")
		}
		(s.items)[s.ptr] = defaultValue
		s.ptr--
	}
}

func (s *stack) peek(numItens int) []orth_types.Operation {
	preview := make([]orth_types.Operation, numItens)
	copy(preview, (s.items)[:numItens])
	return preview
}

func addToMem(stack *[]orth_types.Operation, offset int, itens ...orth_types.Operation) {
	for i, item := range itens {
		(*stack)[offset+i] = item
	}
}

func operateDoubleValueStack(stack *stack, operationGroup doubleOperandsOperationtionGroup) {
	preview := stack.peek(2)

	if err := helpers.OperatingOnEqualTypes(preview...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var operation func(superType string, n1, n2 orth_types.Operand) orth_types.Operand
	if helpers.IsFloat(preview[0].Operator) {
		operation = operationGroup.Float
	} else {
		operation = operationGroup.Integer
	}

	result := operation(preview[0].Operator.SymbolName, preview[0].Operator, preview[1].Operator)
	stack.rmv(2)
	stack.push(orth_types.Operation{
		Instruction: orth_types.InstructionPush,
		Operator:    result,
	})
}

// SimulateStack is an optional step that preceeds compilation, checking for errors, underflows, overflows
// and other things that a programmer like me would do without even thinking
func SimulateStack(program *orth_types.Program) {
	virtualMem := make([]orth_types.Operation, helpers.MEM_MAX_CAP)
	stack := stack{
		ptr:   -1,
		items: make([]orth_types.Operation, 1024),
	}

	for ip, operation := range program.Operations {
		switch operation.Instruction {
		case orth_types.InstructionPush:
			fallthrough
		case orth_types.InstructionPushStr:
			stack.push(operation)
		case orth_types.InstructionMult:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.MultplyIntegers,
				Float:   functions.MultplyFloats,
			})
		case orth_types.InstructionDiv:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.DivideIntegers,
				Float:   functions.DivideFloats,
			})
		case orth_types.InstructionSum:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.SumIntegers,
				Float:   functions.SumFloats,
			})
		case orth_types.InstructionMinus:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.SubIntegers,
				Float:   functions.SubFloats,
			})
		case orth_types.InstructionEqual:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.EqualInts,
				Float:   functions.EqualFloats,
			})
		case orth_types.InstructionLt:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.LowerThanInts,
				Float:   functions.LowerThanFloats,
			})
		case orth_types.InstructionGt:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.GreaterThanInts,
				Float:   functions.GreaterThanFloats,
			})
		case orth_types.InstructionNotEqual:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.EqualInts,
				Float:   functions.EqualFloats,
			})
		case orth_types.InstructionMod:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.ModIntegers,
				Float:   functions.ModFloats,
			})
		case orth_types.InstructionDup:
			preview := stack.peek(2)
			stack.push(preview[1], preview[0], preview[1])
		case orth_types.InstructionTwoDup:
			preview := stack.peek(2)
			stack.push(preview[1], preview[0], preview[1], preview[1])
		case orth_types.FunctionPutU64:
			fallthrough
		case orth_types.FunctionPutString:
			fallthrough
		case orth_types.InstructionDrop:
			stack.rmv(1)
		case orth_types.InstructionSwap:
			preview := stack.peek(2)
			stack.rmv(2)
			stack.push(preview[1], preview[0])
		case orth_types.InstructionMem:
			stack.push(program.Operations[ip])
		case orth_types.InstructionStore:
			preview := stack.peek(2)
			for _, item := range preview {
				if !helpers.IsInt(item.Operator) {
					fmt.Fprintf(os.Stderr, "cannot have type %q used for %q instruction", item.Operator.SymbolName, orth_types.InstructionToStr(orth_types.InstructionStore))
					os.Exit(1)
				}
			}
			// add to mem
			offset := helpers.ToInt(preview[0].Operator)
			memPtr := helpers.ToInt(preview[1].Operator)

			if offset+memPtr < 0 {
				fmt.Fprintf(os.Stderr, "cannot have a negative offset access for %q. Expected x >= 0 got '%d'", orth_types.InstructionToStr(orth_types.InstructionStore), offset+memPtr)
				os.Exit(1)
			}

			if offset > int(helpers.MEM_MAX_CAP) {
				fmt.Fprintf(os.Stderr, "%q offset larger than mem_max_cap: max allowed %d | actual %d", orth_types.InstructionToStr(orth_types.InstructionStore), helpers.MEM_MAX_CAP, offset)
				os.Exit(1)
			}
			// [1:] because index 0 is the offset itself
			addToMem(&virtualMem, offset, preview[1:]...)
			// remove from main stack
			stack.rmv(2)
		case orth_types.InstructionLoad:
			preview := stack.peek(1)
			stack.push(preview...)
			stack.rmv(1)
		case orth_types.InstructionLoadStay:
			preview := stack.peek(1)
			stack.push(preview...)
		case orth_types.InstructionCall:
			callingProcSchema, err := program.FindProc(operation)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			preview := stack.peek(len(callingProcSchema.InParamsAmount))
			for i, stackItem := range preview {
				if callingProcSchema.InParamsAmount[i].Operator.Operand != stackItem.Operator.SymbolName {
					fmt.Fprintf(os.Stderr, "Proc param required type %q but got %q", callingProcSchema.InParamsAmount[i].Operator.Operand, stackItem.Operator.SymbolName)
					os.Exit(1)
				}
			}
			// if param type checking went well, remove params from the main stack
			stack.rmv(len(callingProcSchema.InParamsAmount))
		case orth_types.InstructionEnd:
			procAddress, closingProc := operation.Addresses[orth_types.InstructionProc]
			if closingProc {
				callingProcSchema, err := program.FindProc(program.Operations[procAddress])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				preview := stack.peek(len(callingProcSchema.OutParamsAmount))
				for i, stackItem := range preview {
					if callingProcSchema.OutParamsAmount[i].Operator.Operand != stackItem.Operator.SymbolName {
						fmt.Fprintf(os.Stderr, "Proc return required type %q but got %q", callingProcSchema.OutParamsAmount[i].Operator.Operand, stackItem.Operator.SymbolName)
						os.Exit(1)
					}
				}
			}
		case orth_types.InstructionHold:
			var varAddress int
			if localAddr, ok := operation.Addresses[orth_types.InstructionVar]; ok {
				varAddress = localAddr
			} else if globalAddr, ok := operation.Addresses[orth_types.InstructionConst]; ok {
				varAddress = globalAddr
			}

			stack.push(orth_types.Operation{
				Instruction: orth_types.InstructionPush,
				Operator: orth_types.Operand{
					SymbolName: orth_types.StdAddress,
					Operand:    fmt.Sprint(varAddress),
				},
			})
		case orth_types.FunctionDumpMem:
			preview := stack.peek(2)
			for _, item := range preview {
				if !helpers.IsInt(item.Operator) {
					fmt.Fprintf(os.Stderr, "cannot have type %q used for %q instruction", item.Operator.SymbolName, orth_types.InstructionToStr(orth_types.FunctionDumpMem))
					os.Exit(1)
				}
			}
			stack.rmv(2)
			// future? idk
			// itemsToFetchAmount, err := strconv.Atoi(preview[0].Operator.Operand)
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "type %q has an invalid int value %q", preview[0].Operator.SymbolName, preview[0].Operator.Operand)
			// 	os.Exit(1)
			// }
			// memStartPosition, err := strconv.Atoi(preview[1].Operator.Operand)
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "type %q has an invalid address value %q", preview[1].Operator.SymbolName, preview[1].Operator.Operand)
			// 	os.Exit(1)
			// }
			// stack.push(virtualMem[memStartPosition:memStartPosition+itemsToFetchAmount]...)
		case orth_types.InstructionLShift:
			preview := stack.peek(2)
			stack.rmv(2)

			if !helpers.IsNumeric(preview[0].Operator) {
				fmt.Fprintf(os.Stderr, "cannot have type %q as LShift value", preview[0].Operator.SymbolName)
				os.Exit(1)
			}
			if !helpers.IsInt(preview[1].Operator) {
				fmt.Fprintf(os.Stderr, "cannot have type %q as LShift amount", preview[1].Operator.SymbolName)
				os.Exit(1)
			}

			var shiftResult orth_types.Operand
			if helpers.IsFloat(preview[0].Operator) {
				shiftResult = functions.LeftShiftFloat(orth_types.StdF32, preview[1].Operator, preview[0].Operator)
			} else {
				shiftResult = functions.LeftShiftInt(orth_types.StdF32, preview[1].Operator, preview[0].Operator)
			}
			stack.push(orth_types.Operation{
				Instruction: orth_types.InstructionPush,
				Context:     operation.Context,
				Operator:    shiftResult,
			})
		case orth_types.InstructionRShift:
			preview := stack.peek(2)
			stack.rmv(2)

			if !helpers.IsNumeric(preview[0].Operator) {
				fmt.Fprintf(os.Stderr, "cannot have type %q as %q value", preview[0].Operator.SymbolName, orth_types.InstructionToStr(orth_types.InstructionRShift))
				os.Exit(1)
			}
			if !helpers.IsInt(preview[1].Operator) {
				fmt.Fprintf(os.Stderr, "cannot have type %q for %q amount", preview[1].Operator.SymbolName, orth_types.InstructionToStr(orth_types.InstructionRShift))
				os.Exit(1)
			}

			var shiftResult orth_types.Operand
			if helpers.IsFloat(preview[0].Operator) {
				shiftResult = functions.RightShiftFloat(orth_types.StdF32, preview[1].Operator, preview[0].Operator)
			} else {
				shiftResult = functions.RightShiftInt(orth_types.StdF32, preview[1].Operator, preview[0].Operator)
			}
			stack.push(orth_types.Operation{
				Instruction: orth_types.InstructionPush,
				Context:     operation.Context,
				Operator:    shiftResult,
			})
		case orth_types.InstructionLAnd:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.BitwiseAnd,
				Float:   functions.BitwiseAnd,
			})
		case orth_types.InstructionLOr:
			operateDoubleValueStack(&stack, doubleOperandsOperationtionGroup{
				Integer: functions.BitwiseOr,
				Float:   functions.BitwiseOr,
			})
		case orth_types.InstructionOver:
			preview := stack.peek(2)
			stack.push(preview[0])
		case orth_types.InstructionExit:
			preview := stack.peek(1)
			if !helpers.IsInt(preview[0].Operator) {
				fmt.Fprintln(os.Stderr, "'exit' only accepts integer values")
				os.Exit(1)
			}
			exitCode := helpers.ToInt(preview[0].Operator)
			os.Exit(exitCode)
		case orth_types.InstructionDeref:
			preview := stack.peek(1)
			addr, ok := helpers.ToAddress(preview[0].Operator)

			if !ok {
				fmt.Fprintf(os.Stderr, "cannot non addressable value for instruction %q\n", orth_types.InstructionToStr(orth_types.InstructionDeref))
				os.Exit(1)
			}

			stack.push(program.Operations[addr])
		}
	}
}
