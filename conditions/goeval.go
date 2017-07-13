package condition

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/TIBCOSoftware/mashling-lib/conditions"
	"github.com/aambhaik/resources/util"
	"regexp"
	"strconv"
	"strings"
)

func GoEvaluateCondition(expression string, payload string) (bool, error) {
	originalExpression := expression
	/**
	The condition syntax uses govaluate expressions. Please refer to github.com/Knetic/govaluate for more information

	The value on LHS is identified with "${" at the start and "}" at the end and contains a valid JSONPath

	LHS
		If the LHS clause starts with "trigger.content" then it refers to the trigger's payload. It maps internally to the "$." JSONPath of the payload.

	The expression may contain more than one condition connected by logical operators '&&' or '||' so we need to
	first check the expression for multiple occurrences of expected identifiers.
	*/
	expressionKeyValues := make(map[string]string)

	startPattern := regexp.MustCompile(`\$\{`)
	startMatches := startPattern.FindAllStringIndex(originalExpression, -1)
	for i, match := range startMatches {
		startIndex := match[0]
		if startIndex < 0 {
			return false, fmt.Errorf("Condition LHS expresssion must start with [%v], invalid expression: [%v]", util.Gateway_Link_Condition_LHS_Start_Expr, originalExpression)
		}
		tempExpression := originalExpression[startIndex:]
		endIndex := strings.Index(tempExpression, util.Gateway_Link_Condition_LHS_End_Expr)
		if endIndex < 0 {
			return false, fmt.Errorf("Condition LHS expresssion must end with [%v], invalid expression: [%v]", util.Gateway_Link_Condition_LHS_End_Expr, originalExpression)
		}
		lhsCondition := tempExpression[len(util.Gateway_Link_Condition_LHS_Start_Expr):endIndex]
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
				return false, fmt.Errorf("Failed to evaluate JSONPath expression [%v] on payload [%v]", lhsCondition, payload)
			}
			output = value
		} else if IsXML(payload) {
			lhsCondition = strings.TrimSpace(lhsCondition)
			// jsonpath expression will begin with $. this translates to the root '/' in xpath.
			// subsequent json nodes are identified by '.' which again translates to '/' in xpath
			xpathExpression := strings.Replace(lhsCondition, "$.", "/", -1)
			xpathExpression = strings.Replace(xpathExpression, ".", "/", -1)

			value, err := util.XPathEval(payload, xpathExpression)
			if err != nil {
				panic(fmt.Errorf("Failed to evaluate XPath expression [%v] on payload [%v]", lhsCondition, payload))
			}
			output = value
		} else {
			//unsupported data format. the supported formats are JSON and XML
			panic(fmt.Errorf("Unknown data format on payload [%v] \nSupported formats are JSON and XML", payload))
		}

		//substitute the result of the JSONPath evaluation in to the original LHS expression
		oldString := tempExpression[:endIndex+len(util.Gateway_Link_Condition_LHS_End_Expr)]
		lhs := "lhs" + strconv.Itoa(i)
		expression = strings.Replace(expression, oldString, lhs, -1)

		expressionKeyValues[lhs] = *output
	}

	return goeval(expression, expressionKeyValues)
}

func goeval(exp string, kvMap map[string]string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		fmt.Sprintf("invalid expression [%v]", exp)
	}
	parameters := make(map[string]interface{})
	for k, v := range kvMap {
		number, err := strconv.ParseFloat(v, 64)
		if err == nil {
			//looks like the value is a number
			parameters[k] = number
		} else {
			parameters[k] = v
		}
	}

	r, err := expression.Evaluate(parameters)
	if err != nil {
		return false, fmt.Errorf("invalid expression [%v], err [%v]", exp, err)
	}

	return r.(bool), nil
}
