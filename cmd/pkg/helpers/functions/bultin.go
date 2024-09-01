package functions

import (
	"fmt"
	"log"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers"
	orth_types "orth/cmd/pkg/types"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var Functions map[string]func(
	stack *[]orth_types.Operand,
	mem *[]orth_types.Operand,
	vars map[string]orth_types.Operand)

func init() {
	Functions = make(map[string]func(*[]orth_types.Operand, *[]orth_types.Operand, map[string]orth_types.Operand))

	Functions["to_string"] = func(stack *[]orth_types.Operand, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		o1 := helpers.StackPop(stack)
		res := ToString(o1)
		*stack = append(*stack, res)
	}

	Functions["size_of"] = func(stack *[]orth_types.Operand, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		o1 := helpers.StackPop(stack)
		switch o1.SymbolName {
		case orth_types.StdSTR:
			*stack = append(*stack, orth_types.Operand{
				SymbolName: orth_types.StdI32,
				Operand:    fmt.Sprint(len(o1.Operand)),
			})
		default:
			panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, o1.SymbolName, "Functions[length_of]"))
		}
	}

	Functions["length_of"] = func(stack *[]orth_types.Operand, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		o1 := helpers.StackPop(stack)
		switch o1.SymbolName {
		case orth_types.StdSTR:
			*stack = append(*stack, orth_types.Operand{
				SymbolName: orth_types.StdI32,
				Operand:    fmt.Sprint(utf8.RuneCountInString(o1.Operand)),
			})
		default:
			panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, o1.SymbolName, "Functions[length_of]"))
		}
	}

	Functions["make_array"] = func(stack *[]orth_types.Operand, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		if len(*stack) < 2 {
			panic(orth_debug.StackUnderFlow)
		}

		capacity := helpers.ToInt(helpers.StackPop(stack))
		typ := helpers.StackPop(stack)

		insertCollectionToMem(mem, stack, typ, capacity)
	}

	Functions["free_var"] = func(stack *[]orth_types.Operand, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		vName := helpers.StackPop(stack)
		_, ok := vars[vName.Operand]

		if !ok {
			panic(fmt.Errorf(orth_debug.VariableUndefined, vName.Operand))
		}

		if vName.SymbolName != orth_types.StdConst {
			panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, vName.SymbolName, "free_var"))
		}
		delete(vars, vName.Operand)
	}

	Functions["dump_mem"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		DumpMem(*stack, *mem)
	}

	Functions["dump_stack"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		DumpStack(*stack)
	}

	Functions["dump_vars"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		DumpVars(vars)
	}

	Functions["exit"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		code := helpers.StackPop(stack)
		if helpers.IsInt(code.SymbolName) {
			log.Println("Progam exited with code:", code.Operand)
			os.Exit(helpers.ToInt(code))
		}
		panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, code.SymbolName, "exit"))
	}

	Functions["fill"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		value := helpers.StackPop(stack)
		rangeable := helpers.StackPop(stack)

		if rangeable.SymbolName != orth_types.RNGABL {
			panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, rangeable.SymbolName, "fill"))
		}
		start, end := DissectRangeAsInt(rangeable)
		helpers.SameBaseType((*mem)[start], value)

		for i := start; i <= end; i++ {
			(*mem)[i] = value
		}

		*stack = append(*stack, rangeable)
	}

	Functions["index"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		rangeable := helpers.StackPop(stack)
		if rangeable.SymbolName != orth_types.RNGABL {
			panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, rangeable.SymbolName, "index"))
		}

		nums := strings.Split(rangeable.Operand, "|")
		index, _ := strconv.Atoi(nums[0])

		*stack = append(*stack, orth_types.Operand{
			SymbolName: orth_types.ADDR,
			Operand:    fmt.Sprint(index),
		})
	}

	Functions["grab_at"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		o1 := helpers.StackPop(stack)
		o2 := helpers.StackPop(stack)

		if !helpers.IsInt(o2.SymbolName) {
			panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, o2.SymbolName, "grab_at"))
		}
		o3 := (*stack)[helpers.ToInt(o2)]

		vars[o1.Operand] = o3
	}

	Functions["grab_last"] = func(stack, mem *[]orth_types.Operand, vars map[string]orth_types.Operand) {
		o1 := helpers.StackPop(stack)
		o2 := helpers.StackPop(stack)

		vars[o1.Operand] = o2
	}
}

func DumpVars(vars map[string]orth_types.Operand) {
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

func DumpStack(stack []orth_types.Operand) {
	if len(stack) == 0 {
		fmt.Println("***********************************************************")
		fmt.Println("STACK IS EMPTY")
		fmt.Println("***********************************************************")
		return
	}

	fmt.Println("STACK:")
	for i := len(stack); i > 0; i-- {
		fmt.Printf("Position: %d\t Type: %q\t Value: %#v\n", i-1, stack[i-1].SymbolName, stack[i-1].Operand)
	}
}

func DumpMem(stack, mem []orth_types.Operand) {
	opTo := helpers.StackPop(&stack)
	opFrom := helpers.StackPop(&stack)

	if !helpers.IsInt(opTo.SymbolName) && !helpers.IsInt(opFrom.SymbolName) {
		msg := fmt.Sprintf(orth_debug.InvalidTypeForInstruction+"\n", opTo.Operand, "dump_mem")
		msg += fmt.Sprintf(orth_debug.InvalidTypeForInstruction, opFrom.Operand, "dump_mem")
		panic(msg)
	}

	to, from := helpers.ToInt(opTo), helpers.ToInt(opFrom)

	slicedMem := mem[from:to]
	for i, op := range slicedMem {
		fmt.Printf("%d: %v\n", i, op)
	}
}

func insertCollectionToMem(originalMem, stack *[]orth_types.Operand, typ orth_types.Operand, capacity int) {
	// why go?
	memCopy := make([]orth_types.Operand, len(*originalMem), cap(*originalMem))
	copy(memCopy, *originalMem)

	var start int
	var foundPlace bool

	// look for a place where all members can fit by direct index
	for i := range memCopy {
		fitInAmount := 0
		slice := memCopy[i:(i + capacity)]

		for _, xx := range slice {
			if xx.Operand == "" && xx.SymbolName == "" {
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
			memCopy[i] = orth_types.Operand{
				SymbolName: typ.Operand,
			}
		}
		*originalMem = memCopy
		*stack = append(*stack, orth_types.Operand{
			SymbolName: orth_types.RNGABL,
			Operand:    fmt.Sprintf("%d|%d", start, start+capacity-1),
		})
	}
}
