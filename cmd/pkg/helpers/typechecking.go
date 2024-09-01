package helpers

import (
	"errors"
	"fmt"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"strconv"
)

func IsNumeric(t string) bool {
	return IsInt(t) || IsFloat(t)
}

func ToAddress(op orth_types.Operand) (int, bool) {
	address, err := strconv.Atoi(op.Operand)
	return address, IsInt(op.SymbolName) && err == nil
}

func IsInt(t string) bool {
	_, ok := orth_types.GlobalTypes[orth_types.INTS][t]
	return ok
}

func IsFloat(t string) bool {
	return orth_types.GlobalTypes[orth_types.FLOATS][t] != ""
}

func IsBool(t string) bool {
	return orth_types.GlobalTypes[orth_types.BOOL][t] != ""
}

func IsString(t string) bool {
	return orth_types.GlobalTypes[orth_types.STRING][t] != ""
}

// SameBaseType checks if the 2 variables have the same base type.
// Ex: INT-INT, FLOAT-FLOAT, STRING-INT
func SameBaseType(operands ...orth_types.Operand) error {
	if len(operands) == 0 {
		return errors.New("no operands to loop")
	}
	baseOperand := operands[0]
	for _, operand := range operands {
		if operand.Operand != baseOperand.Operand {
			return fmt.Errorf("mismatch types %q and %q", baseOperand.Operand, operand.Operand)
		}
	}

	return nil
}

func OperatingOnEqualTypes(operations ...orth_types.Operation) error {
	if len(operations) == 0 {
		return errors.New("no operands to loop")
	}
	baseOperand := operations[0]
	for _, operand := range operations {
		if operand.Operator.SymbolName != baseOperand.Operator.SymbolName {
			return fmt.Errorf("mismatch types %q and %q", baseOperand.Operator.SymbolName, operand.Operator.SymbolName)
		}
	}

	return nil
}

// ===================================
//
//	CONVERSIONS
//
// ===================================
func ToInt(o orth_types.Operand) int {
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

	if o.SymbolName == orth_types.ADDR {
		n, err := strconv.Atoi(o.Operand)
		if err != nil {
			panic(orth_debug.DefaultRuntimeException)
		}
		return n
	}

	panic("Non INT or FLOAT variant been used in 'to int' operation!")
}
