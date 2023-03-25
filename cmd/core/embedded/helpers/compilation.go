package embedded_helpers

import (
	"fmt"
	"log"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"runtime"
	"strconv"
	"strings"
)

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

func RetrieveProgramInfo(program orthtypes.Program, outOfOrder orthtypes.OutOfOrder, act func(*orthtypes.Program, orthtypes.Operation, int) []orthtypes.Pair[orthtypes.Operation, orthtypes.Operand]) {
	for i, operation := range program.Operations {
		retreive := act(&program, operation, i)
		for _, op := range retreive {
			outOfOrder.Vars <- op
		}
	}
	close(outOfOrder.Vars)
}

func GetVarsAndValues(program *orthtypes.Program, operation orthtypes.Operation, i int) []orthtypes.Pair[orthtypes.Operation, orthtypes.Operand] {
	retreive := make([]orthtypes.Pair[orthtypes.Operation, orthtypes.Operand], 0, cap(program.Operations))
	if (operation.Instruction == orthtypes.Var || operation.Instruction == orthtypes.Const) &&
		program.Operations[i-1].Instruction == orthtypes.Push {
		retreive = append(retreive, orthtypes.Pair[orthtypes.Operation, orthtypes.Operand]{
			VarName:  operation,
			VarValue: program.Operations[i-1].Operand,
		})
		program.Operations[i].Instruction = orthtypes.Skip
		program.Operations[i-1].Instruction = orthtypes.Skip
	}
	return retreive
}

func VarTypeToAsmType(operand orthtypes.Operand) string {
	var asmTypeInstruction string
	switch operand.VarType {
	case orthtypes.PrimitiveSTR:
		asmTypeInstruction = "db"
	case orthtypes.PrimitiveI8:
		asmTypeInstruction = "byte"
	case orthtypes.PrimitiveI16:
		asmTypeInstruction = "dw"
	case orthtypes.PrimitiveI32:
		asmTypeInstruction = "dd"
	case orthtypes.PrimitiveInt:
		fallthrough
	case orthtypes.PrimitiveI64:
		asmTypeInstruction = "dq"
	case orthtypes.PrimitiveF32:
		asmTypeInstruction = "real4"
	case orthtypes.PrimitiveF64:
		asmTypeInstruction = "real8"
	default:
		if strings.Contains(runtime.GOARCH, "64") {
			asmTypeInstruction = "dw"
		} else {
			asmTypeInstruction = "dd"
		}
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

func VarValueToAsmSyntax(operand orthtypes.Operand, endWithNullByte bool) string {
	var lietralValue string
	switch operand.VarType {
	case orthtypes.PrimitiveSTR:
		lietralValue = StringToByteRep(operand.Operand, endWithNullByte)
	default:
		lietralValue = operand.Operand
	}
	return lietralValue
}

func MangleVarName(o orthtypes.Operation) string {
	var memType string

	if o.Instruction == orthtypes.Var || o.Operand.VarType == orthtypes.PrimitiveVar {
		memType = "Var"
	} else if o.Instruction == orthtypes.Const || o.Operand.VarType == orthtypes.PrimitiveConst {
		memType = "Const"
	} else {
		panic(fmt.Errorf("invalid operation on type %d", o.Instruction))
	}

	return fmt.Sprintf("_@%s@%s@%s", o.Context, memType, o.Operand.Operand)
}

func BuildVarDataSeg(oVar orthtypes.Pair[orthtypes.Operation, orthtypes.Operand]) string {
	return fmt.Sprintf("%s %s %s", MangleVarName(oVar.VarName), VarTypeToAsmType(oVar.VarValue), VarValueToAsmSyntax(oVar.VarValue, true))
}
