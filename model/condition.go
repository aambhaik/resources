package model

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	cutil "github.com/aambhaik/resources/conditions"
	"github.com/aambhaik/resources/util"
	"github.com/pkg/errors"
	"strings"
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

type In struct {
	If
}

type NotIn struct {
	If
}

func (oper Equals) exec() bool {
	return oper.Lhs == oper.Rhs
}

func (oper NotEquals) exec() bool {
	return oper.Lhs != oper.Rhs
}

func (oper In) exec() bool {
	//RHS will be starting with '(' and ending with ')' and the values will be separated by a comma ','
	oper.Rhs = strings.TrimPrefix(oper.Rhs, "(")
	oper.Rhs = strings.TrimSuffix(oper.Rhs, ")")
	values := strings.Split(oper.Rhs, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == oper.Lhs {
			return true
		}
	}
	return false
}

func (oper NotIn) exec() bool {
	//RHS will be starting with '(' and ending with ')' and the values will be separated by a comma ','
	oper.Rhs = strings.TrimPrefix(oper.Rhs, "(")
	oper.Rhs = strings.TrimSuffix(oper.Rhs, ")")
	values := strings.Split(oper.Rhs, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == oper.Lhs {
			return false
		}
	}
	return true
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

	//check if the message content is JSON first. mashling only supports JSON payloads for condition/content evaluation
	if !cutil.IsJSON(content) {
		return nil, errors.New(fmt.Sprintf("Content is not a valid JSON payload [%v]", content))
	}

	if !strings.HasPrefix(conditionStr, util.Gateway_Link_Condition_LHS_Start_Expr) {
		return nil, errors.New("If does not match expected semantics, missing '${' at the start.")
	}
	if !strings.HasSuffix(conditionStr, util.Gateway_Link_Condition_LHS_End_Expr) {
		return nil, errors.New("condition 'If' does not match expected semantics, missing '}' at the end.")
	}

	condition := conditionStr[len(util.Gateway_Link_Condition_LHS_Start_Expr) : len(conditionStr)-len(util.Gateway_Link_Condition_LHS_End_Expr)]
	contentRoot := cutil.GetContentRoot()

	if !strings.HasPrefix(condition, contentRoot) {
		return nil, errors.New(fmt.Sprintf("condition 'If' JSONPath must start with %v", contentRoot))
	}

	condition = strings.Replace(condition, contentRoot, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1)

	condition = strings.TrimSpace(condition)

	operation, err := getOperation(condition, content)
	if err != nil {
		return nil, err
	}
	return operation, nil

}

func getOperation(condition string, content string) (*ConditionalOperation, error) {
	var operation ConditionalOperation
	if index := strings.Index(condition, util.Gateway_Link_Condition_Operator_Equals); index > -1 {
		//operation is Equals
		//find the LHS
		lhs := strings.TrimSpace(condition[:index]) + "+" // Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
		//get the value for LHS
		flogoLogger.Debugf("condition: left hand side found to be [%v], content is [%v]", lhs, content)
		output, err := util.JsonPathEval(content, lhs)
		if err != nil {
			return nil, err
		}
		outputValue := *output
		flogoLogger.Debugf("json path eval output is [%v]", outputValue)

		//find the RHS
		rhs := strings.TrimSpace(condition[index+len(util.Gateway_Link_Condition_Operator_Equals):])
		flogoLogger.Debugf("condition: right hand side found to be [%v]", rhs)

		//create the equals struct instance
		operation = Equals{If{Lhs: outputValue, Rhs: rhs}}
		return &operation, nil

	} else if index := strings.Index(condition, util.Gateway_Link_Condition_Operator_NotEquals); index > -1 {
		//operation is Not Equals

		//find the LHS
		lhs := strings.TrimSpace(condition[:index]) + "+" // Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
		//get the value for LHS
		output, err := util.JsonPathEval(content, lhs)
		if err != nil {
			return nil, err
		}
		outputValue := *output

		//find the RHS
		rhs := strings.TrimSpace(condition[index+len(util.Gateway_Link_Condition_Operator_NotEquals):])

		//create the equals struct instance
		operation = NotEquals{If{Lhs: outputValue, Rhs: rhs}}
		return &operation, nil

	} else if index := strings.Index(condition, util.Gateway_Link_Condition_Operator_NotIn); index > -1 { //NotIn handling must come before In handling as the string 'in' occurs in 'not in'
		//operator is Not In

		//find the LHS
		lhs := strings.TrimSpace(condition[:index]) + "+" // Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
		//get the value for LHS
		output, err := util.JsonPathEval(content, lhs)
		if err != nil {
			return nil, err
		}
		outputValue := *output

		//find the RHS
		rhs := strings.TrimSpace(condition[index+len(util.Gateway_Link_Condition_Operator_NotIn):])

		//create the equals struct instance
		operation = NotIn{If{Lhs: outputValue, Rhs: rhs}}
		return &operation, nil

	} else if index := strings.Index(condition, util.Gateway_Link_Condition_Operator_In); index > -1 {
		//operator is In

		//find the LHS
		lhs := strings.TrimSpace(condition[:index]) + "+" // Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
		//get the value for LHS
		output, err := util.JsonPathEval(content, lhs)
		if err != nil {
			return nil, err
		}
		outputValue := *output

		//find the RHS
		rhs := strings.TrimSpace(condition[index+len(util.Gateway_Link_Condition_Operator_In):])

		//create the equals struct instance
		operation = In{If{Lhs: outputValue, Rhs: rhs}}
		return &operation, nil

	} else {
		//unknown operator?
		operators := []interface{}{util.Gateway_Link_Condition_Operator_Equals, util.Gateway_Link_Condition_Operator_NotEquals}
		return nil, errors.New(fmt.Sprintf("Unsupported operator found in the condition [%v], supported operators are [%v]", condition, operators))
	}
}

func Evaluate(operation *ConditionalOperation) bool {
	oper := *operation
	return oper.exec()
}
