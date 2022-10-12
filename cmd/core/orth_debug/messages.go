package orth_debug

const (
	DefaultRuntimeException   = "RNT_ERR: fatal error while executing program"
	InvalidBoolType           = "RNT_ERR: non bool type used as bool"
	InvalidTypeForIndex       = "RNT_ERR: can not use a non %q as index"
	InvalidInstruction        = "RNT_ERR: invalid instruction of type '%d'"
	InvalidTypeForInstruction = "RNT_ERR: invalid type of '%s' for instruction of type '%s'"
	VariableUndefined         = "RNT_ERR: variable %q does not exist"
	StackUnderFlow            = "RNT_ERR: stack underflow!"
	UndefinedToken            = "RNT_ERR: undefined token %q"
	UndefinedFunction         = "RNT_ERR: undefined function %q"
	StrangeUseOfVariable      = "RNT_ERR: a variable of type %q can not be used in %q statements"
	IndexOutOfBounds          = "RNT_ERR: the index %d is out of bounds [%d, %d]"
)
