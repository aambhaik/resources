package util

import (
	"github.com/NodePrime/jsonpath"
	"fmt"
)

func JsonPathEval(jsonData string, expression string) (*string, error) {
	paths, err := jsonpath.ParsePaths(expression)
	if err != nil {
		return nil, err
	}
	fmt.Sprintf("expression parsed is [%v], json is [%v]", expression, jsonData)

	eval, err := jsonpath.EvalPathsInBytes([]byte(jsonData), paths)
	if err != nil {
		return nil, err
	}

	for {
		if result, ok := eval.Next(); ok {
			//return after the first match
			value := result.Pretty(false)
			fmt.Sprintf("expression parsed is [%v], value is [%v]", expression, value)
			return &value, nil // true -> show keys in pretty string
		} else {
			break
		}
	}
	return nil, nil
}
