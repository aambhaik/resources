package condition

import (
	"strings"
)

var infoIn = &OperatorInfo{
	Name:        "in",
	Description: `Support for 'in' operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&In{info: infoIn})
}

type In struct {
	info *OperatorInfo
}

func (o *In) OperatorInfo() *OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *In) Eval(lhs string, rhs string) bool {
	//RHS will be starting with '(' and ending with ')' and the values will be separated by a comma ','
	rhs = strings.TrimPrefix(rhs, "(")
	rhs = strings.TrimSuffix(rhs, ")")
	values := strings.Split(rhs, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == lhs {
			return true
		}
	}
	return false
}
