package debug

const (
	DefaultRuntimeException   = "RNT_ERR: fatal error while executing program"
	InvalidBoolType           = "RNT_ERR: non bool type used as bool"
	InvalidTypeForIndex       = "RNT_ERR: can not use a non %q as index"
	InvalidInstruction        = "RNT_ERR: invalid instruction of type '%d'"
	InvalidTypeForInstruction = "RNT_ERR: invalid type of '%s' for instruction of type '%s'"
	StackUnderFlow            = "RNT_ERR: stack underflow!"
	UndefinedToken            = "VM_ERR: undefined token %q"
)
