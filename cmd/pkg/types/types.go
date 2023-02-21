package orthtypes

import (
	"reflect"
)

const FileType = "orth"

const (
	StdTrue  = "1"
	StdFalse = "0"
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
	PrimitiveEND     = "end"
	PrimitiveVOID    = "void"
	PrimitiveRNT     = "rnt"
	PrimitiveMem     = "mem"
	PrimitiveType    = "type"
	PrimitiveConst   = "const"
	PrimitiveHold    = "hold"
	PrimitiveProc    = "proc"
	PrimitiveIn      = "in"
	PrimitiveInvalid = ""
	Bitwise          = "bitwise"
)

const (
	INTS        = "ints"
	FLOATS      = "floats"
	STRING      = "string"
	BOOL        = "bool"
	VOID        = "void"
	RNT         = "rnt"
	ADDR        = "address"
	RNGABL      = "rangeable"
	MEM         = "mem"
	TYPE        = "type"
	INVALIDTYPE = ""
)

type Pair[T1, T2 any] struct {
	VarName  T1
	VarValue T2
}

type Operation struct {
	Instruction int
	Operand     Operand
	Context     string
	RefBlock    int
}

type Operand struct {
	VarType string
	Operand string
}

type OutOfOrder struct {
	Vars chan Pair[Operation, Operand]
}

type SliceOf[T comparable] struct {
	Slice *[]T
}

type File[T comparable] struct {
	Name      string
	CodeBlock T
}

// UpdateCodeReference takes an argument of type string and then updated the current codeblock to the passed one
func (f *File[T]) UpdateCodeReference(codeBlock T) {
	isString := reflect.TypeOf(codeBlock).Kind() == reflect.String
	if !isString {
		panic("cannot have a non string as a codeblock")
	}

	f.CodeBlock = codeBlock
}

// IsValidType checks whenever a variable has a know or unknow type
func (o Operation) IsValidType() bool {
	return GlobalTypes[TYPE][o.Operand.VarType] != "" ||
		GlobalTypes[INTS][o.Operand.VarType] != "" ||
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
	Token    string
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

	GlobalTypes[INTS][PrimitiveI64] = "i64"
	GlobalTypes[INTS][PrimitiveI32] = "i32"
	GlobalTypes[INTS][PrimitiveI16] = "i16"
	GlobalTypes[INTS][PrimitiveI8] = "i8"
	GlobalTypes[INTS][PrimitiveInt] = "i"

	GlobalTypes[FLOATS] = make(map[string]string, 0)
	GlobalTypes[FLOATS][PrimitiveF64] = "f64"
	GlobalTypes[FLOATS][PrimitiveF32] = "f32"

	GlobalTypes[STRING] = make(map[string]string, 0)
	GlobalTypes[STRING][PrimitiveSTR] = "s"

	GlobalTypes[BOOL] = make(map[string]string, 0)
	GlobalTypes[BOOL][PrimitiveBOOL] = "b"

	GlobalTypes[VOID] = make(map[string]string)
	GlobalTypes[VOID][PrimitiveVOID] = "v"

	GlobalTypes[RNT] = make(map[string]string, 0)
	GlobalTypes[RNT][PrimitiveRNT] = "rnt"
	GlobalTypes[RNT][ADDR] = "address"
	GlobalTypes[RNT][RNGABL] = "rangeable"

	GlobalTypes[MEM] = make(map[string]string, 0)
	GlobalTypes[MEM][PrimitiveMem] = "rnt"

	GlobalTypes[TYPE] = make(map[string]string, 0)
	GlobalTypes[TYPE]["type"] = "type"
	GlobalTypes[TYPE]["var"] = "var"
	GlobalTypes[TYPE]["hold"] = "hold"
}
