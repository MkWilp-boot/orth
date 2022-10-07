package functions

import (
	"fmt"
	"log"
	"orth/cmd/core/debug"
	"orth/cmd/pkg/helpers"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"strconv"
	"strings"
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

	Functions["dump_mem"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		DumpMem(*stack, *mem)
	}

	Functions["dump_stack"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		DumpStack(*stack)
	}

	Functions["dump_vars"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		DumpVars(vars)
	}

	Functions["exit"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		code := helpers.StackPop(&*stack)
		if helpers.IsInt(code.VarType) {
			log.Println("Progam exited with code:", code.Operand)
			os.Exit(helpers.ToInt(code))
		}
		panic(fmt.Errorf(debug.InvalidTypeForInstruction, code.VarType, "exit"))
	}

	Functions["fill"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		value := helpers.StackPop(&*stack)
		rangeable := helpers.StackPop(&*stack)

		if rangeable.VarType != orthtypes.RNGABL {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, rangeable.VarType, "fill"))
		}
		start, end := DissectRangeAsInt(rangeable)
		helpers.SameBaseType((*mem)[start], value)

		for i := start; i <= end; i++ {
			(*mem)[i] = value
		}

		*stack = append(*stack, rangeable)
	}

	Functions["index"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		rangeable := helpers.StackPop(&*stack)
		if rangeable.VarType != orthtypes.RNGABL {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, rangeable.VarType, "index"))
		}

		nums := strings.Split(rangeable.Operand, "|")
		index, _ := strconv.Atoi(nums[0])

		*stack = append(*stack, orthtypes.Operand{
			VarType: orthtypes.ADDR,
			Operand: fmt.Sprint(index),
		})
	}

	Functions["grab_at"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(&*stack)
		o2 := helpers.StackPop(&*stack)

		if !helpers.IsInt(o2.VarType) {
			panic(fmt.Errorf(debug.InvalidTypeForInstruction, o2.VarType, "grab_at"))
		}
		o3 := (*stack)[helpers.ToInt(o2)]

		vars[o1.Operand] = o3
	}

	Functions["grab_last"] = func(stack, mem *[]orthtypes.Operand, vars map[string]orthtypes.Operand) {
		o1 := helpers.StackPop(&*stack)
		o2 := helpers.StackPop(&*stack)

		vars[o1.Operand] = o2
	}
}

func DumpVars(vars map[string]orthtypes.Operand) {
	if len(vars) == 0 {
		fmt.Println("***********************************************************")
		fmt.Println("VARS IS EMPTY")
		fmt.Println("***********************************************************")
		return
	}

	fmt.Println("VARS:")
	for name, value := range vars {
		fmt.Printf("var name: %q\t var value: %v\n", name, value)
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
		fmt.Printf("Position: %d\t Type: %q\t Value: %#v\n", i-1, stack[i-1].VarType, stack[i-1].Operand)
	}
}

func DumpMem(stack, mem []orthtypes.Operand) {
	opTo := helpers.StackPop(&stack)
	opFrom := helpers.StackPop(&stack)

	if !helpers.IsInt(opTo.VarType) && !helpers.IsInt(opFrom.VarType) {
		msg := fmt.Sprintf(debug.InvalidTypeForInstruction+"\n", opTo.Operand, "dump_mem")
		msg += fmt.Sprintf(debug.InvalidTypeForInstruction, opFrom.Operand, "dump_mem")
		panic(msg)
	}

	to, from := helpers.ToInt(opTo), helpers.ToInt(opFrom)

	slicedMem := mem[from:to]
	for i, op := range slicedMem {
		fmt.Printf("%d: %v\n", i, op)
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
		for i := start; i < capacity; i++ {
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
