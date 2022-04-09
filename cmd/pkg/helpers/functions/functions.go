package functions

import (
	"fmt"
	"strconv"
	"strings"
	orthtypes "t/cmd/pkg/types"
)

func TypesAreEqual(opreands ...orthtypes.Operand) bool {
	t := opreands[0].VarType
	equal := true

	for _, n := range opreands {
		equal = n.VarType == t
	}

	return equal
}

func GetSupersetType(opreands ...orthtypes.Operand) string {
	switch {
	case strings.Contains(opreands[0].VarType, "i"):
		return IntSupersetOfSlice(opreands...)

	case strings.Contains(opreands[0].VarType, "f"):
		return FloatSupersetOfSlice(opreands...)

	case strings.Contains(opreands[0].VarType, orthtypes.PrimitiveSTR):
		return orthtypes.PrimitiveSTR
	default:
		panic("Invalid type")
	}
}

// ===================================
//	BASEDON
// ===================================
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

func SumBasedOnType(superType string) func(string, orthtypes.Operand, orthtypes.Operand) orthtypes.Operand {
	switch {
	case strings.Contains(superType, "i"):
		return SumIntegers
	case strings.Contains(superType, "f"):
		return SumFloats
	case strings.Contains(superType, orthtypes.PrimitiveSTR):
		return ConcatPrimitiveSTR
	default:
		panic("Invalid type")
	}
}

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

func EqualInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		Operand: fmt.Sprintf("%v", o1 == o2),
	}
}

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

func LowerThanInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		Operand: fmt.Sprintf("%v", o2 < o1),
	}
}

func GreaterThanInts(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		Operand: fmt.Sprintf("%v", o2 > o1),
	}
}

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

func EqualFloats(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		Operand: fmt.Sprintf("%v", o1 == o2),
	}
}

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
		Operand: fmt.Sprintf("%v", o2 < o1),
	}
}

func GreaterThanFloats(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
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
		Operand: fmt.Sprintf("%v", o2 > o1),
	}
}

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

func DiffString(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: fmt.Sprintf("%v", n1.Operand != n2.Operand),
	}
}

func EqualString(_ string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	return orthtypes.Operand{
		VarType: orthtypes.PrimitiveBOOL,
		Operand: fmt.Sprintf("%v", n1.Operand == n2.Operand),
	}
}

func ConcatPrimitiveSTR(superType string, n1, n2 orthtypes.Operand) orthtypes.Operand {
	return orthtypes.Operand{
		VarType: superType,
		Operand: n2.Operand + n1.Operand,
	}
}
