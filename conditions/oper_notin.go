package condition

import (
	"strings"
)

var infoNotIn = &OperatorInfo{
	Name:        "not in",
	Description: `Support for 'not in' operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&NotIn{info: infoNotIn})
}

type NotIn struct {
	info *OperatorInfo
}

func (o *NotIn) OperatorInfo() *OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *NotIn) Eval(lhs string, rhs string) bool {
	//RHS will be starting with '(' and ending with ')' and the values will be separated by a comma ','
	rhs = strings.TrimPrefix(rhs, "(")
	rhs = strings.TrimSuffix(rhs, ")")
	values := strings.Split(rhs, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == lhs {
			return false
		}
	}
	return true
}
