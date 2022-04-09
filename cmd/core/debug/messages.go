package debug

const (
	DefaultRuntimeException = "Fatal error while executing program\n"
	InvalidBoolType         = "Non bool type used as bool\n"
	//OFMException            = "Out of memory in %q storage\n"
	InvalidTypeForIndex = "Can not use a non %q as index\n"
	InvalidInstruction  = "Invalid instruction of type '%d'\n"
)
