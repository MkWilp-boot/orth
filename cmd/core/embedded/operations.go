package embedded

import (
	"fmt"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"regexp"
	"strconv"
)

func PopLast[T comparable](root *[]T) T {
	stack := *root
	ret := stack[len(stack)-1]
	*root = stack[:len(stack)-1]
	return ret
}

const MainScope = "_global"

// CrossReferenceBlocks loops over a program and define all inter references
// needed for execution. Ex: if-else-do blocks
func CrossReferenceBlocks(program orthtypes.Program, crossResult chan<- orthtypes.Pair[orthtypes.Program, error]) {

	stack := make([]orthtypes.Pair[int, orthtypes.Operation], 0)

	for ip, currentOperation := range program.Operations {
		pair := orthtypes.Pair[int, orthtypes.Operation]{
			Left:  ip,
			Right: currentOperation,
		}
		switch currentOperation.Instruction {
		case orthtypes.Mem:
			if currentOperation.Context.Name == MainScope {
				msg := fmt.Sprintf(orth_debug.InvalidUsageOfTokenOutside, orthtypes.PrimitiveMem, orthtypes.PrimitiveProc, MainScope)
				crossResult <- orthtypes.Pair[orthtypes.Program, error]{
					Left:  orthtypes.Program{},
					Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_04, "Mem", msg),
				}
				close(crossResult)
				return
			}
		case orthtypes.Hold:
			variableDeclared := currentOperation.Context.HasVariableDeclaredInOrAbove(currentOperation.Operator.Operand)
			if !variableDeclared {
				err := orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_11, currentOperation.Operator.Operand, "hold")
				crossResult <- orthtypes.Pair[orthtypes.Program, error]{
					Left:  orthtypes.Program{},
					Right: err,
				}
				continue
			}

			for i := ip; i >= 0; i-- {
				isVar := program.Operations[i].Operator.VarType == orthtypes.PrimitiveConst ||
					program.Operations[i].Operator.VarType == orthtypes.PrimitiveVar

				if isVar && program.Operations[i].Operator.Operand == currentOperation.Operator.Operand && currentOperation.RefBlock == -1 {
					program.Operations[ip].RefBlock = i
					break
				}
			}
		case orthtypes.SetString:
			if ip-2 < 0 {
				err := orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_09, "set_string")
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}

			holdingVariable := program.Operations[ip-1]
			newValue := program.Operations[ip-2]

			isString := helpers.IsString(newValue.Operator.VarType)

			if !isString {
				err := orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_08,
					"set_string",
					"ptr",
					"string",
					holdingVariable.Operator.Operand,
					newValue.Operator.VarType,
				)
				crossResult <- orthtypes.Pair[orthtypes.Program, error]{
					Left:  orthtypes.Program{},
					Right: err,
				}
				continue
			}
		case orthtypes.If:
			fallthrough
		case orthtypes.Proc:
			fallthrough
		case orthtypes.While:
			stack = append(stack, pair)
		case orthtypes.Else:
			blockIp := PopLast(&stack)

			if program.Operations[blockIp.Left].Instruction != orthtypes.If {
				fmt.Fprintln(os.Stderr, "Invalid Else clause")
				os.Exit(1)
			}

			program.Operations[blockIp.Left].RefBlock = ip + 1
			stack = append(stack, pair)
		case orthtypes.End:
			blockIp := PopLast(&stack)
			switch {
			case program.Operations[blockIp.Left].Instruction == orthtypes.If:
				fallthrough
			case program.Operations[blockIp.Left].Instruction == orthtypes.Else:
				program.Operations[blockIp.Left].RefBlock = ip
				program.Operations[ip].RefBlock = ip + 1 // end block
			case program.Operations[blockIp.Left].Instruction == orthtypes.In:
				//context = globalScope
				fallthrough
			case program.Operations[blockIp.Left].Instruction == orthtypes.Do:
				if program.Operations[blockIp.Left].RefBlock == -1 {
					fmt.Fprintln(os.Stderr, "Not enought arguments for a cross-refernce block operation")
					os.Exit(1)
				}
				program.Operations[ip].RefBlock = program.Operations[blockIp.Left].RefBlock
				program.Operations[blockIp.Left].RefBlock = ip + 1
			default:
				fmt.Fprintln(os.Stderr, "End block can only close [if | else | do | proc in] blocks")
				os.Exit(1)
			}
		case orthtypes.In:
			fallthrough
		case orthtypes.Do:
			blockIp := PopLast(&stack)
			program.Operations[ip].RefBlock = blockIp.Left
			stack = append(stack, pair)
		}
	}
	crossResult <- orthtypes.Pair[orthtypes.Program, error]{
		Left:  program,
		Right: nil,
	}
	close(crossResult)
}

// ParseTokenAsOperation parses an slice of pre-instructions into a runnable program
func ParseTokenAsOperation(tokenFiles []orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]], parseTokenResult chan<- orthtypes.Pair[orthtypes.Program, error]) {
	program := orthtypes.Program{
		Id: len(tokenFiles),
	}
	procNames := make(map[string]int)

	context := &orthtypes.Context{
		Name:         MainScope,
		Order:        0,
		Parent:       nil,
		Declarations: make([]string, 0),
		InnerContext: make([]*orthtypes.Context, 0),
	}

	for fIndex, file := range tokenFiles {
		for i, v := range *file.CodeBlock.Slice {
			preProgram := (*tokenFiles[fIndex].CodeBlock.Slice)

			if v.Content.ValidPos {
				continue
			}
			switch v.Content.Token {
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
				ins := parseToken(v.Content.Token, preProgram[i+1].Content.Token, context, orthtypes.Push)
				program.Operations = append(program.Operations, ins)
			case orthtypes.PrimitiveSTR:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Token[1:len(preProgram[i+1].Content.Token)-1], context, orthtypes.PushStr)
				program.Operations = append(program.Operations, ins)
			case "+":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Sum)
				program.Operations = append(program.Operations, ins)
			case "-":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Minus)
				program.Operations = append(program.Operations, ins)
			case "*":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Mult)
				program.Operations = append(program.Operations, ins)
			case "/":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Div)
				program.Operations = append(program.Operations, ins)
			case "putui":
				ins := parseToken(orthtypes.PrimitiveVOID, "", context, orthtypes.PutU64)
				program.Operations = append(program.Operations, ins)
			case "==":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Equal)
				program.Operations = append(program.Operations, ins)
			case "<>":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.NotEqual)
				program.Operations = append(program.Operations, ins)
			case "<":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Lt)
				program.Operations = append(program.Operations, ins)
			case ">":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Gt)
				program.Operations = append(program.Operations, ins)
			case "if":
				newContext := orthtypes.Context{
					Name:         fmt.Sprintf("c?_if_%d$", len(context.InnerContext)),
					Parent:       context,
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]string, 0),
					InnerContext: make([]*orthtypes.Context, 0),
				}
				context.InnerContext = append(context.InnerContext, &newContext)
				context = &newContext

				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.If)
				program.Operations = append(program.Operations, ins)
			case "else":
				// context is an "if" block that must not have a "else" block as a child, they should be siblings
				// so no appending to context.InnerContext due to this restriction
				newContext := orthtypes.Context{
					Name:         fmt.Sprintf("c?_else_%d$", len(context.InnerContext)),
					Parent:       context.Parent, // else is not a child of "if"
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]string, 0),
					InnerContext: make([]*orthtypes.Context, 0),
				}

				context = &newContext
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Else)
				program.Operations = append(program.Operations, ins)
			case "end":
				if context.Parent != nil {
					context = context.Parent
				}
				ins := parseToken(orthtypes.PrimitiveEND, "", context, orthtypes.End)
				program.Operations = append(program.Operations, ins)
			case "puts":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.PutString)
				program.Operations = append(program.Operations, ins)
			case "over":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Over)
				program.Operations = append(program.Operations, ins)
			case "2dup":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.TwoDup)
				program.Operations = append(program.Operations, ins)
			case "dup":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Dup)
				program.Operations = append(program.Operations, ins)
			case "while":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.While)
				program.Operations = append(program.Operations, ins)
			case "lshift":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.LShift)
				program.Operations = append(program.Operations, ins)
			case "rshift":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.RShift)
				program.Operations = append(program.Operations, ins)
			case "land":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.LAnd)
				program.Operations = append(program.Operations, ins)
			case "lor":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.LOr)
				program.Operations = append(program.Operations, ins)
			case "proc":
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Token

				procNames[pName]++
				if procNames[pName] != 1 {
					parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
						Left:  orthtypes.Program{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_02, "PROC", pName, file.Name, v.Index, v.Content.Index),
					}
					close(parseTokenResult)
					return
				}

				newContext := orthtypes.Context{
					Name:         fmt.Sprintf("c?_proc_%s_%d$", pName, len(context.InnerContext)),
					Parent:       context,
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]string, 0),
					InnerContext: make([]*orthtypes.Context, 0),
				}
				context.InnerContext = append(context.InnerContext, &newContext)
				context = &newContext

				ins := parseToken(orthtypes.PrimitiveProc, pName, context, orthtypes.Proc)
				program.Operations = append(program.Operations, ins)
			case "in":
				ins := parseToken(orthtypes.PrimitiveIn, "", context, orthtypes.In)
				program.Operations = append(program.Operations, ins)
			case "do":
				newContext := orthtypes.Context{
					Name:         fmt.Sprintf("c?_do_%d$", len(context.InnerContext)),
					Parent:       context,
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]string, 0),
					InnerContext: make([]*orthtypes.Context, 0),
				}
				context.InnerContext = append(context.InnerContext, &newContext)
				context = &newContext

				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Do)
				program.Operations = append(program.Operations, ins)
			case "drop":
				ins := parseToken(orthtypes.PrimitiveVOID, "", context, orthtypes.Drop)
				program.Operations = append(program.Operations, ins)
			case "swap":
				ins := parseToken(orthtypes.PrimitiveVOID, "", context, orthtypes.Swap)
				program.Operations = append(program.Operations, ins)
			case "%":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Mod)
				program.Operations = append(program.Operations, ins)
			case orthtypes.PrimitiveMem:
				ins := parseToken(orthtypes.PrimitiveMem, orthtypes.PrimitiveMem, context, orthtypes.Mem)
				program.Operations = append(program.Operations, ins)
			case ".":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Store)
				program.Operations = append(program.Operations, ins)
			case ",":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Load)
				program.Operations = append(program.Operations, ins)
			case "call":
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Token, context, orthtypes.Call)
				program.Operations = append(program.Operations, ins)
			case ",!":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.LoadStay)
				program.Operations = append(program.Operations, ins)
			case "type":
				preProgram[i+1].Content.ValidPos = true

				ins := parseToken(orthtypes.PrimitiveType, preProgram[i+1].Content.Token, context, orthtypes.Push)

				program.Operations = append(program.Operations, ins)
			case "const":
				vValue, vType, vName := grabVariableDefinition(preProgram, i, &program)

				if context.HasVariableDeclaredInOrAbove(vName) {
					parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
						Left:  orthtypes.Program{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "constant", vName, context.Name),
					}
					close(parseTokenResult)
					return
				}

				context.Declarations = append(context.Declarations, vName)

				ins := parseToken(vType, vValue, context, orthtypes.Push)
				program.Operations = append(program.Operations, ins)

				ins = parseToken(orthtypes.PrimitiveConst, vName, context, orthtypes.Const)
				program.Operations = append(program.Operations, ins)
			case "var":
				vValue, vType, vName := grabVariableDefinition(preProgram, i, &program)

				if context.HasVariableDeclaredInOrAbove(vName) {
					parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
						Left:  orthtypes.Program{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "variable", vName, context.Name),
					}
					close(parseTokenResult)
					return
				}

				context.Declarations = append(context.Declarations, vName)

				ins := parseToken(vType, vValue, context, orthtypes.Push)
				program.Operations = append(program.Operations, ins)

				ins = parseToken(orthtypes.PrimitiveVar, vName, context, orthtypes.Var)
				program.Operations = append(program.Operations, ins)
			case "deref":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Deref)
				program.Operations = append(program.Operations, ins)
			case "set_number":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.SetNumber)
				program.Operations = append(program.Operations, ins)
			case "set_string":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.SetString)
				program.Operations = append(program.Operations, ins)
			case orthtypes.PrimitiveHold:
				preProgram[i+1].Content.ValidPos = true
				vName := preProgram[i+1].Content.Token

				ins := parseToken(orthtypes.PrimitiveHold, vName, context, orthtypes.Hold)
				program.Operations = append(program.Operations, ins)
			case "invoke":
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Token

				ins := parseToken(orthtypes.PrimitiveRNT, pName, context, orthtypes.Invoke)
				program.Operations = append(program.Operations, ins)
			case "exit":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Exit)
				program.Operations = append(program.Operations, ins)
			case "out":
				preProgram[i+1].Content.ValidPos = true
				amountStr := preProgram[i+1].Content.Token

				amount, err := strconv.Atoi(amountStr)
				if err != nil {
					errStr := orth_debug.BuildErrorMessage(
						orth_debug.ORTH_ERR_05,
						"out",
						"i~",
						amountStr,
						file.Name, v.Index, v.Content.Index,
					)
					panic(errStr)
				}

				if amount > orthtypes.MAX_PROC_OUTPUT_COUNT {
					errStr := orth_debug.BuildErrorMessage(
						orth_debug.ORTH_ERR_06,
						orthtypes.MAX_PROC_OUTPUT_COUNT,
						amount,
						file.Name, v.Index, v.Content.Index,
					)
					panic(errStr)
				}

				ins := parseToken(orthtypes.PrimitiveRNT, amountStr, context, orthtypes.Out)
				program.Operations = append(program.Operations, ins)
			case "with":
				preProgram[i+1].Content.ValidPos = true
				amountStr := preProgram[i+1].Content.Token

				if amountStr != "cli" {
					amount, err := strconv.Atoi(amountStr)
					if err != nil {
						errStr := orth_debug.BuildErrorMessage(
							orth_debug.ORTH_ERR_05,
							"with",
							"i~ | cli",
							amountStr,
							file.Name, v.Index, v.Content.Index,
						)
						panic(errStr)
					}

					if amount > orthtypes.MAX_PROC_PARAM_COUNT {
						errStr := orth_debug.BuildErrorMessage(
							orth_debug.ORTH_ERR_06,
							orthtypes.MAX_PROC_PARAM_COUNT,
							amount,
							file.Name, v.Index, v.Content.Index,
						)
						panic(errStr)
					}
				}

				ins := parseToken(orthtypes.PrimitiveRNT, amountStr, context, orthtypes.With)
				program.Operations = append(program.Operations, ins)
			case "dump_mem":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.DumpMem)
				program.Operations = append(program.Operations, ins)
			default:
				if !v.Content.ValidPos {
					parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
						Left:  orthtypes.Program{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_01, v.Content.Token, file.Name, v.Index, v.Content.Index),
					}
					close(parseTokenResult)
					return
				}
			}
		}
	}

	parseTokenResult <- orthtypes.Pair[orthtypes.Program, error]{
		Left:  program,
		Right: nil,
	}

	close(parseTokenResult)
}

func grabVariableDefinition(preProgram []orthtypes.StringEnum, i int, program *orthtypes.Program) (string, string, string) {
	re := regexp.MustCompile(`[^\w]`)

	// check name
	if re.Match([]byte(preProgram[i+1].Content.Token)) {
		fmt.Fprintf(os.Stderr, "%s has invalid characters in it's composition\n", "const")
		os.Exit(1)
	}
	// check if has a value
	if preProgram[i+2].Content.Token != "=" {
		switch {
		// used as a func param is currently disabled
		// case preProgram[i+2].Content.Token == "call":
		// 	preProgram[i+1].Content.ValidPos = true
		// 	ins := parseTokenWithContext(orthtypes.PrimitiveConst, preProgram[i+1].Content.Token, context, orthtypes.Push)
		// 	program.Operations = append(program.Operations, ins)
		// 	continue
		default:
			fmt.Fprintln(os.Stderr, "var must be initialized with `=` sign")
			os.Exit(1)
		}
	}

	for x := 1; x < 5; x++ {
		preProgram[i+x].Content.ValidPos = true
	}

	vName := preProgram[i+1].Content.Token
	vType := preProgram[i+3].Content.Token

	var vValue string

	switch vType {
	case orthtypes.PrimitiveSTR:
		fallthrough
	case orthtypes.RNGABL:
		vValue = preProgram[i+4].Content.Token[1 : len(preProgram[i+4].Content.Token)-1]
	default:
		vValue = preProgram[i+4].Content.Token
	}

	return vValue, vType, vName
}

// parseToken parses a single token into a instruction
func parseToken(varType, operand string, context *orthtypes.Context, op int) orthtypes.Operation {
	return orthtypes.Operation{
		Instruction: op,
		Operator: orthtypes.Operand{
			VarType: varType,
			Operand: operand,
		},
		Context:  context,
		RefBlock: -1,
	}
}
