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
	PrimitiveVar     = "var"
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

type Context struct {
	Name         string
	Order        uint
	Parent       *Context
	Declarations []string
	InnerContext []*Context
}

type Pair[T1, T2 any] struct {
	Left  T1
	Right T2
}

type Operation struct {
	Instruction Instruction
	Operator    Operand
	Context     *Context
	RefBlock    int
}

func (op *Operation) IsString() bool {
	_, isString := GlobalTypes[STRING][op.Operator.SymbolName]
	return isString
}

func (op *Operation) IsNumeric() bool {
	_, isInt := GlobalTypes[INTS][op.Operator.SymbolName]
	_, isFloat := GlobalTypes[FLOATS][op.Operator.SymbolName]
	return isInt || isFloat
}

func (op *Operation) IsInt() bool {
	_, ok := GlobalTypes[INTS][op.Operator.SymbolName]
	return ok
}

func (op *Operation) IsFloat() bool {
	_, ok := GlobalTypes[FLOATS][op.Operator.SymbolName]
	return ok
}

func (op *Operation) IsFloat64() bool {
	return op.Operator.SymbolName == PrimitiveF64
}

func (op *Operation) IsFloat32() bool {
	return op.Operator.SymbolName == PrimitiveF32
}

type Operand struct {
	SymbolName string
	Operand    string
}

type SliceOf[T comparable] struct {
	Slice *[]T
}

type File[T comparable] struct {
	Name      string
	CodeBlock T
}

func (ctx *Context) MountFullLengthContext(name string) string {
	if ctx == nil {
		return name
	}
	name += ctx.Parent.MountFullLengthContext(name) + "_" + ctx.Name
	return name
}

func (ctx *Context) HasVariableDeclaredInOrAbove(variable string) bool {
	for ctx != nil {
		for _, v := range ctx.Declarations {
			if v == variable {
				return true
			}
		}
		ctx = ctx.Parent
	}
	return false
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
	return GlobalTypes[TYPE][o.Operator.SymbolName] != "" ||
		GlobalTypes[INTS][o.Operator.SymbolName] != "" ||
		GlobalTypes[FLOATS][o.Operator.SymbolName] != "" ||
		GlobalTypes[STRING][o.Operator.SymbolName] != "" ||
		GlobalTypes[BOOL][o.Operator.SymbolName] != "" ||
		GlobalTypes[VOID][o.Operator.SymbolName] != "" ||
		GlobalTypes[RNT][o.Operator.SymbolName] != "" ||
		GlobalTypes[MEM][o.Operator.SymbolName] != ""
}

func (o Operand) GrabRootType() string {
	var ret string
	switch {
	case GlobalTypes[INTS][o.SymbolName] != INVALIDTYPE:
		ret = INTS
	case GlobalTypes[STRING][o.SymbolName] != INVALIDTYPE:
		ret = STRING
	case GlobalTypes[FLOATS][o.SymbolName] != INVALIDTYPE:
		ret = FLOATS
	case GlobalTypes[RNT][o.SymbolName] != INVALIDTYPE:
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
