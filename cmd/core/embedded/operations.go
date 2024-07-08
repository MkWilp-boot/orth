package embedded

import (
	"fmt"
	embedded_helpers "orth/cmd/core/embedded/helpers"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"os"
	"regexp"
	"strconv"
)

// CrossReferenceBlocks loops over a program and define all inter references
// needed for execution. Ex: if-else-do blocks
func CrossReferenceBlocks(program orth_types.Program) (orth_types.Program, error) {
	stack := make([]embedded_helpers.RefStackItem, 0, len(program.Operations))

	for operationIndex, operation := range program.Operations {
		switch operation.Instruction {
		case orth_types.InstructionVar:
			if operation.Context.Name == embedded_helpers.MainScope {
				program.Variables = append(program.Variables, operation)
			}
		case orth_types.InstructionConst:
			if operation.Context.Name == embedded_helpers.MainScope {
				program.Constants = append(program.Constants, operation)
			}
		case orth_types.InstructionHold:
			variable, err := operation.Context.GetVaraible(operation.Operator.Operand, &program)
			if err != nil {
				panic(err)
			}
			if program.Operations[operationIndex].Links == nil {
				program.Operations[operationIndex].Links = make(map[string]orth_types.Operation)
			}
			if variable.Context.Name == embedded_helpers.MainScope {
				program.Operations[operationIndex].Links["hold_mult"] = *variable
			} else {
				program.Operations[operationIndex].Links["hold_local"] = *variable
			}
		case orth_types.InstructionWhile:
			fallthrough
		case orth_types.InstructionIf:
			fallthrough
		case orth_types.InstructionProc:
			stack = append(stack, embedded_helpers.RefStackItem{
				AbsPosition: uint(operationIndex),
				Instruction: operation.Instruction,
			})
		case orth_types.InstructionDo:
			embedded_helpers.HandleOperationDo(&stack, &program, uint(operationIndex))
			stack = append(stack, embedded_helpers.RefStackItem{
				AbsPosition: uint(operationIndex),
				Instruction: operation.Instruction,
			})
		case orth_types.InstructionElse:
			embedded_helpers.HandleOperationElse(&stack, &program, uint(operationIndex))
			stack = append(stack, embedded_helpers.RefStackItem{
				AbsPosition: uint(operationIndex),
				Instruction: operation.Instruction,
			})
		case orth_types.InstructionEnd:
			embedded_helpers.HandleOperationEnd(&stack, &program, uint(operationIndex))
		}
	}

	return program, nil
}

// ParseTokenAsOperation parses an slice of pre-instructions into a runnable program
func ParseTokenAsOperation(tokenFiles []orth_types.File[orth_types.SliceOf[orth_types.StringEnum]], parsedOperation chan<- orth_types.Pair[orth_types.Operation, error]) {
	procNames := make(map[string]int)

	context := &orth_types.Context{
		Name:         embedded_helpers.MainScope,
		Order:        0,
		Parent:       nil,
		Declarations: make([]orth_types.ContextDeclaration, 0),
		InnerContext: make([]*orth_types.Context, 0),
	}

	var globalInstructionIndex uint = 0
	for fIndex, file := range tokenFiles {
		for i, v := range *file.CodeBlock.Slice {
			preProgram := (*tokenFiles[fIndex].CodeBlock.Slice)

			if v.Content.ValidPos {
				continue
			}
			switch v.Content.Token {
			case orth_types.ADDR:
				fallthrough
			case orth_types.StdRNT:
				fallthrough
			case orth_types.StdINT:
				fallthrough
			case orth_types.StdI8:
				fallthrough
			case orth_types.StdI16:
				fallthrough
			case orth_types.StdI32:
				fallthrough
			case orth_types.StdI64:
				fallthrough
			case orth_types.StdF32:
				fallthrough
			case orth_types.StdF64:
				fallthrough
			case orth_types.StdBOOL:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(v.Content.Token, preProgram[i+1].Content.Token, context, orth_types.InstructionPush)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdSTR:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orth_types.StdSTR, preProgram[i+1].Content.Token[1:len(preProgram[i+1].Content.Token)-1], context, orth_types.InstructionPushStr)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdPlus:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionSum)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdMinus:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionMinus)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdMult:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionMult)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdDiv:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionDiv)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdPutUint:
				ins := parseToken(orth_types.StdVOID, "", context, orth_types.FunctionPutU64)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdEquals:
				ins := parseToken(orth_types.StdBOOL, "", context, orth_types.InstructionEqual)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdNotEquals:
				ins := parseToken(orth_types.StdBOOL, "", context, orth_types.InstructionNotEqual)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdLowerThan:
				ins := parseToken(orth_types.StdBOOL, "", context, orth_types.InstructionLt)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdGreaterThan:
				ins := parseToken(orth_types.StdBOOL, "", context, orth_types.InstructionGt)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdIf:
				newContext := orth_types.Context{
					Name:         fmt.Sprintf("c?_if_%d$", len(context.InnerContext)),
					Parent:       context,
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]orth_types.ContextDeclaration, 0),
					InnerContext: make([]*orth_types.Context, 0),
				}
				context.InnerContext = append(context.InnerContext, &newContext)
				context = &newContext

				ins := parseToken(orth_types.StdBOOL, "", context, orth_types.InstructionIf)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdElse:
				// context is an "if" block that must not have a "else" block as a child, they should be siblings
				// so no appending to context.InnerContext due to this restriction
				newContext := orth_types.Context{
					Name:         fmt.Sprintf("c?_else_%d$", len(context.InnerContext)),
					Parent:       context.Parent, // else is not a child of "if"
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]orth_types.ContextDeclaration, 0),
					InnerContext: make([]*orth_types.Context, 0),
				}

				context = &newContext
				ins := parseToken(orth_types.StdBOOL, "", context, orth_types.InstructionElse)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdEND:
				if context.Parent != nil {
					context = context.Parent
				}
				ins := parseToken(orth_types.StdEND, "", context, orth_types.InstructionEnd)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdPutStr:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionPutString)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdOver:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionOver)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.Std2Dup:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionTwoDup)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdDup:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionDup)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdWhile:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionWhile)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdLeftShift:
				ins := parseToken(orth_types.StdBitwise, "", context, orth_types.InstructionLShift)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdRightShift:
				ins := parseToken(orth_types.StdBitwise, "", context, orth_types.InstructionRShift)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdLogicalAnd:
				ins := parseToken(orth_types.StdBitwise, "", context, orth_types.InstructionLAnd)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdLogicalOr:
				ins := parseToken(orth_types.StdBitwise, "", context, orth_types.InstructionLOr)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdProc:
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Token

				procNames[pName]++
				if procNames[pName] != 1 {
					parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
						Left:  orth_types.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_02, "PROC", pName, file.Name, v.Index, v.Content.Index),
					}
					close(parsedOperation)
					return
				}

				newContext := orth_types.Context{
					Name:         fmt.Sprintf("c?_proc_%s_%d$", pName, len(context.InnerContext)),
					Parent:       context,
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]orth_types.ContextDeclaration, 0),
					InnerContext: make([]*orth_types.Context, 0),
				}
				context.InnerContext = append(context.InnerContext, &newContext)
				context = &newContext

				ins := parseToken(orth_types.StdProc, pName, context, orth_types.InstructionProc)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdIn:
				ins := parseToken(orth_types.StdIn, "", context, orth_types.InstructionIn)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdDo:
				newContext := orth_types.Context{
					Name:         fmt.Sprintf("c?_do_%d$", len(context.InnerContext)),
					Parent:       context,
					Order:        uint(len(context.InnerContext)),
					Declarations: make([]orth_types.ContextDeclaration, 0),
					InnerContext: make([]*orth_types.Context, 0),
				}
				context.InnerContext = append(context.InnerContext, &newContext)
				context = &newContext

				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionDo)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdDrop:
				ins := parseToken(orth_types.StdVOID, "", context, orth_types.InstructionDrop)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdSwap:
				ins := parseToken(orth_types.StdVOID, "", context, orth_types.InstructionSwap)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdMod:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionMod)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdMem:
				ins := parseToken(orth_types.StdMem, orth_types.StdMem, context, orth_types.InstructionMem)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdStore:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionStore)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdLoad:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionLoad)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdCall:
				preProgram[i+1].Content.ValidPos = true
				ins := parseToken(orth_types.StdSTR, preProgram[i+1].Content.Token, context, orth_types.InstructionCall)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdLoadAndStay:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionLoadStay)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdType:
				preProgram[i+1].Content.ValidPos = true

				ins := parseToken(orth_types.StdType, preProgram[i+1].Content.Token, context, orth_types.InstructionPush)

				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdConst:
				vValue, vType, vName := grabVariableDefinition(preProgram, i)

				if context.HasVariableDeclaredInOrAbove(vName) {
					parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
						Left:  orth_types.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "constant", vName, context.Name),
					}
					close(parsedOperation)
					return
				}

				context.Declarations = append(context.Declarations, orth_types.ContextDeclaration{
					Name:  vName,
					Index: globalInstructionIndex,
				})

				value := parseToken(vType, vValue, context, orth_types.InstructionPush)
				constant := parseToken(orth_types.StdConst, vName, context, orth_types.InstructionConst)
				constant.Links = make(map[string]orth_types.Operation)
				constant.Links["variable_value"] = value

				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  constant,
					Right: nil,
				}
			case orth_types.StdVar:
				vValue, vType, vName := grabVariableDefinition(preProgram, i)

				if context.HasVariableDeclaredInOrAbove(vName) {
					parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
						Left:  orth_types.Operation{},
						Right: orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_03, "variable", vName, context.Name),
					}
					close(parsedOperation)
					return
				}

				context.Declarations = append(context.Declarations, orth_types.ContextDeclaration{
					Name:  vName,
					Index: globalInstructionIndex,
				})

				value := parseToken(vType, vValue, context, orth_types.InstructionPush)
				variable := parseToken(orth_types.StdVar, vName, context, orth_types.InstructionVar)
				variable.Links = make(map[string]orth_types.Operation)
				variable.Links["variable_value"] = value

				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  variable,
					Right: nil,
				}
			case orth_types.StdDeref:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionDeref)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdSetNumber:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionSetNumber)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdSetStr:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionSetString)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdHold:
				preProgram[i+1].Content.ValidPos = true
				vName := preProgram[i+1].Content.Token

				ins := parseToken(orth_types.StdHold, vName, context, orth_types.InstructionHold)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdInvoke:
				preProgram[i+1].Content.ValidPos = true
				pName := preProgram[i+1].Content.Token

				ins := parseToken(orth_types.StdRNT, pName, context, orth_types.InstructionInvoke)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdExit:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.InstructionExit)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdOut:
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

				if amount > orth_types.MAX_PROC_OUTPUT_COUNT {
					errStr := orth_debug.BuildErrorMessage(
						orth_debug.ORTH_ERR_06,
						orth_types.MAX_PROC_OUTPUT_COUNT,
						amount,
						file.Name, v.Index, v.Content.Index,
					)
					panic(errStr)
				}

				ins := parseToken(orth_types.StdRNT, amountStr, context, orth_types.InstructionOut)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdWith:
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

					if amount > orth_types.MAX_PROC_PARAM_COUNT {
						errStr := orth_debug.BuildErrorMessage(
							orth_debug.ORTH_ERR_06,
							orth_types.MAX_PROC_PARAM_COUNT,
							amount,
							file.Name, v.Index, v.Content.Index,
						)
						panic(errStr)
					}
				}

				ins := parseToken(orth_types.StdRNT, amountStr, context, orth_types.InstructionWith)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdDumpMem:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionDumpMem)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdAlloc:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionAlloc)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdFree:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionFree)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			case orth_types.StdPutChar:
				ins := parseToken(orth_types.StdRNT, "", context, orth_types.FunctionPutChar)
				parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
					Left:  ins,
					Right: nil,
				}
			default:
				if !v.Content.ValidPos {
					parsedOperation <- orth_types.Pair[orth_types.Operation, error]{
						Left:  orth_types.Operation{},
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

func grabVariableDefinition(preProgram []orth_types.StringEnum, i int) (string, string, string) {
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
		// 	ins := parseTokenWithContext(orth_types.PrimitiveConst, preProgram[i+1].Content.Token, context, orth_types.Push)
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
	case orth_types.StdSTR:
		fallthrough
	case orth_types.RNGABL:
		vValue = preProgram[i+4].Content.Token[1 : len(preProgram[i+4].Content.Token)-1]
	default:
		vValue = preProgram[i+4].Content.Token
	}

	return vValue, vType, vName
}

// parseToken parses a single token into a instruction
func parseToken(varType, operand string, context *orth_types.Context, op orth_types.Instruction) orth_types.Operation {
	return orth_types.Operation{
		Instruction: op,
		Operator: orth_types.Operand{
			SymbolName: varType,
			Operand:    operand,
		},
		Context:   context,
		Addresses: make(map[orth_types.Instruction]int),
	}
}
