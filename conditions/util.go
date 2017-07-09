package condition

import (
	"encoding/json"
	"github.com/aambhaik/resources/util"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func GetOperatorInExpression(expression string) (*Operator, *string, error) {
	var operator *Operator
	var operatorName *string
	operNames := OperatorRegistry.Names()
	for _, name := range operNames {
		if strings.Contains(expression, name) {
			oper, exists := OperatorRegistry.Operator(name)
			if !exists {
				continue
			} else {
				operator = &oper
				operatorName = &name
				break
			}
		}
	}
	if operator == nil {
		return nil, nil, errors.Errorf("invalid operators found in expression [%v]", expression)
	}
	return operator, operatorName, nil
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
