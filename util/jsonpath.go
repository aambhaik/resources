package util

import (
	"github.com/NodePrime/jsonpath"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/pkg/errors"
	"fmt"
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
		flogoLogger.Infof("unable to eval expression, error is [%v]", err)
		return nil, err
	}

	for {
		if result, ok := eval.Next(); ok {
			//return after the first match
			value := string(result.Value)
			flogoLogger.Infof("expression parsed is [%v], value is [%v]", expression, value)
			return &value, nil // true -> show keys in pretty string
		} else {
			flogoLogger.Infof("evaluation failed..")
			break
		}
	}
	return nil, errors.New(fmt.Sprintf("Error with expression[%v] on content [%v]", expression, jsonData))
}
