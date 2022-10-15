package functions

import (
	"fmt"
	"orth/cmd/core/orth_debug"
	orthtypes "orth/cmd/pkg/types"
	"strconv"
	"strings"
)

func PanicErrIfNonNil(err error) {
	if err != nil {
		panic(err)
	}
}

func DissectRangeAsInt(o1 orthtypes.Operand) (int, int) {
	opStart, opEnd := DissectRange(o1)
	iStart, err := strconv.Atoi(opStart.Operand)
	PanicErrIfNonNil(err)
	iEnd, err := strconv.Atoi(opEnd.Operand)
	PanicErrIfNonNil(err)
	return iStart, iEnd
}

func DissectRange(o1 orthtypes.Operand) (orthtypes.Operand, orthtypes.Operand) {
	if o1.VarType != orthtypes.RNGABL {
		panic(fmt.Errorf(orth_debug.InvalidTypeForInstruction, o1.VarType, "DissectRange"))
	}
	nums := strings.Split(o1.Operand, "|")
	start, _ := strconv.Atoi(nums[0])
	end, _ := strconv.Atoi(nums[1])

	return orthtypes.Operand{
			VarType: orthtypes.PrimitiveI32,
			Operand: fmt.Sprint(start),
		}, orthtypes.Operand{
			VarType: orthtypes.PrimitiveI32,
			Operand: fmt.Sprint(end),
		}
}

func CheckAsmType(flagValue string) string {
	available := []string{"nasm", "masm", "fasm"}
	for _, v := range available {
		if flagValue == v {
			return v
		}
	}
	panic("unsupported assembly type")
}

// TypesAreEqual checks if the compared types the same INNER-TYPE variant
func TypesAreEqual(opreands ...orthtypes.Operand) bool {
	t := opreands[0].VarType
	equal := true

	for _, n := range opreands {
		equal = n.VarType == t
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
func GetSupersetType(opreands ...orthtypes.Operand) string {
	switch {
	case strings.Contains(opreands[0].VarType, "i") ||
		strings.Contains(opreands[0].VarType, "i8") ||
		strings.Contains(opreands[0].VarType, "i16") ||
		strings.Contains(opreands[0].VarType, "i32") ||
		strings.Contains(opreands[0].VarType, "i64"):
		return IntSupersetOfSlice(opreands...)
	case strings.Contains(opreands[0].VarType, "f32") ||
		strings.Contains(opreands[0].VarType, "f64"):
		return FloatSupersetOfSlice(opreands...)
	case strings.Contains(opreands[0].VarType, orthtypes.ADDR):
		return orthtypes.ADDR
	case opreands[0].VarType == orthtypes.PrimitiveSTR:
		return orthtypes.PrimitiveSTR

	default:
		panic("Invalid type")
	}
}

// ===================================
//	BASEDON
// ===================================

// ModBasedOnType Modules a set of numbers based on the set's type
func ModBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return ModIntegers
	case strings.Contains(superType, "f"):
		return ModFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		panic("Can not use 'orthtypes.PrimitiveSTR' with '%' operation")
	default:
		panic("Invalid type")
	}
}

// SumBasedOnType sums a set of numbers based on the set's type
func SumBasedOnType(superType string) (func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand, error) {
	switch {
	case strings.Contains(superType, "i") ||
		strings.Contains(superType, "i8") ||
		strings.Contains(superType, "i16") ||
		strings.Contains(superType, "i32") ||
		strings.Contains(superType, "i64"):
		return SumIntegers, nil
	case strings.Contains(superType, "f32") ||
		strings.Contains(superType, "f64"):
		return SumFloats, nil
	case superType == orthtypes.PrimitiveSTR:
		return ConcatPrimitiveSTR, nil
	case strings.Contains(superType, orthtypes.ADDR):
		return nil, fmt.Errorf(orth_debug.StrangeUseOfVariable, orthtypes.ADDR, "PLUS")
	default:
		panic("Invalid type")
	}
}

// SubBasedOnType subs a set of numbers based on the set's type
func SubBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return SubIntegers
	case strings.Contains(superType, "f"):
		return SubFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		panic("Can not apply '-' operation to a string value")
	default:
		panic("Invalid type")
	}
}

// DivBasedOnType divides a set of numbers based on the set's type
func DivBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return DivideIntegers
	case strings.Contains(superType, "f"):
		return DivideFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		panic("Can not apply '/' operation to a string value")
	default:
		panic("Invalid type")
	}
}

// MultBasedOnType multiplies a set of numbers based on the set's type
func MultBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return MultplyIntegers
	case strings.Contains(superType, "f"):
		return MultplyFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		panic("Can not apply '*' operation to a string value")
	default:
		panic("Invalid type")
	}
}

// EqualBasedOnType compare a set of numbers based on the set's type
func EqualBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return EqualInts
	case strings.Contains(superType, "f"):
		return EqualFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		return EqualString
	default:
		panic("Invalid type")
	}
}

// LowerBasedOnType compare a set of numbers based on the set's type
func LowerBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return LowerThanInts
	case strings.Contains(superType, "f"):
		return LowerThanFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		panic("Can not apply '<' in a string literal")
	default:
		panic("Invalid type")
	}
}

// GreaterBasedOnType compare a set of numbers based on the set's type
func GreaterBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return GreaterThanInts
	case strings.Contains(superType, "f"):
		return GreaterThanFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		panic("Can not apply '>' in a string literal")
	default:
		panic("Invalid type")
	}
}

// NotEqualBasedOnType compare a set of numbers based on the set's type
func NotEqualBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return DiffInts
	case strings.Contains(superType, "f"):
		return DiffFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		return DiffString
	default:
		panic("Invalid type")
	}
}

// EqualInts compare a set of integers
func EqualInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		res = "1"
	} else {
		res = "0"
	}

	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: res,
	}
}

// DiffInts compare a set of integers
func DiffInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	o1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}

	o2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: fmt.Sprintf("%v", o1 != o2),
	}
}

// LowerThanInts compare a set of integers
func LowerThanInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		op = "1"
	} else {
		op = "0"
	}
	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: op,
	}
}

// GreaterThanInts compare a set of integers
func GreaterThanInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		op = "1"
	} else {
		op = "0"
	}
	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: op,
	}
}

// MultplyIntegers multiplies a set of integers
func MultplyIntegers(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var mult string
	switch superType {
	case "i64":
		mult = MultI64(n1, n2)
	case "i32":
		mult = MultI32(n1, n2)
	case "i16":
		mult = MultI16(n1, n2)
	case "i8":
		mult = MultI8(n1, n2)
	case "i":
		mult = MultI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: mult,
	}
}

// DivideIntegers divides a set of integers
func DivideIntegers(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var div string
	switch superType {
	case "i64":
		div = DivI64(n1, n2)
	case "i32":
		div = DivI32(n1, n2)
	case "i16":
		div = DivI16(n1, n2)
	case "i8":
		div = DivI8(n1, n2)
	case "i":
		div = DivI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: div,
	}
}

// SumIntegers sums up a set of integers
func SumIntegers(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var sum string
	switch superType {
	case "i64":
		sum = SumI64(n1, n2)
	case "i32":
		sum = SumI32(n1, n2)
	case "i16":
		sum = SumI16(n1, n2)
	case "i8":
		sum = SumI8(n1, n2)
	case "i":
		sum = SumI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: sum,
	}
}

// ModIntegers modules a set of integers
func ModIntegers(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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

	return orthtypes.Operand{
		VarType: superType,
		Operand: sum,
	}
}

// SubIntegers subtracts a set of integers
func SubIntegers(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var sub string
	switch superType {
	case "i64":
		sub = SubI64(n1, n2)
	case "i32":
		sub = SubI32(n1, n2)
	case "i16":
		sub = SubI16(n1, n2)
	case "i8":
		sub = SubI8(n1, n2)
	case "i":
		sub = SubI(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: sub,
	}
}

// IntSupersetOfSlice gets the super type of a slice of integers
func IntSupersetOfSlice(opreands ...orthtypes.Operand) string {
	for _, v := range opreands {
		if v.VarType == "i64" {
			return v.VarType
		}
	}

	for _, v := range opreands {
		if v.VarType == "i32" {
			return v.VarType
		}
	}

	for _, v := range opreands {
		if v.VarType == "i16" {
			return v.VarType
		}
	}

	for _, v := range opreands {
		if v.VarType == "i8" {
			return v.VarType
		}
	}

	return "i"
}

// ModFloats divides a set of floays
func ModFloats(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var sum string
	switch superType {
	case "f64":
		sum = ModF64(n1, n2)
	case "f32":
		sum = ModF32(n1, n2)
	default:
		panic("Not an float")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: sum,
	}
}

// EqualFloats compare a set of floats
func EqualFloats(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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

	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: res,
	}
}

// DiffFloats compare a set of floats
func DiffFloats(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	o1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	o2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: fmt.Sprintf("%v", o1 != o2),
	}
}

// LowerThanFloats compare a set of floats
func LowerThanFloats(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	o1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	o2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: fmt.Sprintf("%v", o1 < o2),
	}
}

// GreaterThanFloats compare a set of floats
func GreaterThanFloats(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: op,
	}
}

// MultplyFloats multiplies a set of floats
func MultplyFloats(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var mult string
	switch superType {
	case "f64":
		mult = MultF64(n1, n2)
	case "f32":
		mult = MultF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: mult,
	}
}

// DivideFloats divides a set of floats
func DivideFloats(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var div string
	switch superType {
	case "f64":
		div = DivF64(n1, n2)
	case "f32":
		div = DivF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: div,
	}
}

// SumFloats sums up a set of floats
func SumFloats(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var sum string
	switch superType {
	case "f64":
		sum = SumF64(n1, n2)
	case "f32":
		sum = SumF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: sum,
	}
}

// SubFloats subtract a set of floats
func SubFloats(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	var sub string
	switch superType {
	case "f64":
		sub = SubF64(n1, n2)
	case "f32":
		sub = SubF32(n1, n2)
	default:
		panic("Not an integer")
	}

	return orthtypes.Operand{
		VarType: superType,
		Operand: sub,
	}
}

// FloatSupersetOfSlice gets the super type of a slice of floats
func FloatSupersetOfSlice(opreands ...orthtypes.Operand) string {
	for _, v := range opreands {
		if v.VarType == "f64" {
			return v.VarType
		}
	}

	for _, v := range opreands {
		if v.VarType == "f32" {
			return v.VarType
		}
	}
	panic("Invalid type")
}

// DiffString check if two string are different
func DiffString(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: fmt.Sprintf("%v", n1.Operand != n2.Operand),
	}
}

// EqualString checks if two variables have the same Operand
func EqualString(_ string, o1, o2 orthtypes.Operand) orthtypes.Operand {

	var res string
	if o1 == o2 {
		res = "1"
	} else {
		res = "0"
	}

	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: res,
	}
}

// ConcatPrimitiveSTR concats two string into 1
func ConcatPrimitiveSTR(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	return orthtypes.Operand{
		VarType: superType,
		Operand: n2.Operand + n1.Operand,
	}
}

// ToString converts any Operand's operand to it's string literal
func ToString(n1 orthtypes.Operand) orthtypes.Operand {
	return orthtypes.Operand{
		VarType: "s",
		Operand: n1.Operand,
	}
}
