package orth_debug

import (
	"fmt"
	orth_types "orth/cmd/pkg/types"
	"strings"
)

func PPrintOperation(op orth_types.Operation) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", orth_types.InstructionToStr(op.Instruction)))
	builder.WriteString(fmt.Sprintf("	operand: %s | symbolName%q\n", op.Operator.Operand, op.Operator.SymbolName))
	for k, v := range op.Links {
		builder.WriteString(fmt.Sprintf("	link_name: %q | link_type: %q | link_value: %q\n", k, v.Operator.SymbolName, v.Operator.Operand))
	}
	for k, v := range op.Addresses {
		builder.WriteString(fmt.Sprintf("\n** %s: %d\n", orth_types.InstructionToStr(k), v))
	}
	builder.WriteString("****************************************************\n")
	return builder.String()
}
