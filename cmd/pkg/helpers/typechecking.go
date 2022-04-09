package helpers

import (
	"fmt"
	"strconv"
	orthtypes "t/cmd/pkg/types"
)

func IsInt(t string) bool {
	return orthtypes.GlobalTypes[orthtypes.INTS][t] != ""
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

func SameBaseType(operands ...orthtypes.Operand) {
	if operands[0].GrabRootType() != operands[1].GrabRootType() {
		panic(fmt.Errorf("Mismatch types! [%q - %q] and [%q - %q]", operands[0].Operand, operands[0].VarType, operands[1].Operand, operands[1].VarType))
	}
}

// ===================================
//	CONVERSIONS
// ===================================
func ToInt(o orthtypes.Operand) int {
	if IsInt(o.VarType) {
		n, err := strconv.Atoi(o.Operand)
		if err != nil {
			panic("Unknow error " + err.Error())
		}
		return n
	}

	if IsFloat(o.VarType) {
		f, err := strconv.ParseFloat(o.Operand, 64)
		if err != nil {
			panic("Unknow error " + err.Error())
		}
		return int(f)
	}

	panic("Non INT or FLOAT variant been used in 'to int' operation!")
}
