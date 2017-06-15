package util

import (
	"github.com/NodePrime/jsonpath"
)

func JsonPathEval(jsonData string, expression string) (*string, error) {
	paths, err := jsonpath.ParsePaths(expression)
	if err != nil {
		return nil, err
	}

	eval, err := jsonpath.EvalPathsInBytes([]byte(jsonData), paths)
	if err != nil {
		return nil, err
	}

	for {
		if result, ok := eval.Next(); ok {
			//return after the first match
			value := result.Pretty(false)
			return &value, nil // true -> show keys in pretty string
		} else {
			break
		}
	}
	return nil, nil
}
