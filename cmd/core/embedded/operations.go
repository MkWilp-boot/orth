package embedded

import (
	"fmt"
	"regexp"
	"t/cmd/core/debug"
	"t/cmd/pkg/helpers/functions"
	orthtypes "t/cmd/pkg/types"
)

func PopLast(root *[]int) int {
	stack := *root
	ret := stack[len(stack)-1]
	*root = stack[:len(stack)-1]
	return ret
}

// CrossReferenceBlocks loops over a program and define all inter references
// needed for execution. Ex: if-else-do blocks
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

// ParseTokenAsOperation parses an slice of pre-instructions into a runnable program
func ParseTokenAsOperation(preProgram []orthtypes.StringEnum) orthtypes.Program {
	program := orthtypes.Program{}

	for i, v := range preProgram {
		if v.Content.ValidPos {
			continue
		}
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
			ins := parseToken(v.Content.Content, preProgram[i+1].Content.Content, orthtypes.Push)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveF32:
			fallthrough
		case orthtypes.PrimitiveF64:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(v.Content.Content, preProgram[i+1].Content.Content, orthtypes.Push)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveSTR:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Content[1:len(preProgram[i+1].Content.Content)-1], orthtypes.Push)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveBOOL:
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(orthtypes.PrimitiveBOOL, preProgram[i+1].Content.Content, orthtypes.Push)
			program.Operations = append(program.Operations, ins)
		case "+":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Sum)
			program.Operations = append(program.Operations, ins)
		case "-":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Minus)
			program.Operations = append(program.Operations, ins)
		case "*":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Mult)
			program.Operations = append(program.Operations, ins)
		case "/":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Div)
			program.Operations = append(program.Operations, ins)
		case "dump":
			ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.Dump)
			program.Operations = append(program.Operations, ins)
		case "==":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Equal)
			program.Operations = append(program.Operations, ins)
		case "<>":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.NotEqual)
			program.Operations = append(program.Operations, ins)
		case "<":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Lt)
			program.Operations = append(program.Operations, ins)
		case ">":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Gt)
			program.Operations = append(program.Operations, ins)
		case "if":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.If)
			program.Operations = append(program.Operations, ins)
		case "else":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.Else)
			program.Operations = append(program.Operations, ins)
		case "end":
			ins := parseToken(orthtypes.PrimitiveBOOL, "", orthtypes.End)
			program.Operations = append(program.Operations, ins)
		case "print":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Print)
			program.Operations = append(program.Operations, ins)
		case "dup":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Dup)
			program.Operations = append(program.Operations, ins)
		case "while":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.While)
			program.Operations = append(program.Operations, ins)
		case "do":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Do)
			program.Operations = append(program.Operations, ins)
		case "drop":
			ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.Drop)
			program.Operations = append(program.Operations, ins)
		case "swap":
			ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.Swap)
			program.Operations = append(program.Operations, ins)
		case "%":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Mod)
			program.Operations = append(program.Operations, ins)
		case orthtypes.PrimitiveMem:
			ins := parseToken(orthtypes.PrimitiveMem, orthtypes.PrimitiveMem, orthtypes.Mem)
			program.Operations = append(program.Operations, ins)
		case ".":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Store)
			program.Operations = append(program.Operations, ins)
		case ",":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Load)
			program.Operations = append(program.Operations, ins)
		case "call":
			_, ok := functions.Functions[preProgram[i+1].Content.Content]
			if !ok {
				panic(fmt.Errorf(debug.UndefinedToken, preProgram[i+1].Content.Content))
			}
			preProgram[i+1].Content.ValidPos = true
			ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Content, orthtypes.Call)
			program.Operations = append(program.Operations, ins)
		case ",!":
			ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.LoadStay)
			program.Operations = append(program.Operations, ins)
		case "type":
			preProgram[i+1].Content.ValidPos = true

			ins := parseToken(orthtypes.PrimitiveType, preProgram[i+1].Content.Content, orthtypes.Push)

			program.Operations = append(program.Operations, ins)
		case "var":
			re := regexp.MustCompile(`[^\w]`)

			// check name
			if re.Match([]byte(preProgram[i+1].Content.Content)) {
				panic("var has invalid characters in it's composition")
			}
			// check if has a value
			if preProgram[i+2].Content.Content != "=" {
				switch {
				// used as a func param
				case preProgram[i+2].Content.Content == "call":
					preProgram[i+1].Content.ValidPos = true

					ins := parseToken(orthtypes.PrimitiveVar, preProgram[i+1].Content.Content, orthtypes.Push)
					program.Operations = append(program.Operations, ins)
					continue

				default:
					panic("var must be initialized with `=` sign")
				}
			}

			for x := 1; x < 5; x++ {
				preProgram[i+x].Content.ValidPos = true
			}

			vName := preProgram[i+1].Content.Content
			vType := preProgram[i+3].Content.Content

			var vValue string

			if vType == orthtypes.PrimitiveSTR {
				vValue = preProgram[i+4].Content.Content[1 : len(preProgram[i+4].Content.Content)-1]
			} else {
				vValue = preProgram[i+4].Content.Content
			}

			ins := parseToken(vType, vValue, orthtypes.Push)
			program.Operations = append(program.Operations, ins)

			ins = parseToken(orthtypes.PrimitiveVar, vName, orthtypes.Var)
			program.Operations = append(program.Operations, ins)
		case "hold":
			preProgram[i+1].Content.ValidPos = true
			vName := preProgram[i+1].Content.Content

			ins := parseToken(orthtypes.PrimitiveHold, vName, orthtypes.Hold)
			program.Operations = append(program.Operations, ins)
		default:
			if !v.Content.ValidPos {
				panic(fmt.Errorf("unknow token %q at line: %d colum: %d", v.Content.Content, v.Index, v.Content.Index))
			}
		}
	}
	return program
}

// parseToken parses a single token into a instruction
func parseToken(varType, operand string, op int) orthtypes.Operation {
	return orthtypes.Operation{
		Instruction: op,
		Operand: orthtypes.Operand{
			VarType: varType,
			Operand: operand,
		},
		RefBlock: -1,
	}
}
