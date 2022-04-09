package orthtypes

const FileType = "orth"

const (
	StdTrue  = "true"
	StdFalse = "false"
)

const (
	PrimitiveI64     = "i64"
	PrimitiveI32     = "i32"
	PrimitiveI16     = "i16"
	PrimitiveI8      = "i8"
	PrimitiveInt     = "i"
	PrimitiveF64     = "f64"
	PrimitiveF32     = "f32"
	PrimitiveSTR     = "s"
	PrimitiveBOOL    = "b"
	PrimitiveVOID    = "void"
	PrimitiveRNT     = "rnt"
	PrimitiveMem     = "mem"
	PrimitiveInvalid = ""
)

const (
	INTS        = "ints"
	FLOATS      = "floats"
	STRING      = "string"
	BOOL        = "bool"
	VOID        = "void"
	RNT         = "rnt"
	MEM         = "mem"
	INVALIDTYPE = ""
)

type Operation struct {
	Instruction int
	Operand     Operand
	RefBlock    int
}

type Operand struct {
	VarType string
	Operand string
}

func (o Operation) IsValidType() bool {
	return GlobalTypes[INTS][o.Operand.VarType] != "" ||
		GlobalTypes[FLOATS][o.Operand.VarType] != "" ||
		GlobalTypes[STRING][o.Operand.VarType] != "" ||
		GlobalTypes[BOOL][o.Operand.VarType] != "" ||
		GlobalTypes[VOID][o.Operand.VarType] != "" ||
		GlobalTypes[RNT][o.Operand.VarType] != "" ||
		GlobalTypes[MEM][o.Operand.VarType] != ""
}

func (o Operand) GrabRootType() string {
	var ret string
	switch {
	case GlobalTypes[INTS][o.VarType] != INVALIDTYPE:
		ret = INTS
	case GlobalTypes[STRING][o.VarType] != INVALIDTYPE:
		ret = STRING
	case GlobalTypes[FLOATS][o.VarType] != INVALIDTYPE:
		ret = FLOATS
	case GlobalTypes[RNT][o.VarType] != INVALIDTYPE:
		ret = RNT
	}
	return ret
}

type Vec2DString struct {
	Index    int
	ValidPos bool
	Content  string
}

type StringEnum struct {
	Index   int
	Content Vec2DString
}

type (
	Type   map[string]string
	Ints   Type
	Floats Type
	String Type
)

var GlobalTypes map[string]Type

func init() {
	GlobalTypes = make(map[string]Type, 0)
	GlobalTypes[INTS] = make(map[string]string, 0)

	GlobalTypes[INTS][PrimitiveI64] = "int64"
	GlobalTypes[INTS][PrimitiveI32] = "int32"
	GlobalTypes[INTS][PrimitiveI16] = "int16"
	GlobalTypes[INTS][PrimitiveI8] = "int8"
	GlobalTypes[INTS][PrimitiveInt] = "int"

	GlobalTypes[FLOATS] = make(map[string]string, 0)
	GlobalTypes[FLOATS][PrimitiveF64] = "float64"
	GlobalTypes[FLOATS][PrimitiveF32] = "float32"

	GlobalTypes[STRING] = make(map[string]string, 0)
	GlobalTypes[STRING][PrimitiveSTR] = "string"

	GlobalTypes[BOOL] = make(map[string]string, 0)
	GlobalTypes[BOOL][PrimitiveBOOL] = "bool"

	GlobalTypes[VOID] = make(map[string]string, 0)
	GlobalTypes[VOID][PrimitiveVOID] = "void"

	GlobalTypes[RNT] = make(map[string]string, 0)
	GlobalTypes[RNT][PrimitiveRNT] = "rnt"

	GlobalTypes[MEM] = make(map[string]string, 0)
	GlobalTypes[MEM][PrimitiveMem] = "rnt"
}
