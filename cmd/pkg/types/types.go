package orth_types

import (
	"errors"
	"reflect"
)

const FileType = "orth"

const (
	StdTrue  = "1"
	StdFalse = "0"
)

// types
const (
	StdI64     = "i64"
	StdI32     = "i32"
	StdI16     = "i16"
	StdI8      = "i8"
	StdINT     = "i"
	StdF64     = "f64"
	StdF32     = "f32"
	StdSTR     = "s"
	StdBOOL    = "b"
	StdINVALID = ""
)

// keywords
const (
	StdPlus          = "+"
	StdMinus         = "-"
	StdMult          = "*"
	StdDiv           = "/"
	StdEquals        = "=="
	StdNotEquals     = "<>"
	StdLowerThan     = "<"
	StdGreaterThan   = ">"
	StdMod           = "%"
	StdEND           = "end"
	StdVOID          = "void"
	StdRNT           = "rnt"
	StdParam         = "param"
	StdMem           = "mem"
	StdType          = "type"
	StdConst         = "const"
	StdVar           = "var"
	StdHold          = "hold"
	StdProc          = "proc"
	StdIn            = "in"
	StdIf            = "if"
	StdElse          = "else"
	StdOver          = "over"
	Std2Dup          = "2dup"
	StdDup           = "dup"
	StdWhile         = "while"
	StdLeftShift     = "lshift"
	StdRightShift    = "rshift"
	StdLogicalAnd    = "land"
	StdLogicalOr     = "lor"
	StdDo            = "do"
	StdDrop          = "drop"
	StdSwap          = "swap"
	StdStore         = "."
	StdLoad          = ","
	StdCall          = "call"
	StdLoadAndStay   = ",!"
	StdInvoke        = "invoke"
	StdProcOutParams = "--"
	StdProcInParams  = ":"
	StdBitwise       = "bitwise"
)

// builtin functions/symbols
const (
	StdPutUint   = "putui"
	StdPutStr    = "puts"
	StdSetNumber = "set_number"
	StdSetStr    = "set_string"
	StdDumpMem   = "dump_mem"
	StdPutChar   = "put_char"
	StdDeref     = "deref"
	StdExit      = "exit"
	StdAlloc     = "alloc"
	StdFree      = "free"
)

// some shit I don't remember
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

type ContextDeclaration struct {
	Name  string
	Index uint
}

type Context struct {
	Name          string
	Order         uint
	Parent        *Context
	Declarations  []ContextDeclaration
	InnerContexts []*Context
}

type Pair[T1, T2 any] struct {
	Left  T1
	Right T2
}

type Operation struct {
	Instruction Instruction
	Operator    Operand
	Context     *Context
	Links       map[string]Operation
	Addresses   map[Instruction]int
}

func (op *Operation) PrioritizeAddress() int {
	priorities := instructionJumpAddressPriority[op.Instruction]
	for _, instruction := range priorities {
		jumpAddress, ok := op.Addresses[instruction]
		if ok {
			return jumpAddress
		}
	}
	return -1
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
	return op.Operator.SymbolName == StdF64
}

func (op *Operation) IsFloat32() bool {
	return op.Operator.SymbolName == StdF32
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

func (ctx *Context) GetVaraible(variable string, program *Program) (*Operation, error) {
	for ctx != nil {
		for _, decls := range ctx.Declarations {
			if decls.Name == variable {
				return &program.Operations[decls.Index], nil
			}
		}
		ctx = ctx.Parent
	}
	return nil, errors.New("could not find variable, it's either out of scope or was not declared")
}

func (ctx *Context) GetNestedVariables(program *Program) ([]Operation, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}
	variables := make([]Operation, 0, len(ctx.Declarations))
	for _, variable := range ctx.Declarations {
		variables = append(variables, program.Operations[variable.Index])
	}
	for _, innerContext := range ctx.InnerContexts {
		nestedVariables, err := innerContext.GetNestedVariables(program)
		// unlikely to have an "inner context" set to nil, but let's be safe...
		if err == nil {
			variables = append(variables, nestedVariables...)
		}
	}
	return variables, nil
}

func (ctx *Context) HasVariableDeclaredInOrAbove(variable string) bool {
	for ctx != nil {
		for _, v := range ctx.Declarations {
			if v.Name == variable {
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

// IsValidTypeOp checks whenever a variable has a know or unknow type
func IsValidTypeOp(o *Operation) bool {
	return GlobalTypes[TYPE][o.Operator.SymbolName] != "" ||
		GlobalTypes[INTS][o.Operator.SymbolName] != "" ||
		GlobalTypes[FLOATS][o.Operator.SymbolName] != "" ||
		GlobalTypes[STRING][o.Operator.SymbolName] != "" ||
		GlobalTypes[BOOL][o.Operator.SymbolName] != "" ||
		GlobalTypes[VOID][o.Operator.SymbolName] != "" ||
		GlobalTypes[RNT][o.Operator.SymbolName] != "" ||
		GlobalTypes[MEM][o.Operator.SymbolName] != ""
}

// IsValidTypeSybl checks whenever a variable has a know or unknow type
func IsValidTypeSybl(s string) bool {
	return GlobalTypes[TYPE][s] != "" ||
		GlobalTypes[INTS][s] != "" ||
		GlobalTypes[FLOATS][s] != "" ||
		GlobalTypes[STRING][s] != "" ||
		GlobalTypes[BOOL][s] != "" ||
		GlobalTypes[VOID][s] != "" ||
		GlobalTypes[RNT][s] != "" ||
		GlobalTypes[MEM][s] != ""
}

func GrabType(o string) string {
	switch {
	case GlobalTypes[INTS][o] != INVALIDTYPE:
		return GlobalTypes[INTS][o]
	case GlobalTypes[STRING][o] != INVALIDTYPE:
		return GlobalTypes[STRING][o]
	case GlobalTypes[FLOATS][o] != INVALIDTYPE:
		return GlobalTypes[FLOATS][o]
	case GlobalTypes[RNT][o] != INVALIDTYPE:
		return GlobalTypes[RNT][o]
	default:
		return ""
	}
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
var instructionJumpAddressPriority map[Instruction][]Instruction

func init() {
	instructionJumpAddressPriority = make(map[Instruction][]Instruction)
	instructionJumpAddressPriority[InstructionIf] = []Instruction{InstructionElse, InstructionEnd}
	instructionJumpAddressPriority[InstructionElse] = []Instruction{InstructionEnd}

	GlobalTypes = make(map[string]Type, 0)
	GlobalTypes[INTS] = make(map[string]string, 0)

	GlobalTypes[INTS][StdI64] = "i64"
	GlobalTypes[INTS][StdI32] = "i32"
	GlobalTypes[INTS][StdI16] = "i16"
	GlobalTypes[INTS][StdI8] = "i8"
	GlobalTypes[INTS][StdINT] = "i"

	GlobalTypes[FLOATS] = make(map[string]string, 0)
	GlobalTypes[FLOATS][StdF64] = "f64"
	GlobalTypes[FLOATS][StdF32] = "f32"

	GlobalTypes[STRING] = make(map[string]string, 0)
	GlobalTypes[STRING][StdSTR] = "s"

	GlobalTypes[BOOL] = make(map[string]string, 0)
	GlobalTypes[BOOL][StdBOOL] = "b"

	GlobalTypes[VOID] = make(map[string]string)
	GlobalTypes[VOID][StdVOID] = "v"

	GlobalTypes[RNT] = make(map[string]string, 0)
	GlobalTypes[RNT][StdRNT] = "rnt"
	GlobalTypes[RNT][ADDR] = "address"
	GlobalTypes[RNT][RNGABL] = "rangeable"

	GlobalTypes[MEM] = make(map[string]string, 0)
	GlobalTypes[MEM][StdMem] = "rnt"

	GlobalTypes[TYPE] = make(map[string]string, 0)
	GlobalTypes[TYPE]["type"] = "type"
	GlobalTypes[TYPE]["var"] = "var"
	GlobalTypes[TYPE]["hold"] = "hold"
}
