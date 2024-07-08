package embedded_helpers

import (
	"fmt"
	"log"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var AsmVariablePriority = map[string]uint8{
	"REAL10": 10,
	"QWORD":  8,
	"REAL8":  8,
	"REAL4":  4,
	"DWORD":  4,
	"WORD":   2,
	"BYTE":   1,
}

func CleanUp() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("CleanUp Error:", fmt.Sprintf("%v", err))
		}
	}()

	var rmvFiles = []string{
		"mllink$.lnk",
		*orth_debug.ObjectName + ".ilk",
		*orth_debug.ObjectName + ".obj",
		*orth_debug.ObjectName + ".pdb",
	}

	for _, file := range rmvFiles {
		orth_debug.LogStep("[CMD] Deleting extra files")
		if err := os.Remove(file); err != nil {
			panic(err)
		}
	}
}

func VarTypeToLocalAsmType(operand orth_types.Operand) string {
	switch operand.SymbolName {
	case orth_types.StdSTR:
		panic("string not supported for local scopes")
	case orth_types.StdI8:
		return "BYTE"
	case orth_types.StdI16:
		return "WORD"
	case orth_types.StdI32:
		return "DWORD"
	case orth_types.StdINT:
		if strings.Contains(runtime.GOARCH, "64") {
			return "QWORD"
		} else {
			return "DWORD"
		}
	case orth_types.StdI64:
		return "QWORD"
	case orth_types.StdF32:
		return "REAL4"
	case orth_types.StdF64:
		return "REAL8"
	default:
		fmt.Fprintf(os.Stderr, "ivalid type od %q\n", operand.SymbolName)
		os.Exit(1)
		return ""
	}
}

func VarTypeToAsmType(operand orth_types.Operand) string {
	var asmTypeInstruction string
	switch operand.SymbolName {
	case orth_types.StdSTR:
		asmTypeInstruction = "db"
	case orth_types.StdI8:
		asmTypeInstruction = "byte"
	case orth_types.StdI16:
		asmTypeInstruction = "dw"
	case orth_types.StdI32:
		asmTypeInstruction = "dd"
	case orth_types.StdINT:
		if strings.Contains(runtime.GOARCH, "64") {
			asmTypeInstruction = "dq"
		} else {
			asmTypeInstruction = "dd"
		}
	case orth_types.StdI64:
		asmTypeInstruction = "dq"
	case orth_types.StdF32:
		asmTypeInstruction = "real4"
	case orth_types.StdF64:
		asmTypeInstruction = "real8"
	default:
		fmt.Fprintf(os.Stderr, "ivalid type of %q\n", operand.SymbolName)
		os.Exit(1)
	}
	return asmTypeInstruction
}

func StringToByteRep(s string, endWithNullByte bool) (lietralValue string) {
	// Unquote/unescape the string
	unquoted, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		panic(err)
	}
	// convert to a byte array so we can convert each byte to a string representation
	unquotedBytes := []byte(unquoted)

	// allocate the buffer
	unquotedBF := make([]string, len(unquotedBytes)+1)
	if endWithNullByte {
		unquotedBF[len(unquotedBF)-1] = "0" // add null to the end
	}

	for i, byt := range unquotedBytes {
		unquotedBF[i] = fmt.Sprint(byt)
	}

	lietralValue = strings.Join(unquotedBF, ",")
	if !endWithNullByte {
		lietralValue = lietralValue[:len(lietralValue)-1]
	}
	return
}

func VarValueToAsmSyntax(operand orth_types.Operand, endWithNullByte bool) string {
	var lietralValue string
	switch operand.SymbolName {
	case orth_types.StdSTR:
		lietralValue = StringToByteRep(operand.Operand, endWithNullByte)
	default:
		lietralValue = operand.Operand
	}
	return lietralValue
}

func MangleVarName(o orth_types.Operation) string {
	var memType string

	if o.Instruction == orth_types.InstructionVar || o.Operator.SymbolName == orth_types.StdVar {
		memType = "Var"
	} else if o.Instruction == orth_types.InstructionConst || o.Operator.SymbolName == orth_types.StdConst {
		memType = "Const"
	} else {
		panic(fmt.Errorf("invalid operation on type %d", o.Instruction))
	}

	return fmt.Sprintf("%s@%s@%s", o.Context.Name, memType, o.Operator.Operand)
}

func BuildVarDataSeg(variable orth_types.Operation) string {
	variableValue := variable.Links["variable_value"].Operator
	return fmt.Sprintf("%s %s %s",
		MangleVarName(variable),
		VarTypeToAsmType(variableValue),
		VarValueToAsmSyntax(variableValue, true))
}
