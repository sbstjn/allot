package allot

import (
	"regexp"
	"strings"
)

var regexpMapping = map[string]string{
	"string":  "[^\\s]+",
	"integer": "[0-9]+",
}

// Expression returns the regexp for a data type
func Expression(datatype string) *regexp.Regexp {
	if exp, ok := regexpMapping[datatype]; ok {
		return regexp.MustCompile(exp)
	}

	return nil
}

// ParameterInterface describes how to access a Parameter
type ParameterInterface interface {
	Equals(param ParameterInterface) bool
	Expression() *regexp.Regexp
	Name() string
	Datatype() string
}

// Parameter is the Parameter definition
type Parameter struct {
	name     string
	datatype string
	expr     *regexp.Regexp
}

// Expression returns the regexp behind the type
func (p Parameter) Expression() *regexp.Regexp {
	return p.expr
}

// Name returns the Parameter name
func (p Parameter) Name() string {
	return p.name
}

// Data returns the Parameter name
func (p Parameter) Datatype() string {
	return p.datatype
}

// Equals checks if two parameter are equal
func (p Parameter) Equals(param ParameterInterface) bool {
	return p.Name() == param.Name() && p.Expression().String() == param.Expression().String()
}

// NewParameterWithType returns
func NewParameterWithType(name string, datatype string) Parameter {
	return Parameter{name, datatype, Expression(datatype)}
}

// Parse parses parameter info
func Parse(text string) Parameter {
	var splits []string
	var name, datatype string

	name = strings.Replace(text, "<", "", -1)
	name = strings.Replace(name, ">", "", -1)
	datatype = "string"

	if strings.Contains(name, ":") {
		splits = strings.Split(name, ":")

		name = splits[0]
		datatype = splits[1]
	}

	return NewParameterWithType(name, datatype)
}
