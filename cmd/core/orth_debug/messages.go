package orth_debug

import (
	"errors"
	"fmt"
)

const (
	DefaultRuntimeException    = "RNT_ERR: fatal error while executing program"
	InvalidBoolType            = "RNT_ERR: non bool type used as bool"
	InvalidTypeForIndex        = "RNT_ERR: can not use a non %q as index"
	InvalidInstruction         = "RNT_ERR: invalid instruction of type '%d'"
	InvalidTypeForInstruction  = "RNT_ERR: invalid type of '%s' for instruction of type '%s'"
	VariableUndefined          = "RNT_ERR: variable %q does not exist"
	StackUnderFlow             = "RNT_ERR: stack underflow!"
	UndefinedToken             = "RNT_ERR: undefined token %q"
	UndefinedFunction          = "RNT_ERR: undefined function %q"
	StrangeUseOfVariable       = "RNT_ERR: a variable of type %q can not be used in %q statements"
	IndexOutOfBounds           = "RNT_ERR: the index %d is out of bounds [%d, %d]"
	InvalidUsageOfTokenOutside = "COMP_ERR: The token %q can only be used inside a %q context, rigth now it is been used in %q"
)

const commomFileSpecificationStruct = "in %q at line: %d colum: %d\n"

const (
	ORTH_ERR_01 = "[ERROR] Undefined/unknow token %q " + commomFileSpecificationStruct
	ORTH_ERR_02 = "[ERROR] Redeclaration of %q -> %q " + commomFileSpecificationStruct
	ORTH_ERR_03 = "[ERROR] Redeclaration of %q -> %q in block %q"
	ORTH_ERR_04 = "[ERROR] Invalid operation of type %q\nDetails:\n%s"
	ORTH_ERR_05 = "[ERROR] The instruction of type %q requires a parameter of type %q, but found token %q\n\t" + commomFileSpecificationStruct
	ORTH_ERR_06 = "[ERROR] A procedure can only have a maximum of %d but found %d\n\t" + commomFileSpecificationStruct
	ORTH_ERR_07 = "[ERROR] Syntax error. %s"
	ORTH_ERR_08 = "[ERROR] Instruction %q requires: (%q, %q). But found: (%q, %q)\n"
	ORTH_ERR_09 = "[ERROR] Stack underflow!. Instruction %q requires values that are not part of the stack!\n"
	ORTH_ERR_10 = "[ERROR] The instruction of type %q requires a parameter of type %q, but found %q\n"
	ORTH_ERR_11 = "[ERROR] Variable %q is undefined for instruction %q\n"
	ORTH_ERR_12 = "[ERROR] Invalid type %q %s" + commomFileSpecificationStruct
	ORTH_ERR_13 = "[ERROR] Incorrect number of arguments for instruction %q, required '%d' and got '%d' " + commomFileSpecificationStruct
	ORTH_ERR_14 = "[ERROR] Incorrect number of arguments for instruction %q, required '%s' and got '%d' " + commomFileSpecificationStruct
	ORTH_ERR_15 = "[ERROR] Could not find include file %q on paths"
)

const (
	ORTH_WARN_01 = "[WARN] Performin operation %q on values with distinct types (%q, %q)\n"
)

func BuildMessage(message string, params ...interface{}) string {
	return fmt.Sprintf(message, params...)
}

func BuildErrorMessage(message string, params ...interface{}) error {
	return errors.New(BuildMessage(message, params...))
}
