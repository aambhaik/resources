package util

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-cli/env"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling-lib/util"
	"io/ioutil"
	"os"
	"strings"
)

func GetGithubResource(gitHubPath string, resourceFile string) ([]byte, error) {
	gbProject := env.NewGbProjectEnv()

	gbProject.Init(os.Getenv("GOPATH"))

	resourceDir := gbProject.GetVendorSrcDir()
	resourcePath := resourceDir + "/" + gitHubPath + "/" + resourceFile

	gbProject.InstallDependency(gitHubPath, "")

	data, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return nil, err
	}

	err = gbProject.UninstallDependency(gitHubPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetTriggerMetadata(gitHubPath string) (*ftrigger.Metadata, error) {
	gbProject := env.NewGbProjectEnv()

	gbProject.Init(os.Getenv("GOPATH"))

	resourceDir := gbProject.GetVendorSrcDir()
	triggerPath := resourceDir + "/" + gitHubPath + "/" + util.Gateway_Trigger_Metadata_JSON_Name

	gbProject.InstallDependency(gitHubPath, "")
	data, err := ioutil.ReadFile(triggerPath)
	if err != nil {
		return nil, err
	}
	triggerMetadata := &ftrigger.Metadata{}
	json.Unmarshal(data, triggerMetadata)
	return triggerMetadata, nil
}

func IsValidTriggerSetting(metadata *ftrigger.Metadata, property string) bool {
	settings := metadata.Settings
	for key := range settings {
		if key == property {
			return true
		}
	}

	return false
}

func IsValidTriggerHandlerSetting(metadata *ftrigger.Metadata, property string) bool {
	settings := metadata.Handler.Settings

	for _, element := range settings {
		if element.Name == property {
			return true
		}
	}

	return false
}

func ValidateTriggerConfigExpr(expression *string) (bool, *string) {
	if expression == nil {
		return false, nil
	}

	exprValue := *expression
	if strings.HasPrefix(exprValue, util.Gateway_Trigger_Config_Prefix) && strings.HasSuffix(exprValue, util.Gateway_Trigger_Config_Suffix) {
		//get name of the config
		str := exprValue[17 : len(exprValue)-1]
		return true, &str
	} else {
		return false, &exprValue
	}
}

func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}
