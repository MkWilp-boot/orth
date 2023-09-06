package embedded

import (
	"fmt"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"regexp"
	"strconv"
)

// CrossReferenceBlocks loops over a program and define all inter references
// needed for execution. Ex: if-else-do blocks
func CrossReferenceBlocks(program orthtypes.Program) (orthtypes.Program, error) {
	var err error
	stack := make([]orthtypes.Pair[int, orthtypes.Operation], 0)

	for ip, currentOperation := range program.Operations {
		if err != nil {
			break
		}
		pair := orthtypes.Pair[int, orthtypes.Operation]{
			Left:  ip,
			Right: currentOperation,
		}
		switch currentOperation.Instruction {
		case orthtypes.Mem:
			if currentOperation.Context.Name == embedded_helpers.MainScope {
				msg := fmt.Sprintf(orth_debug.InvalidUsageOfTokenOutside, orthtypes.PrimitiveMem, orthtypes.PrimitiveProc, embedded_helpers.MainScope)
				err = orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_04, "Mem", msg)
			}
		case orthtypes.Hold:
			variableDeclared := currentOperation.Context.HasVariableDeclaredInOrAbove(currentOperation.Operator.Operand)
			if !variableDeclared {
				err = orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_11, currentOperation.Operator.Operand, "hold")
				continue
			}

			declarationContext, err := embedded_helpers.GetVariableContext(currentOperation.Operator.Operand, currentOperation.Context)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			var lookTable *[]orthtypes.Operation
			if currentOperation.Operator.SymbolName == orthtypes.PrimitiveVar {
				lookTable = &program.Variables
			} else {
				lookTable = &program.Constants
			}

			for i := len(*lookTable) - 1; i >= 0; i-- {
				isVar := (*lookTable)[i].Operator.SymbolName == orthtypes.PrimitiveConst || (*lookTable)[i].Operator.SymbolName == orthtypes.PrimitiveVar
				nameMatch := (*lookTable)[i].Operator.Operand == currentOperation.Operator.Operand
				contextMatch := (*lookTable)[i].Context.Name == declarationContext

				if isVar && nameMatch && contextMatch && currentOperation.RefBlock == -1 {
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

			isString := helpers.IsString(newValue.Operator.SymbolName)

			if !isString {
				err = orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_08,
					"set_string",
					"ptr",
					"string",
					holdingVariable.Operator.Operand,
					newValue.Operator.SymbolName,
				)
				continue
			}
		case orthtypes.If:
			fallthrough
		case orthtypes.Proc:
			fallthrough
		case orthtypes.While:
			stack = append(stack, pair)
		case orthtypes.Else:
			blockIp := embedded_helpers.PopLast(&stack)

			if program.Operations[blockIp.Left].Instruction != orthtypes.If {
				fmt.Fprintln(os.Stderr, "Invalid Else clause")
				os.Exit(1)
			}

			program.Operations[blockIp.Left].RefBlock = ip + 1
			stack = append(stack, pair)
		case orthtypes.End:
			blockIp := embedded_helpers.PopLast(&stack)
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
			blockIp := embedded_helpers.PopLast(&stack)
			program.Operations[ip].RefBlock = blockIp.Left
			stack = append(stack, pair)
		}
	}

	return program, err
}

// ParseTokenAsOperation parses an slice of pre-instructions into a runnable program
func ParseTokenAsOperation(tokenFiles []orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]], parsedOperation chan<- orthtypes.Pair[orthtypes.Operation, error]) {
	procNames := make(map[string]int)

	context := &orthtypes.Context{
		Name:         embedded_helpers.MainScope,
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orthtypes.PrimitiveSTR:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Token[1:len(preProgram[i+1].Content.Token)-1], context, orthtypes.PushStr)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "+":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Sum)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "-":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Minus)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "*":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Mult)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "/":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Div)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "putui":
				ins := parseToken(orthtypes.PrimitiveVOID, "", context, orthtypes.PutU64)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "==":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Equal)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "<>":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.NotEqual)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "<":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Lt)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case ">":
				ins := parseToken(orthtypes.PrimitiveBOOL, "", context, orthtypes.Gt)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "end":
				if context.Parent != nil {
					context = context.Parent
				}
				ins := parseToken(orthtypes.PrimitiveEND, "", context, orthtypes.End)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "puts":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.PutString)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "over":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Over)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "2dup":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.TwoDup)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "dup":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Dup)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "while":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.While)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "lshift":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.LShift)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "rshift":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.RShift)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "land":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.LAnd)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "lor":
				ins := parseToken(orthtypes.Bitwise, "", context, orthtypes.LOr)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "proc":
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Token

				procNames[pName]++
				if procNames[pName] != 1 {
					parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
						Left:  orthtypes.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_02, "PROC", pName, file.Name, v.Index, v.Content.Index),
					}
					close(parsedOperation)
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "in":
				ins := parseToken(orthtypes.PrimitiveIn, "", context, orthtypes.In)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "drop":
				ins := parseToken(orthtypes.PrimitiveVOID, "", context, orthtypes.Drop)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "swap":
				ins := parseToken(orthtypes.PrimitiveVOID, "", context, orthtypes.Swap)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "%":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Mod)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orthtypes.PrimitiveMem:
				ins := parseToken(orthtypes.PrimitiveMem, orthtypes.PrimitiveMem, context, orthtypes.Mem)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case ".":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Store)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case ",":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Load)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "call":
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orthtypes.PrimitiveSTR, preProgram[i+1].Content.Token, context, orthtypes.Call)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case ",!":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.LoadStay)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "type":
				preProgram[i+1].Content.ValidPos = true

				ins := parseToken(orthtypes.PrimitiveType, preProgram[i+1].Content.Token, context, orthtypes.Push)

				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "const":
				vValue, vType, vName := grabVariableDefinition(preProgram, i)

				if context.HasVariableDeclaredInOrAbove(vName) {
					parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
						Left:  orthtypes.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "constant", vName, context.Name),
					}
					close(parsedOperation)
					return
				}

				context.Declarations = append(context.Declarations, vName)

				ins := parseToken(vType, vValue, context, orthtypes.Push)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}

				ins = parseToken(orthtypes.PrimitiveConst, vName, context, orthtypes.Const)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "var":
				vValue, vType, vName := grabVariableDefinition(preProgram, i)

				if context.HasVariableDeclaredInOrAbove(vName) {
					parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
						Left:  orthtypes.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "variable", vName, context.Name),
					}
					close(parsedOperation)
					return
				}

				context.Declarations = append(context.Declarations, vName)

				ins := parseToken(vType, vValue, context, orthtypes.Push)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}

				ins = parseToken(orthtypes.PrimitiveVar, vName, context, orthtypes.Var)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "deref":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Deref)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "set_number":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.SetNumber)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "set_string":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.SetString)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orthtypes.PrimitiveHold:
				preProgram[i+1].Content.ValidPos = true
				vName := preProgram[i+1].Content.Token

				ins := parseToken(orthtypes.PrimitiveHold, vName, context, orthtypes.Hold)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "invoke":
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Token

				ins := parseToken(orthtypes.PrimitiveRNT, pName, context, orthtypes.Invoke)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "exit":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Exit)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
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
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "dump_mem":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.DumpMem)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "alloc":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Alloc)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "free":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.Free)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case "put_char":
				ins := parseToken(orthtypes.PrimitiveRNT, "", context, orthtypes.PutChar)
				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			default:
				if !v.Content.ValidPos {
					parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
						Left:  orthtypes.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_01, v.Content.Token, file.Name, v.Index, v.Content.Index),
					}
					close(parsedOperation)
					return
				}
			}
		}
	}

	close(parsedOperation)
}

func grabVariableDefinition(preProgram []orthtypes.StringEnum, i int) (string, string, string) {
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
			SymbolName: varType,
			Operand:    operand,
		},
		Context:  context,
		RefBlock: -1,
	}
}
