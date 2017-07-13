package util

import (
	"bytes"
	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

func XPathEval(xmlPayload string, xpathExpression string) (*string, error) {
	var xpExec = goxpath.MustParse(xpathExpression)
	xTree, err := xmltree.ParseXML(bytes.NewBufferString(xmlPayload))
	if err != nil {
		return nil, err
	}

	res, err := xpExec.Exec(xTree)
	if err != nil {
		return nil, err
	}
	value := res.String()
	return &value, nil
}
