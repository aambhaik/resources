package condition

import (
	"errors"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/aambhaik/resources/util"
	"strings"
	"sync"
)

var fLogger = logger.GetLogger("event-link-operator")

var (
	OperatorRegistry = NewOperatorRegistry()
)

type Condition struct {
	Operator
	LHS string
	RHS string
}

type Operator interface {
	HasOperatorInfo
	Eval(lhs string, rhs string) bool
}

// OperRegistry is a registry for condition operators
type OperRegistry struct {
	operatorsMu   sync.Mutex
	operatorNames []string
	operators     map[string]Operator
}

// NewOperatorRegistry creates a new operator registry
func NewOperatorRegistry() *OperRegistry {
	return &OperRegistry{operators: make(map[string]Operator)}
}

// RegisterOperator registers an operator
func (r *OperRegistry) RegisterOperator(operator Operator) {

	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	if operator == nil {
		panic("OperatorRegistry: operator cannot be nil")
	}

	operatorName := operator.OperatorInfo().Name

	if _, exists := r.operators[operatorName]; exists {
		panic("OperatorRegistry: operator [" + operatorName + "] already registered")
	}
	if r.operatorNames == nil {
		r.operatorNames = append(r.operatorNames, operatorName)
	} else {
		inserted := false
		for index, o := range r.operatorNames {
			if strings.Contains(o, operatorName) {
				//the current operator is already part of another operator, so bump up the current operator
				r.operatorNames = Insert(r.operatorNames, index, operatorName)
				inserted = true
			}
		}
		if !inserted {
			r.operatorNames = append(r.operatorNames, operatorName)
		}
	}

	r.operators[operatorName] = operator
}

// Operator gets the specified operator
func (r *OperRegistry) Operator(operatorName string) (o Operator, exists bool) {

	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	operator, exists := r.operators[operatorName]
	return operator, exists
}

// Operators gets all the registered operators
func (r *OperRegistry) Operators() []Operator {

	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	var opers []Operator
	for _, v := range r.operators {
		opers = append(opers, v)
	}

	return opers
}

// Names gets all the registered operator names
func (r *OperRegistry) Names() []string {

	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	return r.operatorNames
}

// EvaluateExpression evaluates the specified expression
func EvaluateExpression(expression string, content string) bool {
	condition, err := getCondition(expression)
	if err != nil {
		fLogger.Debugf("Error getting the condition from expression [%v], [%v]", expression, err)
		return false
	}
	operator := condition.Operator
	lhsExpression := condition.LHS

	//evaluate the lhs against the content
	lhs, err := util.JsonPathEval(content, lhsExpression)
	if err != nil {
		return false
	}

	return operator.Eval(*lhs, condition.RHS)
}

// Evaluate evaluates the specified operator
func EvaluateCondition(condition Condition, content string) (bool, error) {
	//check if the message content is JSON first. mashling only supports JSON payloads for condition/content evaluation
	if !IsJSON(content) {
		return false, errors.New(fmt.Sprintf("Content is not a valid JSON payload [%v]", content))
	}

	operator := condition.Operator
	lhsExpression := condition.LHS

	//evaluate the lhs against the content
	lhs, err := util.JsonPathEval(content, lhsExpression)
	if err != nil {
		fLogger.Debugf("Error evaluating lhs jsonpath [%v] on the content [%v], [%v]", lhsExpression, content, err)
		return false, nil
	}

	return operator.Eval(*lhs, condition.RHS), nil
}

func getCondition(conditionExpr string) (*Condition, error) {
	oper, name, err := GetOperatorInExpression(conditionExpr)
	if err != nil {
		fLogger.Debugf("Error getting the operator from expression [%v], [%v]", conditionExpr, err)
		return nil, err
	}
	//found the operation!
	index := strings.Index(conditionExpr, *name)
	// find the LHS
	// Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
	lhs := strings.TrimSpace(conditionExpr[:index]) + "+"
	//get the value for LHS
	fLogger.Debugf("condition: left hand side found to be [%v", lhs)

	//find the RHS
	rhs := strings.TrimSpace(conditionExpr[index+len(*name):])
	fLogger.Debugf("condition: right hand side found to be [%v]", rhs)

	//create the condition
	condition := Condition{*oper, lhs, rhs}
	return &condition, nil

}

func GetConditionOperation(conditionStr string) (*Condition, error) {
	/**
	Content based conditions rules

	The condition identifier is "${" at the start and "}" at the end.

	If LHS
		If the condition clause starts with "trigger.content" then it refers to the trigger's payload. It maps internally to the "$." JSONPath of the payload.
		The above examples of JSONPath can be expressed as "${trigger.content.phoneNumbers[:1].type" and "${trigger.content.address.city" respectively.
		If the condition clause does not start with "trigger.content": TBD
		If it starts with "env" then it is evaluated as an environment variable. So, "${env.PROD_ENV == true}" will be evaluated as a condition based on the environment variable.
	If Operator
		The condition must evaluate to a boolean output. Example operators are "==" and "!=".
	If RHS
		The condition RHS will be interpreted as follows
		If the value on the RHS starts and ends with a single-quote (''), then it is accessed as a string
		If the value starts and ends without the single quote, then it is treated as an integer or a boolean.
	*/
	if !strings.HasPrefix(conditionStr, util.Gateway_Link_Condition_LHS_Start_Expr) {
		return nil, errors.New("If does not match expected semantics, missing '${' at the start.")
	}
	if !strings.HasSuffix(conditionStr, util.Gateway_Link_Condition_LHS_End_Expr) {
		return nil, errors.New("condition 'If' does not match expected semantics, missing '}' at the end.")
	}

	condition := conditionStr[len(util.Gateway_Link_Condition_LHS_Start_Expr) : len(conditionStr)-len(util.Gateway_Link_Condition_LHS_End_Expr)]
	contentRoot := GetContentRoot()

	if !strings.HasPrefix(condition, contentRoot) {
		return nil, errors.New(fmt.Sprintf("condition 'If' JSONPath must start with %v", contentRoot))
	}

	condition = strings.Replace(condition, contentRoot, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1)

	condition = strings.TrimSpace(condition)

	condOperation, err := getCondition(condition)
	if err != nil {
		return nil, err
	}
	return condOperation, nil

}

// Insert inserts the value into the slice at the specified index,
// which must be in range.
// The slice must have room for the new element.
func Insert(slice []string, index int, value string) []string {
	// Grow the slice by one element.
	slice = slice[0 : len(slice)+1]
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	return slice
}
