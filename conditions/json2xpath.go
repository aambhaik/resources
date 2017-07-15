package condition

import (
	"regexp"
	"strconv"
	"strings"
)

var jp []string
var jp2xp map[string]string
var xpDependence map[string]string
var xpArrayIndexIncrement = 1

func init() {
	/**
	Please refer to http://goessner.net/articles/JsonPath/ for more details and limitations on the mappings
	*/
	jp = []string{"(@.length-1)", "$.", "?(@.", ")]", "-1:", ".", "@", "*", "?(", "..", "[:"}

	jp2xp = make(map[string]string)
	jp2xp["(@.length-1)"] = "last()"
	jp2xp["$."] = "/"
	jp2xp["?(@."] = ""
	jp2xp[")]"] = "]"
	jp2xp["-1:"] = "last()"
	jp2xp["."] = "/"
	jp2xp["@"] = "."
	jp2xp["*"] = "*"
	jp2xp["?("] = "["
	jp2xp[".."] = "//"
	jp2xp["[:"] = "[position()<?]"

	xpDependence = make(map[string]string)
	xpDependence["(@.length-1)"] = ")]"
	xpDependence["(@.length-1)"] = ")]"
}

func ConvertJsonPathToXPath(jsonPathExpression string) (xPathExpression *string, err error) {
	originalExpression := jsonPathExpression
	//increment all array index numbers by the index-increment. find all numbers between '[' and ']' characters
	bracketPattern := `\[(.\[?)\]`
	bracketStartPattern := regexp.MustCompile(bracketPattern)
	bracketStartMatches := bracketStartPattern.FindAllStringSubmatch(jsonPathExpression, -1)
	if bracketStartMatches != nil && len(bracketStartMatches) > 0 {
		for _, match := range bracketStartMatches {
			numberStr := match[1]
			i, err := strconv.Atoi(numberStr)
			if err == nil {
				//increment the number by the factor
				replaceString := match[0]
				newNumberString := strconv.Itoa(i + xpArrayIndexIncrement)
				replaceString = strings.Replace(replaceString, numberStr, newNumberString, 1)
				jsonPathExpression = strings.Replace(jsonPathExpression, match[0], replaceString, 1)
			}
		}
	} else {
		bracketPattern := `\[:(.\[?)\]`
		bracketStartPattern := regexp.MustCompile(bracketPattern)
		bracketStartMatches = bracketStartPattern.FindAllStringSubmatch(jsonPathExpression, -1)
		if bracketStartMatches != nil && len(bracketStartMatches) > 0 {
			for _, match := range bracketStartMatches {
				//this is a special case with position less than operator
				numberStr := match[1]
				i, err := strconv.Atoi(numberStr)
				if err == nil {
					//increment the number by the factor
					newNumberString := strconv.Itoa(i + xpArrayIndexIncrement)
					positionStr := jp2xp["[:"]
					positionStr = strings.Replace(positionStr, "?", newNumberString, 1)
					jsonPathExpression = strings.Replace(jsonPathExpression, match[0], positionStr, 1)
				}
			}
		}
	}

	var skipJPKeys []string
	//replace all occurrences of jp2xp keys
	for _, k := range jp {
		if strings.Contains(originalExpression, k) && !stringInSlice(k, skipJPKeys) {
			//string contains the jsoanpath key, check if the key has any dependent key that we need to skip in further iteration
			if dependent, ok := xpDependence[k]; ok {
				skipJPKeys = append(skipJPKeys, dependent)
			}
			jsonPathExpression = strings.Replace(jsonPathExpression, k, jp2xp[k], -1)
		}
	}
	//return the xpath expression
	return &jsonPathExpression, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
