package helpers

import (
	"errors"
	"fmt"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"strconv"
)

func IsNumeric(op orth_types.Operand) bool {
	return IsInt(op) || IsFloat(op)
}

func ToAddress(op orth_types.Operand) (int, bool) {
	address, err := strconv.Atoi(op.Operand)
	return address, IsInt(op) && err == nil
}

func IsInt(op orth_types.Operand) bool {
	_, ok := orth_types.GlobalTypes[orth_types.INTS][op.SymbolName]
	_, err := strconv.Atoi(op.Operand)
	return ok && err == nil
}

func IsFloat(op orth_types.Operand) bool {

	return orth_types.GlobalTypes[orth_types.FLOATS][op.SymbolName] != ""
}

func IsBool(op orth_types.Operand) bool {
	return orth_types.GlobalTypes[orth_types.BOOL][op.SymbolName] != "" &&
		(op.Operand == orth_types.StdTrue || op.Operand == orth_types.StdFalse)
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
	if IsInt(o) {
		n, err := strconv.Atoi(o.Operand)
		if err != nil {
			panic("Unknow error " + err.Error())
		}
		return n
	}

	if IsFloat(o) {
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
