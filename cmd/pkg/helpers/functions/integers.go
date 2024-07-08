package functions

import (
	orth_types "orth/cmd/pkg/types"
	"strconv"
)

// =======================================
// SUM
// =======================================
func SumI64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int64(r1) + int64(r2)))
}

func SumI32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int32(r1) + int32(r2)))
}

func SumI16(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int16(r1) + int16(r2)))
}

func SumI8(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int8(r1) + int8(r2)))
}

func SumI(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int(r1) + int(r2)))
}

// =======================================
// MOD
// =======================================
func ModI64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int64(r1) % int64(r2)))
}

func ModI32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int32(r1) % int32(r2)))
}

func ModI16(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int16(r1) % int16(r2)))
}

func ModI8(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int8(r1) % int8(r2)))
}

func ModI(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int(r1) % int(r2)))
}

// =======================================
// SUBTRACTION
// =======================================
func SubI64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int64(r1) - int64(r2)))
}

func SubI32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int32(r1) - int32(r2)))
}

func SubI16(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int16(r1) - int16(r2)))
}

func SubI8(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int8(r1) - int8(r2)))
}

func SubI(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int(r1) - int(r2)))
}

// =======================================
// DIVISION
// =======================================
func DivI64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(int64(r1) / int64(r2)))
}

func DivI32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int32(r1) / int32(r2)))
}

func DivI16(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int16(r1) / int16(r2)))
}

func DivI8(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(int8(r1) / int8(r2)))
}

func DivI(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int(r1) / int(r2)))
}

// =======================================
// MULTIPLICATION
// =======================================
func MultI64(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(int64(r1) * int64(r2)))
}

func MultI32(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(int32(r1) * int32(r2)))
}

func MultI16(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(int16(r1) * int16(r2)))
}

func MultI8(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(int8(r1) * int8(r2)))
}

func MultI(n1, n2 orth_types.Operand) string {
	r1, err := strconv.Atoi(n1.Operand)
	if err != nil {
		panic(err)
	}
	r2, err := strconv.Atoi(n2.Operand)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(int(int(r1) * int(r2)))
}
