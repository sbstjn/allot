package allot

import (
	"regexp"
	"testing"
)

func TestExpression(t *testing.T) {
	var data = []struct {
		data       string
		expression string
	}{
		{"string", "[^\\s]+"},
		{"integer", "[0-9]+"},
		{"unknown", ""},
	}

	for _, set := range data {
		exp, err := regexp.Compile(set.expression)

		if err != nil {
			t.Errorf("TextExpression regexp does not compile: %s", set.expression)
		}

		test := Expression(set.data)

		if test == nil && set.expression != "" {
			t.Errorf("Expression for data \"%s\" is not valid", set.data)
		}

		if test != nil && test.String() != exp.String() {
			t.Errorf("Expression() not matching test data! got \"%s\", expected \"%s\"", test.String(), exp.String())
		}
	}
}

func TestParameterExpression(t *testing.T) {
	var data = []struct {
		name       string
		data       string
		expression string
	}{
		{"lorem", "string", "[^\\s]+"},
		{"ipsum", "integer", "[0-9]+"},
	}

	for _, set := range data {
		param := NewParameterWithType(set.name, set.data)
		exp := regexp.MustCompile(set.expression)

		pExp := param.Expression()
		if pExp.String() != exp.String() {
			t.Errorf("Expression() not matching test data! got \"%s\", expected \"%s\"", pExp.String(), exp.String())
		}

	}
}

func TestParse(t *testing.T) {
	var data = []struct {
		text string
		name string
		data string
	}{
		{"<lorem>", "lorem", "string"},
		{"<ipsum:integer>", "ipsum", "integer"},
	}

	var param Parameter
	for _, set := range data {
		param = Parse(set.text)

		if param.Name() != set.name {
			t.Errorf("param.Name() should be \"%s\", but is \"%s\"", set.name, param.Name())
		}
	}
}

func TestParameter(t *testing.T) {
	var data = []struct {
		name string
		data string
	}{
		{"lorem", "string"},
		{"ipsum", "integer"},
	}

	var param Parameter
	for _, set := range data {
		param = NewParameterWithType(set.name, set.data)

		if param.Name() != set.name {
			t.Errorf("param.Name() should be \"%s\", but is \"%s\"", set.name, param.Name())
		}
	}
}
