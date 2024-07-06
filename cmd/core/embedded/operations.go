package embedded

import (
	"fmt"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"regexp"
	"strconv"
)

// CrossReferenceBlocks loops over a program and define all inter references
// needed for execution. Ex: if-else-do blocks
func CrossReferenceBlocks(program orthtypes.Program) (orthtypes.Program, error) {
	stack := make([]embedded_helpers.RefStackItem, 0, len(program.Operations))

	for operationIndex, operation := range program.Operations {
		switch operation.Instruction {
		case orthtypes.Var:
			if operation.Context.Name == embedded_helpers.MainScope {
				program.Variables = append(program.Variables, operation)
			}
		case orthtypes.Const:
			if operation.Context.Name == embedded_helpers.MainScope {
				program.Constants = append(program.Constants, operation)
			}
		case orthtypes.Hold:
			variable, err := operation.Context.GetVaraible(operation.Operator.Operand, &program)
			if err != nil {
				panic(err)
			}
			if program.Operations[operationIndex].Links == nil {
				program.Operations[operationIndex].Links = make(map[string]orthtypes.Operation)
			}
			if variable.Context.Name == embedded_helpers.MainScope {
				program.Operations[operationIndex].Links["hold_mult"] = *variable
			} else {
				program.Operations[operationIndex].Links["hold_local"] = *variable
			}
		case orthtypes.While:
			fallthrough
		case orthtypes.If:
			fallthrough
		case orthtypes.Proc:
			stack = append(stack, embedded_helpers.RefStackItem{
				AbsPosition: uint(operationIndex),
				Instruction: operation.Instruction,
			})
		case orthtypes.Do:
			embedded_helpers.HandleOperationDo(&stack, &program, uint(operationIndex))
			stack = append(stack, embedded_helpers.RefStackItem{
				AbsPosition: uint(operationIndex),
				Instruction: operation.Instruction,
			})
		case orthtypes.Else:
			embedded_helpers.HandleOperationElse(&stack, &program, uint(operationIndex))
			stack = append(stack, embedded_helpers.RefStackItem{
				AbsPosition: uint(operationIndex),
				Instruction: operation.Instruction,
			})
		case orthtypes.End:
			embedded_helpers.HandleOperationEnd(&stack, &program, uint(operationIndex))
		}
	}

	// for _, v := range program.Operations {
	// 	pp := orthtypes.PPrintOperation(v)
	// 	fmt.Println(pp)
	// }

	return program, nil
}

// ParseTokenAsOperation parses an slice of pre-instructions into a runnable program
func ParseTokenAsOperation(tokenFiles []orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]], parsedOperation chan<- orthtypes.Pair[orthtypes.Operation, error]) {
	procNames := make(map[string]int)

	context := &orthtypes.Context{
		Name:         embedded_helpers.MainScope,
		Order:        0,
		Parent:       nil,
		Declarations: make([]orthtypes.ContextDeclaration, 0),
		InnerContext: make([]*orthtypes.Context, 0),
	}

	var globalInstructionIndex uint = 0
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
					Declarations: make([]orthtypes.ContextDeclaration, 0),
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
					Declarations: make([]orthtypes.ContextDeclaration, 0),
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
					Declarations: make([]orthtypes.ContextDeclaration, 0),
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
					Declarations: make([]orthtypes.ContextDeclaration, 0),
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

				context.Declarations = append(context.Declarations, orthtypes.ContextDeclaration{
					Name:  vName,
					Index: globalInstructionIndex,
				})

				value := parseToken(vType, vValue, context, orthtypes.Push)
				constant := parseToken(orthtypes.PrimitiveConst, vName, context, orthtypes.Const)
				constant.Links = make(map[string]orthtypes.Operation)
				constant.Links["variable_value"] = value

				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  constant,
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

				context.Declarations = append(context.Declarations, orthtypes.ContextDeclaration{
					Name:  vName,
					Index: globalInstructionIndex,
				})

				value := parseToken(vType, vValue, context, orthtypes.Push)
				variable := parseToken(orthtypes.PrimitiveVar, vName, context, orthtypes.Var)
				variable.Links = make(map[string]orthtypes.Operation)
				variable.Links["variable_value"] = value

				parsedOperation <- orthtypes.Pair[orthtypes.Operation, error]{
					Left:  variable,
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
			globalInstructionIndex++
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
func parseToken(varType, operand string, context *orthtypes.Context, op orthtypes.Instruction) orthtypes.Operation {
	return orthtypes.Operation{
		Instruction: op,
		Operator: orthtypes.Operand{
			SymbolName: varType,
			Operand:    operand,
		},
		Context:   context,
		Addresses: make(map[orthtypes.Instruction]int),
	}
}
