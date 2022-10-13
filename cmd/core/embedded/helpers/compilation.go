package embedded_helpers

import (
	"fmt"
	orthtypes "orth/cmd/pkg/types"
	"strconv"
	"strings"
)

func RetreiveProgramInfo(program orthtypes.Program, ops chan<- orthtypes.Pair[orthtypes.Operand, orthtypes.Operand], act func(*orthtypes.Program, orthtypes.Operation, int) []orthtypes.Pair[orthtypes.Operand, orthtypes.Operand]) {
	for i, operation := range program.Operations {
		retreive := act(&program, operation, i)
		for _, op := range retreive {
			ops <- op
		}
	}
	close(ops)
}

func GetVarsAndValues(program *orthtypes.Program, operation orthtypes.Operation, i int) []orthtypes.Pair[orthtypes.Operand, orthtypes.Operand] {
	retreive := make([]orthtypes.Pair[orthtypes.Operand, orthtypes.Operand], 0, cap(program.Operations))
	if operation.Instruction == orthtypes.Var && program.Operations[i-1].Instruction == orthtypes.Push {
		retreive = append(retreive, orthtypes.Pair[orthtypes.Operand, orthtypes.Operand]{
			VarName:  operation.Operand,
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
	default:
		asmTypeInstruction = "dd"
	}
	return asmTypeInstruction
}

func VarValueToAsmSyntax(operand orthtypes.Operand) string {
	var lietralValue string
	switch operand.VarType {
	case orthtypes.PrimitiveSTR:
		// Unquote/unescape the string
		unquoted, err := strconv.Unquote(`"` + operand.Operand + `"`)
		if err != nil {
			panic(err)
		}
		// convert to a byte array so we can convert each byte to a string representation
		unquotedBytes := []byte(unquoted)

		// allocate the buffer
		unquotedBF := make([]string, len(unquotedBytes)+1, cap(unquotedBytes)+1)
		unquotedBF[len(unquotedBF)-1] = "0" // add null to the end

		for i, byt := range unquotedBytes {
			unquotedBF[i] = fmt.Sprint(byt)
		}

		lietralValue = strings.Join(unquotedBF, ",")
	default:
		lietralValue = operand.Operand
	}
	return lietralValue
}

func BuildVarDataSeg(oVar orthtypes.Pair[orthtypes.Operand, orthtypes.Operand]) string {
	return oVar.VarName.Operand + " " + VarTypeToAsmType(oVar.VarValue) + " " + VarValueToAsmSyntax(oVar.VarValue)
}
