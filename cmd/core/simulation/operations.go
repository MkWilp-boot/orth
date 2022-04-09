package simulation

import (
	"fmt"
	orthtypes "t/cmd/pkg/types"
)

func PopLast(root *[]int) int {
	stack := *root
	ret := stack[len(stack)-1]
	*root = stack[:len(stack)-1]
	return ret
}

func CrossReferenceBlocks(program orthtypes.Program) orthtypes.Program {
	stack := make([]int, 0)

	// program.Operations[ip] is actually the current position in loop
	for ip, v := range program.Operations {
		switch v.Instruction {
		case orthtypes.If:
			stack = append(stack, ip)
		case orthtypes.Else:
			blockIp := PopLast(&stack)

			if program.Operations[blockIp].Instruction != orthtypes.If {
				panic("Invalid Else clause")
			}

			program.Operations[blockIp].RefBlock = ip + 1
			stack = append(stack, ip)
		case orthtypes.End:
			blockIp := PopLast(&stack)
			if program.Operations[blockIp].Instruction == orthtypes.If ||
				program.Operations[blockIp].Instruction == orthtypes.Else {

				program.Operations[blockIp].RefBlock = ip
				program.Operations[ip].RefBlock = ip + 1 // end block
			} else if program.Operations[blockIp].Instruction == orthtypes.Do {
				if program.Operations[blockIp].RefBlock == -1 {
					panic("Not enought arguments for a cross-refernce block operation")
				}
				program.Operations[ip].RefBlock = program.Operations[blockIp].RefBlock
				program.Operations[blockIp].RefBlock = ip + 1
			} else {
				panic("End block can only close [if | else | do] blocks")
			}
		case orthtypes.Do:
			blockIp := PopLast(&stack)
			program.Operations[ip].RefBlock = blockIp
			stack = append(stack, ip)
		case orthtypes.While:
			stack = append(stack, ip)
		}
	}
	return program
}

func ParseTokenAsOperation(preProgram []orthtypes.StringEnum) orthtypes.Program {
	program := orthtypes.Program{}

	for i, v := range preProgram {
		switch v.Content.Content {
		case orthtypes.PrimitiveInt:
			fallthrough
		case orthtypes.PrimitiveI8:
			fallthrough
		case orthtypes.PrimitiveI16:
			fallthrough
		case orthtypes.PrimitiveI32:
			fallthrough
		case orthtypes.PrimitiveI64:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(v.Content.Content, preProgram[i+1].Content.Content, orthtypes.Push, -1)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveF32:
			fallthrough
		case orthtypes.PrimitiveF64:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(v.Content.Content, preProgram[i+1].Content.Content, orthtypes.Push, -1)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveSTR:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Content[1:len(preProgram[i+1].Content.Content)-1], orthtypes.Push, -1)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveBOOL:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(orthtypes.PrimitiveBOOL, preProgram[i+1].Content.Content, orthtypes.Push, -1)
			program.Operations = append(program.Operations, ins)
		case "+":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Sum, -1)
			program.Operations = append(program.Operations, ins)
		case "-":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Minus, -1)
			program.Operations = append(program.Operations, ins)
		case "*":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Mult, -1)
			program.Operations = append(program.Operations, ins)
		case "/":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Div, -1)
			program.Operations = append(program.Operations, ins)
		case "dump":
			ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.Dump, -1)
			program.Operations = append(program.Operations, ins)
		case "=":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Equal, -1)
			program.Operations = append(program.Operations, ins)
		case "<>":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.NotEqual, -1)
			program.Operations = append(program.Operations, ins)
		case "<":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Lt, -1)
			program.Operations = append(program.Operations, ins)
		case ">":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Gt, -1)
			program.Operations = append(program.Operations, ins)
		case "if":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.If, -1)
			program.Operations = append(program.Operations, ins)
		case "else":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Else, -1)
			program.Operations = append(program.Operations, ins)
		case "end":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.End, -1)
			program.Operations = append(program.Operations, ins)
		case "print":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Print, -1)
			program.Operations = append(program.Operations, ins)
		case "dup":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Dup, -1)
			program.Operations = append(program.Operations, ins)
		case "while":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.While, -1)
			program.Operations = append(program.Operations, ins)
		case "do":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Do, -1)
			program.Operations = append(program.Operations, ins)
		case "drop":
			ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.Drop, -1)
			program.Operations = append(program.Operations, ins)
		case "swap":
			ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.Swap, -1)
			program.Operations = append(program.Operations, ins)
		case "%":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Mod, -1)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveMem:
			ins := parseToken(orthtypes.PrimitiveMem, orthtypes.PrimitiveMem, orthtypes.Mem, -1)
			program.Operations = append(program.Operations, ins)
		case ".":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Store, -1)
			program.Operations = append(program.Operations, ins)
		case ",":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Load, -1)
			program.Operations = append(program.Operations, ins)
		default:
			if !v.Content.ValidPos {
				panic(fmt.Errorf("Unknow token %q at line: %d colum: %d\n", v.Content.Content, v.Index, v.Content.Index))
			}
		}
	}
	return program
}

func parseToken(varType, operand string, op, refBlock int) orthtypes.Operation {
	return orthtypes.Operation{
		Instruction: op,
		Operand: orthtypes.Operand{
			VarType: varType,
			Operand: operand,
		},
		RefBlock: refBlock,
	}
}
