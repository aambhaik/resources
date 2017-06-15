package util

import (
	"github.com/NodePrime/jsonpath"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var flogoLogger = logger.GetLogger("trigger-tibco-kafkasubv2")

func JsonPathEval(jsonData string, expression string) (*string, error) {
	paths, err := jsonpath.ParsePaths(expression)
	if err != nil {
		return nil, err
	}
	flogoLogger.Infof("expression parsed is [%v], json is [%v]", expression, jsonData)

	eval, err := jsonpath.EvalPathsInBytes([]byte(jsonData), paths)
	if err != nil {
		return nil, err
	}

	for {
		if result, ok := eval.Next(); ok {
			//return after the first match
			value := result.Pretty(false)
			flogoLogger.Infof("expression parsed is [%v], value is [%v]", expression, value)
			return &value, nil // true -> show keys in pretty string
		} else {
			break
		}
	}
	return nil, nil
}
