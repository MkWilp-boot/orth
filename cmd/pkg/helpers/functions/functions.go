package functions

import (
	"errors"
	"fmt"
	"orth/cmd/core/orth_debug"
	"orth/cmd/pkg/helpers"
	orth_types "orth/cmd/pkg/types"
	"os"
	"strconv"
	"strings"
)

func panicErrIfNonNil(err error) {
	if err != nil {
		panic(err)
	}
}

func DissectRangeAsInt(o1 orth_types.Operand) (int, int) {
	opStart, opEnd := DissectRange(o1)
	iStart, err := strconv.Atoi(opStart.Operand)
	panicErrIfNonNil(err)
	iEnd, err := strconv.Atoi(opEnd.Operand)
	panicErrIfNonNil(err)
	return iStart, iEnd
}

func DissectRange(o1 orth_types.Operand) (orth_types.Operand, orth_types.Operand) {
	if o1.SymbolName != orth_types.RNGABL {
		panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, o1.SymbolName, "DissectRange"))
	}
	nums := strings.Split(o1.Operand, "|")
	start, _ := strconv.Atoi(nums[0])
	end, _ := strconv.Atoi(nums[1])

	return orth_types.Operand{
			SymbolName: orth_types.StdI32,
			Operand:    fmt.Sprint(start),
		}, orth_types.Operand{
			SymbolName: orth_types.StdI32,
			Operand:    fmt.Sprint(end),
		}
}

func CheckAsmType(flagValue string) (string, error) {
	available := []string{"nasm", "masm", "fasm"}
	for _, v := range available {
		if flagValue == v {
			return v, nil
		}
	}
	return "", errors.New("unsupported assembly type")
}

// TypesAreEqual checks if the compared types the same INNER-TYPE variant
func TypesAreEqual(opreands ...orth_types.Operand) bool {
	t := opreands[0].SymbolName
	equal := true

	for _, n := range opreands {
		equal = n.SymbolName == t
	}

	return equal
}

// GetSupersetType gets the supertype of a set of variables.
// Ex:
// If you pass a set of [i,i,i,i32] then the SPS will be i32 because
// i32 has the largest cap of all INT variants present in this context.
// More:
// | [i32, i32, i64] -> i64
// | [i,i,i] -> i
// | [f32, f64] -> f64
// | [s,s] -> s
// | [s, i] -> panic()
func GetSupersetType(opreands ...orth_types.Operand) string {
	switch {
	case helpers.IsInt(opreands[0]):
		return IntSupersetOfSlice(opreands...)
	case helpers.IsFloat(opreands[0]):
		return FloatSupersetOfSlice(opreands...)
	case opreands[0].SymbolName == orth_types.StdSTR:
		return orth_types.StdSTR

	default:
		panic("Invalid type")
	}
}

// ===================================
//	BASEDON
// ===================================

// ModBasedOnType Modules a set of numbers based on the set's type
func ModBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return ModIntegers
	case strings.Contains(superType, "f"):
		return ModFloats
	case strings.Contains(superType, orth_types.StdSTR):
		panic("Can not use 'orth_types.PrimitiveSTR' with '%' operation")
	default:
		panic("Invalid type")
	}
}

// SubBasedOnType subs a set of numbers based on the set's type
func SubBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return SubIntegers
	case strings.Contains(superType, "f"):
		return SubFloats
	case strings.Contains(superType, orth_types.StdSTR):
		panic("Can not apply '-' operation to a string value")
	default:
		panic("Invalid type")
	}
}

// DivBasedOnType divides a set of numbers based on the set's type
func DivBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return DivideIntegers
	case strings.Contains(superType, "f"):
		return DivideFloats
	case strings.Contains(superType, orth_types.StdSTR):
		panic("Can not apply '/' operation to a string value")
	default:
		panic("Invalid type")
	}
}

// MultBasedOnType multiplies a set of numbers based on the set's type
func MultBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return MultplyIntegers
	case strings.Contains(superType, "f"):
		return MultplyFloats
	case strings.Contains(superType, orth_types.StdSTR):
		panic("Can not apply '*' operation to a string value")
	default:
		panic("Invalid type")
	}
}

// EqualBasedOnType compare a set of numbers based on the set's type
func EqualBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return EqualInts
	case strings.Contains(superType, "f"):
		return EqualFloats
	case strings.Contains(superType, orth_types.StdSTR):
		return EqualString
	default:
		panic("Invalid type")
	}
}

// LowerBasedOnType compare a set of numbers based on the set's type
func LowerBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return LowerThanInts
	case strings.Contains(superType, "f"):
		return LowerThanFloats
	case strings.Contains(superType, orth_types.StdSTR):
		panic("Can not apply '<' in a string literal")
	default:
		panic("Invalid type")
	}
}

// GreaterBasedOnType compare a set of numbers based on the set's type
func GreaterBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return GreaterThanInts
	case strings.Contains(superType, "f"):
		return GreaterThanFloats
	case strings.Contains(superType, orth_types.StdSTR):
		panic("Can not apply '>' in a string literal")
	default:
		panic("Invalid type")
	}
}

// NotEqualBasedOnType compare a set of numbers based on the set's type
func NotEqualBasedOnType(superType string) func(string, orth_types.Operand, orth_types.Operand) orth_types.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return DiffInts
	case strings.Contains(superType, "f"):
		return DiffFloats
	case strings.Contains(superType, orth_types.StdSTR):
		return DiffString
	default:
		panic("Invalid type")
	}
}

// EqualInts compare a set of integers
func EqualInts(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}

	o2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	var res string
	if o1 == o2 {
		res = orth_types.StdTrue
	} else {
		res = orth_types.StdFalse
	}

	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    res,
	}
}

// DiffInts compare a set of integers
func DiffInts(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}

	o2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	var op string
	if o1 != o2 {
		op = orth_types.StdTrue
	} else {
		op = orth_types.StdFalse
	}

	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    op,
	}
}

// LowerThanInts compare a set of integers
func LowerThanInts(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}

	o2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	var op string
	if o1 < o2 {
		op = orth_types.StdTrue
	} else {
		op = orth_types.StdFalse
	}
	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    op,
	}
}

// GreaterThanInts compare a set of integers
func GreaterThanInts(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}

	o2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	var op string
	if o1 > o2 {
		op = orth_types.StdTrue
	} else {
		op = orth_types.StdFalse
	}
	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    op,
	}
}

// MultplyIntegers multiplies a set of integers
func MultplyIntegers(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var mult string
	switch superType {
	case orth_types.StdI64:
		mult = MultI64(n1, n2)
	case orth_types.StdI32:
		mult = MultI32(n1, n2)
	case orth_types.StdI16:
		mult = MultI16(n1, n2)
	case orth_types.StdI8:
		mult = MultI8(n1, n2)
	case orth_types.StdINT:
		mult = MultI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    mult,
	}
}

func BitwiseAnd(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	if !helpers.IsInt(n1) || !helpers.IsInt(n2) {
		fmt.Fprintln(os.Stderr, "cannot perform 'logical and' on values that are not integers")
		os.Exit(1)
	}

	left := helpers.ToInt(n1)
	right := helpers.ToInt(n2)

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    fmt.Sprint(left & right),
	}
}

func BitwiseOr(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	if _, ok := orth_types.GlobalTypes[orth_types.INTS][superType]; !ok {
		fmt.Fprintln(os.Stderr, "cannot perform 'logical or' on values that are not integers")
		os.Exit(1)
	}

	left := helpers.ToInt(n1)
	right := helpers.ToInt(n2)

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    fmt.Sprint(left | right),
	}
}

// DivideIntegers divides a set of integers
func DivideIntegers(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var div string
	switch superType {
	case orth_types.StdI64:
		div = DivI64(n1, n2)
	case orth_types.StdI32:
		div = DivI32(n1, n2)
	case orth_types.StdI16:
		div = DivI16(n1, n2)
	case orth_types.StdI8:
		div = DivI8(n1, n2)
	case orth_types.StdINT:
		div = DivI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    div,
	}
}

// Left and shift functions are basiclly copies of each other, don't ask

func LeftShiftFloat(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	shiftAmount, err := strconv.Atoi(n1.Operand)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bitSize := 32

	if n2.SymbolName == orth_types.StdF64 {
		bitSize = 64
	}
	floatValue, err := strconv.ParseFloat(n2.Operand, bitSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	decimalDigits := (floatValue - float64(int(floatValue)))
	// shift the int part and sum with decimals to make another float (silly me)
	floatValue = float64(int(floatValue)<<shiftAmount) + decimalDigits

	precision := 0
	if index := strings.Index(n2.Operand, "."); index != -1 {
		precision = len(n2.Operand[index+1:])
	}

	return orth_types.Operand{
		SymbolName: n2.SymbolName,
		Operand:    strconv.FormatFloat(floatValue, 'f', precision, bitSize),
	}
}

func LeftShiftInt(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	shiftAmount, _ := strconv.Atoi(n1.Operand)
	intValue, err := strconv.Atoi(n2.Operand)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	intValue = intValue << shiftAmount

	return orth_types.Operand{
		SymbolName: n2.SymbolName,
		Operand:    strconv.Itoa(intValue),
	}
}

func RightShiftFloat(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	shiftAmount, err := strconv.Atoi(n1.Operand)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bitSize := 32

	if n2.SymbolName == orth_types.StdF64 {
		bitSize = 64
	}
	floatValue, err := strconv.ParseFloat(n2.Operand, bitSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	decimalDigits := (floatValue - float64(int(floatValue)))
	// shift the int part and sum with decimals to make another float (silly me)
	floatValue = float64(int(floatValue)>>shiftAmount) + decimalDigits

	precision := 0
	if index := strings.Index(n2.Operand, "."); index != -1 {
		precision = len(n2.Operand[index+1:])
	}

	return orth_types.Operand{
		SymbolName: n2.SymbolName,
		Operand:    strconv.FormatFloat(floatValue, 'f', precision, bitSize),
	}
}

func RightShiftInt(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	shiftAmount, _ := strconv.Atoi(n1.Operand)
	intValue, err := strconv.Atoi(n2.Operand)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	intValue = intValue >> shiftAmount

	return orth_types.Operand{
		SymbolName: n2.SymbolName,
		Operand:    strconv.Itoa(intValue),
	}
}

// SumIntegers sums up a set of integers
func SumIntegers(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var sum string
	var err error

	switch superType {
	case orth_types.StdAddress:
		sum, err = SumAddress(n1, n2)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case orth_types.StdI64:
		sum = SumI64(n1, n2)
	case orth_types.StdI32:
		sum = SumI32(n1, n2)
	case orth_types.StdI16:
		sum = SumI16(n1, n2)
	case orth_types.StdI8:
		sum = SumI8(n1, n2)
	case orth_types.StdINT:
		sum = SumI(n1, n2)
	default:
		fmt.Fprintln(os.Stderr, "not an integer")
		os.Exit(1)
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    sum,
	}
}

// ModIntegers modules a set of integers
func ModIntegers(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var sum string
	switch superType {
	case "i64":
		sum = ModI64(n1, n2)
	case "i32":
		sum = ModI32(n1, n2)
	case "i16":
		sum = ModI16(n1, n2)
	case "i8":
		sum = ModI8(n1, n2)
	case "i":
		sum = ModI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    sum,
	}
}

// SubIntegers subtracts a set of integers
func SubIntegers(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var sub string
	switch superType {
	case orth_types.StdI64:
		sub = SubI64(n1, n2)
	case orth_types.StdI32:
		sub = SubI32(n1, n2)
	case orth_types.StdI16:
		sub = SubI16(n1, n2)
	case orth_types.StdI8:
		sub = SubI8(n1, n2)
	case orth_types.StdINT:
		sub = SubI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    sub,
	}
}

// IntSupersetOfSlice gets the super type of a slice of integers
func IntSupersetOfSlice(opreands ...orth_types.Operand) string {
	for _, v := range opreands {
		if v.SymbolName == orth_types.StdI64 {
			return v.SymbolName
		}
	}

	for _, v := range opreands {
		if v.SymbolName == orth_types.StdI32 {
			return v.SymbolName
		}
	}

	for _, v := range opreands {
		if v.SymbolName == orth_types.StdI16 {
			return v.SymbolName
		}
	}

	for _, v := range opreands {
		if v.SymbolName == orth_types.StdI8 {
			return v.SymbolName
		}
	}

	return orth_types.StdINT
}

// ModFloats divides a set of floays
func ModFloats(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var sum string
	switch superType {
	case "f64":
		sum = ModF64(n1, n2)
	case "f32":
		sum = ModF32(n1, n2)
	default:
		panic("Not an float")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    sum,
	}
}

// EqualFloats compare a set of floats
func EqualFloats(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	o2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	var res string
	if o1 == o2 {
		res = "1"
	} else {
		res = "0"
	}

	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    res,
	}
}

// DiffFloats compare a set of floats
func DiffFloats(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	o2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    fmt.Sprintf("%v", o1 != o2),
	}
}

// LowerThanFloats compare a set of floats
func LowerThanFloats(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	o2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    fmt.Sprintf("%v", o1 < o2),
	}
}

// GreaterThanFloats compare a set of floats
func GreaterThanFloats(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	o1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	o2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	var op string
	if o1 > o2 {
		op = "1"
	} else {
		op = "0"
	}
	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    op,
	}
}

// MultplyFloats multiplies a set of floats
func MultplyFloats(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var mult string
	switch superType {
	case "f64":
		mult = MultF64(n1, n2)
	case "f32":
		mult = MultF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    mult,
	}
}

// DivideFloats divides a set of floats
func DivideFloats(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var div string
	switch superType {
	case "f64":
		div = DivF64(n1, n2)
	case "f32":
		div = DivF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    div,
	}
}

// SumFloats sums up a set of floats
func SumFloats(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var sum string
	switch superType {
	case "f64":
		sum = SumF64(n1, n2)
	case "f32":
		sum = SumF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    sum,
	}
}

// SubFloats subtract a set of floats
func SubFloats(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	var sub string
	switch superType {
	case "f64":
		sub = SubF64(n1, n2)
	case "f32":
		sub = SubF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orth_types.Operand{
		SymbolName: superType,
		Operand:    sub,
	}
}

// FloatSupersetOfSlice gets the super type of a slice of floats
func FloatSupersetOfSlice(opreands ...orth_types.Operand) string {
	for _, v := range opreands {
		if v.SymbolName == "f64" {
			return v.SymbolName
		}
	}

	for _, v := range opreands {
		if v.SymbolName == "f32" {
			return v.SymbolName
		}
	}
	panic("Invalid type")
}

// DiffString check if two string are different
func DiffString(_ string, n1, n2 orth_types.Operand) orth_types.Operand {
	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    fmt.Sprintf("%v", n1.Operand != n2.Operand),
	}
}

// EqualString checks if two variables have the same Operand
func EqualString(_ string, o1, o2 orth_types.Operand) orth_types.Operand {

	var res string
	if o1 == o2 {
		res = "1"
	} else {
		res = "0"
	}

	return orth_types.Operand{
		SymbolName: orth_types.StdBOOL,
		Operand:    res,
	}
}

// ConcatPrimitiveSTR concats two string into 1
func ConcatPrimitiveSTR(superType string, n1, n2 orth_types.Operand) orth_types.Operand {
	return orth_types.Operand{
		SymbolName: superType,
		Operand:    n2.Operand + n1.Operand,
	}
}

// ToString converts any Operand's operand to it's string literal
func ToString(n1 orth_types.Operand) orth_types.Operand {
	return orth_types.Operand{
		SymbolName: "s",
		Operand:    n1.Operand,
	}
}
