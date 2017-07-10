package condition

import (
	"encoding/json"
	"fmt"
	"github.com/aambhaik/resources/util"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strings"
)

func GetOperatorInExpression(expression string) (*Operator, *string, error) {
	var oper *Operator
	var operatorName *string
	operators := OperatorRegistry.Operators()
	for _, o := range operators {
		names := o.OperatorInfo().Names
		for _, name := range names {
			// Find words in the expression that *start* with operator
			pattern := `\b` + " " + name + " "
			r, _ := regexp.Compile(pattern)

			if r.MatchString(expression) {
				op, exists := OperatorRegistry.Operator(name)
				if !exists {
					continue
				} else {
					oper = &op
					operatorName = &name
					break
				}
			}
		}

	}
	if oper == nil {
		return nil, nil, errors.Errorf("invalid operators found in expression [%v]", expression)
	}
	return oper, operatorName, nil
}

func ValidateOperatorInExpression(expression string) {
	originalExpression := expression
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
	if !strings.HasPrefix(expression, util.Gateway_Link_Condition_LHS_Start_Expr) {
		panic(fmt.Errorf("Condition expresssion must start with [%v], invalid expression: [%v]", util.Gateway_Link_Condition_LHS_Start_Expr, originalExpression))
	}
	if !strings.HasSuffix(expression, util.Gateway_Link_Condition_LHS_End_Expr) {
		panic(fmt.Errorf("Condition expresssion must end with [%v], invalid expression: [%v]", util.Gateway_Link_Condition_LHS_End_Expr, originalExpression))
	}

	expression = expression[len(util.Gateway_Link_Condition_LHS_Start_Expr) : len(expression)-len(util.Gateway_Link_Condition_LHS_End_Expr)]
	contentRoot := GetContentRoot()

	if !strings.HasPrefix(expression, contentRoot) {
		panic(fmt.Errorf("Condition expression must start with prefix [%v]", contentRoot))
	}

	expression = strings.Replace(expression, contentRoot, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1)

	expression = strings.TrimSpace(expression)

	operFound := false
	operators := OperatorRegistry.Operators()
	for _, o := range operators {
		names := o.OperatorInfo().Names
		for _, name := range names {
			// Find words in the expression that *start* with operator
			pattern := `\b` + " " + name + " "
			r, _ := regexp.Compile(pattern)

			if r.MatchString(expression) {
				_, exists := OperatorRegistry.Operator(name)
				if !exists {
					continue
				} else {
					if !operFound {
						operFound = true
					} else {
						//already one operator was found in the expression. here's another!
						//multiple operators are not allowed in a single expression
						panic(fmt.Errorf("Multiple operators not allowed in expression: [%v]", originalExpression))
					}
				}
			}
		}

	}

	if !operFound {
		operators := make([]string, 0, len(OperatorRegistry.operators))
		for k := range OperatorRegistry.operators {
			operators = append(operators, k)
		}
		panic(fmt.Errorf("No valid operator found in expression: [%v], supported operators are %v", originalExpression, "["+strings.Join(operators, ", ")+"]"))
	}
}

func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func GetContentRoot() string {
	contentRoot := os.Getenv(util.Gateway_JSON_Content_Root_Env_Key)
	if contentRoot == "" {
		//use the default value
		contentRoot = util.Gateway_Link_Condition_LHS_JSON_Content_Prefix_Default
	}
	return contentRoot
}
