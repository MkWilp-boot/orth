package functions

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"t/cmd/core/debug"
	"t/cmd/pkg/helpers"
	orthtypes "t/cmd/pkg/types"
	"unicode/utf8"
)

var Functions map[string]func(
	stack *[]orthtypes.Operand,
	mem *[]orthtypes.Operand,
	vars map[string]orthtypes.Operand)

func init() {
	Functions = make(map[string]func(*[]orthtypes.Operand, *[]orthtypes.Operand, map[string]orthtypes.Operand))

	Functions["to_string"] = func(stack *[]orthtypes.Operand, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(stack)
		res := ToString(o1)
		*stack = append(*stack, res)
	}

	Functions["size_of"] = func(stack *[]orthtypes.Operand, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(stack)
		switch o1.VarType {
		case orthtypes.PrimitiveSTR:
			*stack = append(*stack, orthtypes.Operand{
				VarType: orthtypes.PrimitiveI32,
				Operand: fmt.Sprint(len(o1.Operand)),
			})
		default:
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, o1.VarType, "Functions[length_of]"))
		}
	}

	Functions["length_of"] = func(stack *[]orthtypes.Operand, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(stack)
		switch o1.VarType {
		case orthtypes.PrimitiveSTR:
			*stack = append(*stack, orthtypes.Operand{
				VarType: orthtypes.PrimitiveI32,
				Operand: fmt.Sprint(utf8.RuneCountInString(o1.Operand)),
			})
		default:
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, o1.VarType, "Functions[length_of]"))
		}
	}

	Functions["make_array"] = func(stack *[]orthtypes.Operand, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		if len(*stack) < 2 {
			panic(debug.StackUnderFlow)
		}

		capacity := helpers.ToInt(helpers.StackPop(stack))
		typ := helpers.StackPop(stack)

		insertCollectionToMem(mem, stack, typ, capacity)
	}

	Functions["free_var"] = func(stack *[]orthtypes.Operand, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		vName := helpers.StackPop(stack)
		_, ok := vars[vName.Operand]

		if !ok {
			panic(fmt.Errorf(debug.VariableUndefined, vName.Operand))
		}

		if vName.VarType != orthtypes.PrimitiveVar {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, vName.VarType, "free_var"))
		}
		delete(vars, vName.Operand)
	}

	Functions["dump_stack"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		DumpStack(*stack)
	}

	Functions["write"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		numBytes := helpers.StackPop(&*stack)
		content := helpers.StackPop(&*stack)
		absPos := helpers.StackPop(&*stack)

		if absPos.VarType != orthtypes.ADDR {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, absPos.VarType, "write"))
		}
		if !helpers.IsInt(numBytes.VarType) {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, numBytes.VarType, "write"))
		}

		at := helpers.ToInt(absPos)

		reader := bytes.NewReader([]byte(content.Operand))

		bs := make([]byte, helpers.ToInt(numBytes))
		_, err := reader.Read(bs)
		if err != nil {
			log.Fatal(err)
		}

		(*mem)[at] = orthtypes.Operand{
			VarType: content.VarType,
			Operand: string(bs),
		}
		*stack = append(*stack, orthtypes.Operand{
			VarType: orthtypes.PrimitiveInt,
			Operand: fmt.Sprint(len(bs)),
		})
	}

	Functions["at"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(&*stack)
		if !helpers.IsInt(o1.VarType) {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, o1.VarType, "read"))
		}

		at := helpers.ToInt(o1)
		var retValue orthtypes.Operand

		source := helpers.StackPop(&*stack)

		if source.VarType != orthtypes.RNGABL {
			retValue = orthtypes.Operand{
				VarType: orthtypes.ADDR,
				Operand: fmt.Sprint(at),
			}
		} else {
			points := strings.Split(source.Operand, "|")
			start, _ := strconv.Atoi(points[0])
			end, _ := strconv.Atoi(points[1])

			if start+at >= end {
				panic(fmt.Errorf(debug.IndexOutOfBounds, at, start, end))
			}
			retValue = orthtypes.Operand{
				VarType: orthtypes.ADDR,
				Operand: fmt.Sprint(start + at),
			}
		}

		*stack = append(*stack, retValue)
	}

	Functions["read"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(&*stack)
		if !helpers.IsInt(o1.VarType) {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, o1.VarType, "read"))
		}

		at := helpers.ToInt(o1)
		var retValue orthtypes.Operand

		source := helpers.StackPop(&*stack)

		if source.VarType != orthtypes.RNGABL {
			retValue = (*mem)[at]
		} else {
			points := strings.Split(source.Operand, "|")
			start, _ := strconv.Atoi(points[0])
			end, _ := strconv.Atoi(points[1])

			if start+at >= end {
				panic(fmt.Errorf(debug.IndexOutOfBounds, at, start, end))
			}
			retValue = (*mem)[start+at]
		}

		*stack = append(*stack, retValue)
	}
}

func DumpStack(stack []orthtypes.Operand) {
	if len(stack) == 0 {
		fmt.Println("***********************************************************")
		fmt.Println("STACK IS EMPTY")
		fmt.Println("***********************************************************")
		return
	}

	fmt.Println("STACK:")
	for i := len(stack); i > 0; i-- {
		fmt.Printf("Position: %d\t Type: %q\t Value: %#v\n", i, stack[i-1].VarType, stack[i-1].Operand)
	}
}

func insertCollectionToMem(originalMem, stack *[]orthtypes.Operand, typ orthtypes.Operand, capacity int) {
	// why go?
	memCopy := make([]orthtypes.Operand, len(*originalMem), cap(*originalMem))
	copy(memCopy, *originalMem)

	var start int
	var foundPlace bool

	// look for a place where all members can fit by direct index
	for i := range memCopy {
		fitInAmount := 0
		slice := memCopy[i:(i + capacity)]

		for _, xx := range slice {
			if xx.Operand == "" && xx.VarType == "" {
				fitInAmount++
			}
		}

		if fitInAmount == capacity {
			foundPlace = true
			start = i
			break
		}
	}

	if foundPlace {
		for i := start + 1; i < capacity; i++ {
			memCopy[i] = orthtypes.Operand{
				VarType: typ.Operand,
			}
		}
		*originalMem = memCopy
		*stack = append(*stack, orthtypes.Operand{
			VarType: orthtypes.RNGABL,
			Operand: fmt.Sprintf("%d|%d", start, start+capacity-1),
		})
	}
}
