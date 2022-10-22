package embedded

import (
	"fmt"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers/functions"
	orthtypes "orth/cmd/pkg/types"
	"regexp"
)

const globalScope = "Global"

func PopLast[T comparable](root *[]T) T {
	stack := *root
	ret := stack[len(stack)-1]
	*root = stack[:len(stack)-1]
	return ret
}

// CrossReferenceBlocks loops over a program and define all inter references
// needed for execution. Ex: if-else-do blocks
func CrossReferenceBlocks(program orthtypes.Program, crossResult chan<- orthtypes.Pair[orthtypes.Program, error]) {

	stack := make([]orthtypes.Pair[int, orthtypes.Operation], 0)
	orthVars := make(map[string]map[string]int) // context -> var_name -> vars_declared
	context := globalScope

	// program.Operations[ip] is actually the current position in loop
	for ip, v := range program.Operations {
		pair := orthtypes.Pair[int, orthtypes.Operation]{
			VarName:  ip,
			VarValue: v,
		}
		switch v.Instruction {
		case orthtypes.Var:
			_, ok := orthVars[context]
			if !ok {
				orthVars[context] = make(map[string]int)
			}
			orthVars[context][v.Operand.Operand]++

			if orthVars[context][v.Operand.Operand] != 1 || (context != globalScope && orthVars[globalScope][v.Operand.Operand] == 1) {
				crossResult <- orthtypes.Pair[orthtypes.Program, error]{
					VarName:  orthtypes.Program{},
					VarValue: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "Var", v.Operand.Operand, context),
				}
				close(crossResult)
				return
			}
			program.Operations[ip].Context = context
		case orthtypes.Hold:
			holds := program.Filter(func(op orthtypes.Operation) bool {
				return op == v
			})
			vD := program.Filter(func(op orthtypes.Operation) bool {
				return op.Operand.VarType == orthtypes.PrimitiveVar && op.Operand.Operand == holds[0].VarValue.Operand.Operand
			})

			program.Operations[holds[0].VarName].Context = vD[0].VarValue.Context
		case orthtypes.If:
			fallthrough
		case orthtypes.Proc:
			context = v.Operand.Operand
			fallthrough
		case orthtypes.While:
			stack = append(stack, pair)
		case orthtypes.Else:
			blockIp := PopLast(&stack)

			if program.Operations[blockIp.VarName].Instruction != orthtypes.If {
				panic("Invalid Else clause")
			}

			program.Operations[blockIp.VarName].RefBlock = ip + 1
			stack = append(stack, pair)
		case orthtypes.End:
			blockIp := PopLast(&stack)
			switch {
			case program.Operations[blockIp.VarName].Instruction == orthtypes.If:
				fallthrough
			case program.Operations[blockIp.VarName].Instruction == orthtypes.Else:
				program.Operations[blockIp.VarName].RefBlock = ip
				program.Operations[ip].RefBlock = ip + 1 // end block
			case program.Operations[blockIp.VarName].Instruction == orthtypes.In:
				context = globalScope
				fallthrough
			case program.Operations[blockIp.VarName].Instruction == orthtypes.Do:
				if program.Operations[blockIp.VarName].RefBlock == -1 {
					panic("Not enought arguments for a cross-refernce block operation")
				}
				program.Operations[ip].RefBlock = program.Operations[blockIp.VarName].RefBlock
				program.Operations[blockIp.VarName].RefBlock = ip + 1
			default:
				panic("End block can only close [if | else | do | proc in] blocks")
			}
		case orthtypes.In:
			fallthrough
		case orthtypes.Do:
			blockIp := PopLast(&stack)
			program.Operations[ip].RefBlock = blockIp.VarName
			stack = append(stack, pair)
		}
	}
	crossResult <- orthtypes.Pair[orthtypes.Program, error]{
		VarName:  program,
		VarValue: nil,
	}
	close(crossResult)
}

// ParseTokenAsOperation parses an slice of pre-instructions into a runnable program
func ParseTokenAsOperation(tokenFiles []orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]], parseTokenResult chan<- orthtypes.Pair[orthtypes.Program, error]) {
	program := orthtypes.Program{
		Id: len(tokenFiles),
	}
	var context string
	procNames := make(map[string]int)

	for fIndex, file := range tokenFiles {
		for i, v := range *file.CodeBlock.Slice {
			preProgram := (*tokenFiles[fIndex].CodeBlock.Slice)

			if v.Content.ValidPos {
				continue
			}
			switch v.Content.Content {
			case orthtypes.ADDR:
				fallthrough
			case orthtypes.PrimitiveRNT:
				fallthrough
			case orthtypes.PrimitiveInt:
				fallthrough
			case orthtypes.PrimitiveI8:
				fallthrough
			case orthtypes.PrimitiveI16:
				fallthrough
			case orthtypes.PrimitiveI32:
				fallthrough
			case orthtypes.PrimitiveI64:
				fallthrough
			case orthtypes.PrimitiveF32:
				fallthrough
			case orthtypes.PrimitiveF64:
				fallthrough
			case orthtypes.PrimitiveBOOL:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(v.Content.Content, preProgram[i+1].Content.Content, orthtypes.Push)
				program.Operations = append(program.Operations, ins)
			case orthtypes.PrimitiveSTR:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Content[1:len(preProgram[i+1].Content.Content)-1], orthtypes.Push)
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
			case "put_u64":
				ins := parseToken(orthtypes.PrimitiveVOID, "", orthtypes.PutU64)
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
				ins := parseToken(orthtypes.PrimitiveEND, "", orthtypes.End)
				program.Operations = append(program.Operations, ins)
			case "put_string":
				ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.PutString)
				program.Operations = append(program.Operations, ins)
			case "dup":
				ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.Dup)
				program.Operations = append(program.Operations, ins)
			case "while":
				ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.While)
				program.Operations = append(program.Operations, ins)
			case "proc":
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Content

				procNames[pName]++
				if procNames[pName] != 1 {
					parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
						VarName:  orthtypes.Program{},
						VarValue: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_02, "PROC", pName, file.Name, v.Index, v.Content.Index),
					}
					close(parseTokenResult)
					return
				}

				context = pName

				ins := parseToken(orthtypes.PrimitiveProc, pName, orthtypes.Proc)
				program.Operations = append(program.Operations, ins)
			case "in":
				ins := parseToken(orthtypes.PrimitiveIn, "", orthtypes.In)
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
					panic(fmt.Errorf(orth_debug.UndefinedFunction, preProgram[i+1].Content.Content))
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
						ins := parseTokenWithContext(orthtypes.PrimitiveVar, preProgram[i+1].Content.Content, context, orthtypes.Push)
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

				switch vType {
				case orthtypes.PrimitiveSTR:
					fallthrough
				case orthtypes.RNGABL:
					vValue = preProgram[i+4].Content.Content[1 : len(preProgram[i+4].Content.Content)-1]
				default:
					vValue = preProgram[i+4].Content.Content
				}

				ins := parseTokenWithContext(vType, vValue, context, orthtypes.Push)
				program.Operations = append(program.Operations, ins)

				ins = parseTokenWithContext(orthtypes.PrimitiveVar, vName, context, orthtypes.Var)
				program.Operations = append(program.Operations, ins)
			case "hold":
				preProgram[i+1].Content.ValidPos = true
				vName := preProgram[i+1].Content.Content

				ins := parseToken(orthtypes.PrimitiveHold, vName, orthtypes.Hold)
				program.Operations = append(program.Operations, ins)
			case "invoke":
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Content

				ins := parseToken(orthtypes.PrimitiveRNT, pName, orthtypes.Invoke)
				program.Operations = append(program.Operations, ins)
			case "dump_mem":
				ins := parseToken(orthtypes.PrimitiveRNT, "", orthtypes.DumpMem)
				program.Operations = append(program.Operations, ins)
			default:
				if !v.Content.ValidPos {
					parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
						VarName:  orthtypes.Program{},
						VarValue: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_01, v.Content.Content, file.Name, v.Index, v.Content.Index),
					}
					close(parseTokenResult)
					return
				}
			}
		}
	}
	parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
		VarName:  program,
		VarValue: nil,
	}
	close(parseTokenResult)
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

func parseTokenWithContext(varType, operand, context string, op int) orthtypes.Operation {
	return orthtypes.Operation{
		Instruction: op,
		Operand: orthtypes.Operand{
			VarType: varType,
			Operand: operand,
		},
		Context:  context,
		RefBlock: -1,
	}
}
