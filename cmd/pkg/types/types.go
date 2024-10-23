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
	StdI64     string = "i64"
	StdI32     string = "i32"
	StdI16     string = "i16"
	StdI8      string = "i8"
	StdINT     string = "i"
	StdF64     string = "f64"
	StdF32     string = "f32"
	StdSTR     string = "s"
	StdBOOL    string = "b"
	StdINVALID string = ""
)

// keywords
const (
	StdPlus          string = "+"
	StdMinus         string = "-"
	StdMult          string = "*"
	StdDiv           string = "/"
	StdEquals        string = "=="
	StdNotEquals     string = "<>"
	StdLowerThan     string = "<"
	StdGreaterThan   string = ">"
	StdMod           string = "%"
	StdEND           string = "end"
	StdVOID          string = "void"
	StdRNT           string = "rnt"
	StdParam         string = "param"
	StdMem           string = "mem"
	StdType          string = "type"
	StdConst         string = "const"
	StdVar           string = "var"
	StdHold          string = "hold"
	StdProc          string = "proc"
	StdIn            string = "in"
	StdIf            string = "if"
	StdElse          string = "else"
	StdOver          string = "over"
	Std2Dup          string = "2dup"
	StdDup           string = "dup"
	StdWhile         string = "while"
	StdLeftShift     string = "lshift"
	StdRightShift    string = "rshift"
	StdLogicalAnd    string = "land"
	StdLogicalOr     string = "lor"
	StdDo            string = "do"
	StdDrop          string = "drop"
	StdSwap          string = "swap"
	StdStore         string = "."
	StdLoad          string = ","
	StdCall          string = "call"
	StdLoadAndStay   string = ",!"
	StdInvoke        string = "invoke"
	StdProcOutParams string = "--"
	StdProcInParams  string = ":"
	StdAddress       string = "addr"
	StdBitwise       string = "bitwise"
)

// builtin functions/symbols
const (
	StdPutUint   string = "putui"
	StdPutStr    string = "puts"
	StdSetNumber string = "set_number"
	StdSetStr    string = "set_string"
	StdDumpMem   string = "dump_mem"
	StdPutChar   string = "put_char"
	StdDeref     string = "deref"
	StdExit      string = "exit"
	StdAlloc     string = "alloc"
	StdFree      string = "free"
)

// some shit I don't remember
const (
	INTS        string = "ints"
	FLOATS      string = "floats"
	STRING      string = "string"
	BOOL        string = "bool"
	VOID        string = "void"
	RNT         string = "rnt"
	ADDR        string = "address"
	RNGABL      string = "rangeable"
	MEM         string = "mem"
	TYPE        string = "type"
	INVALIDTYPE string = ""
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

func (op *Operation) PrioritizeAddress() (int, error) {
	priorities := instructionJumpAddressPriority[op.Instruction]
	for _, instruction := range priorities {
		jumpAddress, ok := op.Addresses[instruction]
		if ok {
			return jumpAddress, nil
		}
	}
	return 0, errors.New("no addresses")
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
	GlobalTypes[INTS][StdAddress] = "addr"

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
