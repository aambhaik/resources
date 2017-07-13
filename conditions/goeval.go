package condition

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/TIBCOSoftware/mashling-lib/conditions"
	"github.com/aambhaik/resources/util"
	"strconv"
	"strings"
)

//func main() {
//	//exp := 	"!(${trigger.content.country} IN ('USA','IND','CHN','JPN'))"
//	//payload := `{"country":"MEX"}`
//
//	exp := "!(${trigger.content.address.country} == 'IND')"
//	payload := `{
//		"address": {
//			"country": "IND"
//		}
//	}`
//
//	//exp := "(${trigger.content.address.country} == 'USA')"
//	////payload := `{"amount":"3.14"}`
//	//payload := `<address><country>USA</country></address>`
//
//	GoEvaluateCondition(exp, payload)
//}

func GoEvaluateCondition(expression string, payload string) (bool, error) {
	originalExpression := expression
	/**
	The condition syntax uses govaluate expressions. Please refer to github.com/Knetic/govaluate for more information

	The value on LHS is identified with "${" at the start and "}" at the end and contains a valid JSONPath

	LHS
		If the LHS clause starts with "trigger.content" then it refers to the trigger's payload. It maps internally to the "$." JSONPath of the payload.
	*/
	start := strings.Index(expression, util.Gateway_Link_Condition_LHS_Start_Expr)
	if start < 0 {
		return false, fmt.Errorf("Condition LHS expresssion must start with [%v], invalid expression: [%v]", util.Gateway_Link_Condition_LHS_Start_Expr, originalExpression)
	}
	end := strings.Index(expression, util.Gateway_Link_Condition_LHS_End_Expr)
	if end < 0 {
		return false, fmt.Errorf("Condition LHS expresssion must end with [%v], invalid expression: [%v]", util.Gateway_Link_Condition_LHS_End_Expr, originalExpression)
	}

	lhsCondition := expression[start+len(util.Gateway_Link_Condition_LHS_Start_Expr) : end]
	contentRoot := condition.GetContentRoot()

	if !strings.HasPrefix(lhsCondition, contentRoot) {
		return false, fmt.Errorf("condition 'If' JSONPath must start with %v", contentRoot)
	}

	lhsCondition = strings.Replace(lhsCondition, contentRoot, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1)

	var output *string
	if condition.IsJSON(payload) {
		lhsCondition = strings.TrimSpace(lhsCondition) + "+"
		value, err := util.JsonPathEval(payload, lhsCondition)
		if err != nil {
			return false, fmt.Errorf("Failed to goeval JSONPath expression [%v] on payload [%v]", lhsCondition, payload)
		}
		output = value
	} else if IsXML(payload) {
		lhsCondition = strings.TrimSpace(lhsCondition)
		// jsonpath expression will begin with $. this translates to the root '/' in xpath.
		// subsequent json nodes are identified by '.' which again translates to '/' in xpath
		xpathExpression := strings.Replace(lhsCondition, "$.", "/", -1)
		xpathExpression = strings.Replace(xpathExpression, ".", "/", -1)

		value, err := XpathEval(payload, xpathExpression)
		if err != nil {
			return false, fmt.Errorf("Failed to goeval XPath expression [%v] on payload [%v]", lhsCondition, payload)
		}
		output = value
	} else {
		//unsupported data format. the supported formats are JSON and XML
		return false, fmt.Errorf("Unknown data format on payload [%v] \nSupported formats are JSON and XML", payload)
	}

	//substitute the result of the JSONPath evaluation in to the original LHS expression
	oldString := expression[start : end+len(util.Gateway_Link_Condition_LHS_End_Expr)]
	lhs := "lhs"
	expression = strings.Replace(expression, oldString, lhs, -1)

	return goeval(expression, lhs, *output)
}

func goeval(exp string, paramName string, paramValue string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		return false, fmt.Errorf("invalid expression [%v], err [%v]", exp, err)
	}
	parameters := make(map[string]interface{})
	number, err := strconv.ParseFloat(paramValue, 64)
	if err == nil {
		//looks like the value is a number
		parameters[paramName] = number
	} else {
		parameters[paramName] = paramValue
	}

	r, err := expression.Evaluate(parameters)
	if err != nil {
		return false, fmt.Errorf("invalid expression [%v], err [%v]", exp, err)
	}

	return r.(bool), nil
}
