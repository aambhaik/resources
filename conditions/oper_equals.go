package condition

import (
	"strings"
)

/**
specify the exact name of the operator such that the operator can be
used directly in an expression in mashling event-links. the operator
must be preceded by a space (' ') and succeeded by a space (' ') when
used in an expression.

e.g. ${trigger.content.country == USA}
*/
var infoEquals = &OperatorInfo{
	Name:        "==",
	Description: `Support for equals operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&Equals{info: infoEquals})
}

type Equals struct {
	info *OperatorInfo
}

func (o *Equals) OperatorInfo() *OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *Equals) Eval(lhs string, rhs string) bool {
	if strings.Compare(lhs, rhs) == 0 {
		return true
	}
	return false
}
