package operators

import (
	"github.com/aambhaik/resources/conditions"
	"strings"
)

var infoEquals = &condition.OperatorInfo{
	Name:        "==",
	Description: `Support for equals operation to be used in the conditions`,
}

func init() {
	condition.OperatorRegistry.RegisterOperator(&Equals{info: infoEquals})
}

type Equals struct {
	info *condition.OperatorInfo
}

func (o *Equals) OperatorInfo() *condition.OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *Equals) Eval(lhs string, rhs string) bool {
	if strings.Compare(lhs, rhs) == 0 {
		return true
	}
	return false
}
