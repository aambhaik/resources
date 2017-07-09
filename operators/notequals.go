package operators

import (
	"github.com/aambhaik/resources/conditions"
	"strings"
)

var infoNotEquals = &condition.OperatorInfo{
	Name:        "!=",
	Description: `Support for not-equals operation to be used in the conditions`,
}

func init() {
	condition.OperatorRegistry.RegisterOperator(&NotEquals{info: infoNotEquals})
}

type NotEquals struct {
	info *condition.OperatorInfo
}

func (o *NotEquals) OperatorInfo() *condition.OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *NotEquals) Eval(lhs string, rhs string) bool {
	if strings.Compare(lhs, rhs) == 0 {
		return false
	}
	return true
}
