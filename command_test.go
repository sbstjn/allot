package allot

import (
	"testing"
)

var resultCommand bool

func BenchmarkMatches(b *testing.B) {
	var r bool
	var cmd Command

	for n := 0; n < b.N; n++ {
		cmd = NewCommand("command <lorem:integer> <ipsum:string>")

		cmd.Matches("command 12345 abcdef")
	}

	resultCommand = r
}

func TestGetString(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter string
		value     string
	}{
		{"command <param1>", "command lorem", "param1", "lorem"},
		{"deploy <project:string> to <environment:string>", "deploy example to stage", "project", "example"},
		{"deploy <project:string> to <environment:string>", "deploy example to stage", "environment", "stage"},
	}

	for _, set := range data {
		cmd := NewCommand(set.command)

		value, err := cmd.GetString(set.request, set.parameter)

		if err != nil {
			t.Errorf("Parsign command returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%s\", expected \"%s\"", value, set.value)
		}
	}
}

func TestGetInteger(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter string
		value     int
	}{
		{"command <param1:integer>", "command 1234", "param1", 1234},
		{"revert from <project:string> last <commits:integer> commits", "revert from example last 51 commits", "commits", 51},
	}

	for _, set := range data {
		cmd := NewCommand(set.command)

		value, err := cmd.GetInteger(set.request, set.parameter)

		if err != nil {
			t.Errorf("Parsign command returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%d\", expected \"%d\"", value, set.value)
		}
	}
}

func TestMatches(t *testing.T) {
	var data = []struct {
		command string
		request string
		matches bool
	}{
		{"command", "example", false},
		{"command", "command", true},
		{"command", "command example", false},
		{"command <lorem>", "command", false},
		{"command <lorem>", "command exmaple", true},
		{"command <lorem>", "command 1234567", true},
		{"command <lorem>", "command command", true},
		{"command <lorem>", "example command", false},
		{"command <lorem:integer>", "command exmaple", false},
		{"command <lorem:integer>", "command 1234567", true},
		{"command <lorem>", "command command command", false},
	}

	for _, set := range data {
		cmd := NewCommand(set.command)

		if cmd.Matches(set.request) != set.matches {
			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"", cmd.Matches(set.request), set.matches, cmd.Expression().String(), set.request)
		}
	}
}

func TestPosition(t *testing.T) {
	var data = []struct {
		command string
		param   Parameter
		postion int
	}{
		{"command <lorem>", NewParameter("lorem", Expression("string")), 0},
		{"command <lorem> <ipsum> <dolor> <sit> <amet>", NewParameter("dolor", Expression("string")), 2},
		{"command <lorem> <ipsum> <dolor:integer> <sit> <amet>", NewParameter("dolor", Expression("string")), -1},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", NewParameter("lorem", Expression("string")), -1},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", NewParameter("lorem", Expression("integer")), 0},
		{"command <lorem:integer> <ipsum> <lorem:string> <sit> <amet>", NewParameter("lorem", Expression("integer")), 0},
		{"command <lorem:integer> <ipsum> <lorem:string> <sit> <amet>", NewParameter("lorem", Expression("string")), 2},
	}

	var cmd Command
	for _, set := range data {
		cmd = NewCommand(set.command)

		if cmd.Position(set.param) != set.postion {
			t.Errorf("Position() should be \"%d\", but is \"%d\"", set.postion, cmd.Position(set.param))
		}
	}
}

func TestGetParameter(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter Parameter
		has       bool
		value     string
	}{
		{"command <lorem>", "command test", NewParameterWithType("lorem", "string"), true, "test"},
		{"command <lorem>", "command 1234", NewParameterWithType("lorem", "string"), true, "1234"},
		{"command <lorem> <ipsum>", "command foo bar", NewParameterWithType("ipsum", "string"), true, "bar"},
		{"command <lorem> <ipsum>", "command foo bar", NewParameterWithType("lorem", "string"), true, "foo"},
		{"command <lorem> <ipsum>", "command foo bar", NewParameterWithType("example", "string"), false, "foo"},
		{"command <lorem:integer> <ipsum>", "command 123 bar", NewParameterWithType("lorem", "string"), false, "foo"},
	}

	var cmd CommandInterface
	for _, set := range data {
		cmd = NewCommand(set.command)

		if cmd.HasParameter(set.parameter) != set.has {
			t.Errorf("HasParameter is \"%v\", expected \"%v\"", cmd.HasParameter(set.parameter), set.has)
		}

		value, err := cmd.GetParameter(set.request, set.parameter)

		if err == nil && set.has && value != set.value {
			t.Errorf("GetParameter is \"%v\", expected \"%v\"", value, set.value)
		}
	}
}

func TestHasParameter(t *testing.T) {
	var data = []struct {
		command   string
		parameter Parameter
		has       bool
	}{
		{"command <lorem>", NewParameterWithType("lorem", "string"), true},
		{"command <lorem>", NewParameterWithType("lorem", "integer"), false},
	}

	var cmd CommandInterface
	for _, set := range data {
		cmd = NewCommand(set.command)

		if cmd.HasParameter(set.parameter) != set.has {
			t.Errorf("HasParameter is \"%v\", expected \"%v\"", cmd.HasParameter(set.parameter), set.has)
		}
	}
}

func TestParameters(t *testing.T) {
	var data = []struct {
		command    string
		name       string
		parameters []Parameter
	}{
		{"command <lorem>", "command", []Parameter{NewParameterWithType("lorem", "string")}},
		{"cmd <lorem:string>", "cmd", []Parameter{NewParameterWithType("lorem", "string")}},
		{"command <lorem:integer>", "command", []Parameter{NewParameterWithType("lorem", "integer")}},
		{"example <lorem> <ipsum> <dolor>", "example", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem> <ipsum> <dolor:string>", "command", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem> <ipsum:string> <dolor>", "command", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor>", "command", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor:string>", "command", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor:integer>", "command", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "integer")}},
		{"command <lorem:integer> <ipsum:string> <dolor:integer>", "command", []Parameter{NewParameterWithType("lorem", "integer"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "integer")}},
	}

	var cmd Command
	for _, set := range data {
		cmd = NewCommand(set.command)

		if cmd.Name() != set.name {
			t.Errorf("cmd.Name() must be \"%s\", but is \"%s\"", set.name, cmd.Name())
		}

		if cmd.Text() != set.command {
			t.Errorf("cmd.Text() must be \"%s\", but is \"%s\"", set.command, cmd.Text())
		}

		for _, param := range set.parameters {
			if !cmd.HasParameter(param) {
				t.Errorf("\"%s\" missing parameter.Item \"%s\"", cmd.Text(), param)
			}
		}
	}
}
