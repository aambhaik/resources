package util

import (
	"bytes"
	"github.com/pkg/errors"
	"gopkg.in/xmlpath.v1"
)

//func XPathEval(xmlPayload string, xpathExpression string) (*string, error) {
//	var xpExec = goxpath.MustParse(xpathExpression)
//	xTree, err := xmltree.ParseXML(bytes.NewBufferString(xmlPayload))
//	if err != nil {
//		return nil, err
//	}
//
//	res, err := xpExec.Exec(xTree)
//	if err != nil {
//		return nil, err
//	}
//	value := res.String()
//	return &value, nil
//}

func XPathEval(xmlPayload string, xpathExpression string) (*string, error) {
	path := xmlpath.MustCompile(xpathExpression)
	root, err := xmlpath.Parse(bytes.NewBufferString(xmlPayload))
	if err != nil {
		return nil, err
	}

	if value, ok := path.String(root); ok {
		return &value, nil
	}

	return nil, errors.Errorf("Unable to evaluate the XPath %v on payload %v ", xpathExpression, xmlPayload)
}
