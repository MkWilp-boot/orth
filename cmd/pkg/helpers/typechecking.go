package helpers

import (
	"fmt"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
	"strconv"
)

func IsInt(t string) bool {
	_, ok := orthtypes.GlobalTypes[orthtypes.INTS][t]
	return ok
}

func IsFloat(t string) bool {
	return orthtypes.GlobalTypes[orthtypes.FLOATS][t] != ""
}

func IsBool(t string) bool {
	return orthtypes.GlobalTypes[orthtypes.BOOL][t] != ""
}

func IsString(t string) bool {
	return orthtypes.GlobalTypes[orthtypes.STRING][t] != ""
}

// SameBaseType checks if the 2 variables have the same base type.
// Ex: INT-INT, FLOAT-FLOAT, STRING-INT
func SameBaseType(operands ...orthtypes.Operand) {
	if operands[0].GrabRootType() != operands[1].GrabRootType() {
		panic(fmt.Errorf("mismatch types! [%q - %q] and [%q - %q]", operands[0].Operand, operands[0].SymbolName, operands[1].Operand, operands[1].SymbolName))
	}
}

// ===================================
//
//	CONVERSIONS
//
// ===================================
func ToInt(o orthtypes.Operand) int {
	if IsInt(o.SymbolName) {
		n, err := strconv.Atoi(o.Operand)
		if err != nil {
			panic("Unknow error " + err.Error())
		}
		return n
	}

	if IsFloat(o.SymbolName) {
		f, err := strconv.ParseFloat(o.Operand, 64)
		if err != nil {
			panic("Unknow error " + err.Error())
		}
		return int(f)
	}

	if o.SymbolName == orthtypes.ADDR {
		n, err := strconv.Atoi(o.Operand)
		if err != nil {
			panic(orth_debug.DefaultRuntimeException)
		}
		return n
	}

	panic("Non INT or FLOAT variant been used in 'to int' operation!")
}
