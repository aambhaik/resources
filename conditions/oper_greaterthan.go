package condition

import (
	"strconv"
)

/**
specify the exact name of the operator such that the operator can be
used directly in an expression in mashling event-links. the operator
must be preceded by a space (' ') and succeeded by a space (' ') when
used in an expression.

e.g. ${trigger.content.amount > 20}
*/
var infoGreater = &OperatorInfo{
	Name:        ">",
	Description: `Support for 'greater than' numerical operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&GreaterThan{info: infoGreater})
}

type GreaterThan struct {
	info *OperatorInfo
}

func (o *GreaterThan) OperatorInfo() *OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *GreaterThan) Eval(lhs string, rhs string) bool {
	//convert the LHS and RHS to numbers
	lhsNum, err := strconv.Atoi(lhs)
	if err != nil {
		fLogger.Errorf("can not convert LHS value [%v] to an integer", lhs)
	}
	rhsNum, err := strconv.Atoi(rhs)
	if err != nil {
		fLogger.Errorf("can not convert RHS value [%v] to an integer", rhs)
	}

	if lhsNum > rhsNum {
		return true
	}
	return false
}
