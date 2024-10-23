package functions

import (
	"math"
	orth_types "orth/cmd/pkg/types"
	"strconv"
)

// =======================================
// SUM
// =======================================
func SumF64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1+r2, 'f', 2, 64)
}

func SumF32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 32)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 32)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1+r2, 'f', 2, 32)
}

func ModF64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(math.Mod(r1, r2), 'f', 2, 64)
}

func ModF32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 32)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 32)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(math.Mod(r1, r2), 'f', 2, 32)
}

// =======================================
// SUBTRACTION
// =======================================
func SubF64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1-r2, 'f', 2, 64)
}

func SubF32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 32)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 32)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1-r2, 'f', 2, 32)
}

// =======================================
// DIVISION
// =======================================
func DivF64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1/r2, 'f', 2, 64)
}

func DivF32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 32)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 32)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1/r2, 'f', 2, 32)
}

// =======================================
// MULTIPLICATION
// =======================================
func MultF64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 64)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 64)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1*r2, 'f', 2, 64)
}

func MultF32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.ParseFloat(n1.Operand, 32)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.ParseFloat(n2.Operand, 32)
	if err != nil {
		panic(err)
	}

	return strconv.FormatFloat(r1*r2, 'f', 2, 32)
}
