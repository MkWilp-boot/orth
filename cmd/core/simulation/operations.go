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
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Push,
				Operand: orthtypes.Operand{
					VarType: v.Content.Content,
					Operand: preProgram[i+1].Content.Content,
				},
				RefBlock: -1,
			})
		case orthtypes.PrimitiveF32:
			fallthrough
		case orthtypes.PrimitiveF64:
			preProgram[i+1].Content.ValidPos = true
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Push,
				Operand: orthtypes.Operand{
					VarType: v.Content.Content,
					Operand: preProgram[i+1].Content.Content,
				},
				RefBlock: -1,
			})
		case orthtypes.PrimitiveSTR:
			preProgram[i+1].Content.ValidPos = true
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Push,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveSTR,
					Operand: preProgram[i+1].Content.Content[1 : len(preProgram[i+1].Content.Content)-1],
				},
				RefBlock: -1,
			})
		case orthtypes.PrimitiveBOOL:
			preProgram[i+1].Content.ValidPos = true
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Push,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: preProgram[i+1].Content.Content,
				},
				RefBlock: -1,
			})
		case "+":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Sum,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "-":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Minus,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "*":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Mult,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "/":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Div,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "dump":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Dump,
				Operand: orthtypes.Operand{
					VarType: orthtypes.VOID,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "=":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Equal,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "<>":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.NotEqual,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "<":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Lt,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case ">":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Gt,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "if":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.If,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "else":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Else,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "end":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.End,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveBOOL,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "print":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Print,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveRNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "dup":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Dup,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveRNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "while":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.While,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveRNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "do":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Do,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "drop":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Drop,
				Operand: orthtypes.Operand{
					VarType: orthtypes.VOID,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "swap":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Swap,
				Operand: orthtypes.Operand{
					VarType: orthtypes.VOID,
					Operand: "",
				},
				RefBlock: -1,
			})
		case "%":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Mod,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case orthtypes.PrimitiveMem:
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Mem,
				Operand: orthtypes.Operand{
					VarType: orthtypes.PrimitiveMem,
					Operand: orthtypes.PrimitiveMem,
				},
				RefBlock: -1,
			})
		case ".":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Store,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		case ",":
			program.Operations = append(program.Operations, orthtypes.Operation{
				Instruction: orthtypes.Load,
				Operand: orthtypes.Operand{
					VarType: orthtypes.RNT,
					Operand: "",
				},
				RefBlock: -1,
			})
		default:
			if !v.Content.ValidPos {
				panic(fmt.Errorf("Unknow token %q at line: %d colum: %d\n", v.Content.Content, v.Index, v.Content.Index))
			}
		}
	}
	return program
}
