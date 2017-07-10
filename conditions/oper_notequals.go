package condition

import (
	"strings"
)

/**
specify the exact name of the operator such that the operator can be
used directly in an expression in mashling event-links. the operator
must be preceded by a space (' ') and succeeded by a space (' ') when
used in an expression.

e.g. ${trigger.content.country != USA}
*/
var infoNotEquals = &OperatorInfo{
	Name:        "!=",
	Description: `Support for not-equals operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&NotEquals{info: infoNotEquals})
}

type NotEquals struct {
	info *OperatorInfo
}

func (o *NotEquals) OperatorInfo() *OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *NotEquals) Eval(lhs string, rhs string) bool {
	if strings.Compare(lhs, rhs) == 0 {
		return false
	}
	return true
}
