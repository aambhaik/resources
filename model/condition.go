package model

import (
	"fmt"
	"github.com/aambhaik/resources/util"
	"github.com/pkg/errors"
	"strings"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var flogoLogger = logger.GetLogger("trigger-tibco-kafkasubv2")

type ConditionalOperation interface {
	exec() bool
}

type If struct {
	Lhs string
	Rhs string
}

type Equals struct {
	If
}

type NotEquals struct {
	If
}

func (oper Equals) exec() bool {
	return oper.Lhs == oper.Rhs
}

func (oper NotEquals) exec() bool {
	return oper.Lhs != oper.Rhs
}

func GetConditionOperation(conditionStr string, content string) (*ConditionalOperation, error) {
	/**
	Content based conditions rules

	The condition identifier is "${" at the start and "}" at the end.

	If LHS
		If the condition clause starts with "trigger.content" then it refers to the trigger's payload. It maps internally to the "$." JSONPath of the payload.
		The above examples of JSONPath can be expressed as "${trigger.content.phoneNumbers[:1].type" and "${trigger.content.address.city" respectively.
		If the condition clause does not start with "trigger.content": TBD
		If it starts with "env" then it is evaluated as an environment variable. So, "${env.PROD_ENV == true}" will be evaluated as a condition based on the environment variable.
	If Operator
		The condition must evaluate to a boolean output. To start with, we can support only "==" and "!=" operators.
	If RHS
		The condition RHS will be interpreted as follows
		If the value on the RHS starts and ends with a single-quote (''), then it is accessed as a string
		If the value starts and ends without the single quote, then it is treated as an integer or a boolean.
	*/

	if !strings.HasPrefix(conditionStr, util.Gateway_Link_Condition_LHS_Start_Expr) {
		return nil, errors.New("If does not match expected semantics, missing '${' at the start.")
	}
	if !strings.HasSuffix(conditionStr, util.Gateway_Link_Condition_LHS_End_Expr) {
		return nil, errors.New("If does not match expected semantics, missing '}' at the end.")
	}

	condition := conditionStr[len(util.Gateway_Link_Condition_LHS_Start_Expr) : len(conditionStr)-len(util.Gateway_Link_Condition_LHS_End_Expr)]
	if !strings.HasPrefix(condition, util.Gateway_Link_Condition_LHS_JSON_Content) {
		return nil, errors.New(fmt.Sprintf("If JSONPath must start with %v", util.Gateway_Link_Condition_LHS_JSON_Content))
	}

	condition = strings.Replace(condition, util.Gateway_Link_Condition_LHS_JSON_Content, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1)

	var operation ConditionalOperation
	flogoLogger.Infof("condition is [%v]", condition)

	if index := strings.Index(condition, util.Gateway_Link_Condition_Operator_Equals); index > -1 {
		//operation is Equals
		//find the LHS
		lhs := condition[:index]
		//get the value for LHS
		flogoLogger.Infof("left hand side found to be [%v], content is [%v]", lhs, content)
		output, err := util.JsonPathEval(content, lhs)
		if err != nil {
			return nil, err
		}
		flogoLogger.Infof("json path eval output is [%v]", output)

		outputValue := *output

		//find the RHS
		rhs := condition[index+len(util.Gateway_Link_Condition_Operator_Equals):]
		flogoLogger.Infof("right hand side found to be [%v]", rhs)

		//create the equals struct instance
		operation = Equals{If{Lhs: outputValue, Rhs: rhs}}

	} else if index := strings.Index(condition, util.Gateway_Link_Condition_Operator_NotEquals); index > -1 {
		//operation is Not Equals

		//find the LHS
		lhs := condition[:index]
		//get the value for LHS
		output, err := util.JsonPathEval(content, lhs)
		if err != nil {
			return nil, err
		}
		outputValue := *output

		//find the RHS
		rhs := condition[index+len(util.Gateway_Link_Condition_Operator_NotEquals):]

		//create the equals struct instance
		operation = NotEquals{If{Lhs: outputValue, Rhs: rhs}}
	} else {
		//unknown operator?
		operators := []interface{}{util.Gateway_Link_Condition_Operator_Equals, util.Gateway_Link_Condition_Operator_NotEquals}
		return nil, errors.New(fmt.Sprintf("Unsupported operator found in the condition [%v], supported operators are [%v]", condition, operators))
	}

	return &operation, nil

}

func Evaluate(operation *ConditionalOperation) bool {
	oper := *operation
	return oper.exec()
}