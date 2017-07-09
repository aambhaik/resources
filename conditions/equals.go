package condition

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strings"
)

var fLogger = logger.GetLogger("event-link-operator")

var infoEquals = &OperatorInfo{
	Name:        "==",
	Description: `Support for equals operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&Equals{info: infoEquals})
	fLogger.Debugf("Operator equals registered!")
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
