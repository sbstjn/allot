package allot

import (
	"testing"
)

var resultCommand bool

func BenchmarkMatches(b *testing.B) {
	var r bool

	for n := 0; n < b.N; n++ {
		cmd := New("command <lorem:integer> <ipsum:string>")
		result, err := cmd.Matches("command 12345 abcdef")
		if err != nil {
			b.Error(err)
			return
		}

		r = result
	}

	resultCommand = r
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
		{"command <lorem>", "command example", true},
		{"command <lorem>", "command 1234567", true},
		{"command <lorem>", "command command", true},
		{"command <lorem>", "example command", false},
		{"command <lorem:integer>", "command example", false},
		{"command <lorem:integer>", "command 1234567", true},
		{"command <lorem>", "command command command", false},
	}

	for _, set := range data {
		cmd := New(set.command)

		matches, err := cmd.Matches(set.request)
		if err != nil {
			t.Error(err)
		}

		if matches != set.matches {
			expr, err := cmd.Expression()
			if err != nil {
				t.Error(err)
			}

			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"", matches, set.matches, expr.String(), set.request)
		}
	}
}

func TestNotCompiling(t *testing.T) {
	cmd := New("command with []() invalid regex syntax")
	_, err := cmd.Matches("what ever")

	if err == nil {
		t.Error("Compilation of regex should have had failed")
	}
}

func TestEscapeMatches(t *testing.T) {
	var data = []struct {
		command string
		request string
		matches bool
	}{
		{"[command]", "example", false},
		{"[command]", "[command]", true},
		{"command", "command example", false},
		{"[command] (<lorem>)", "command", false},
		{"[command] (<lorem>)", "[command] (example)", true},
		{"[command] (<lorem>)", "[command] (1234)", true},
		{"[command] (<lorem:integer>)", "[command] (1234)", true},
	}

	for _, set := range data {
		cmd := NewWithEscaping(set.command)

		matches, err := cmd.Matches(set.request)
		if err != nil {
			t.Error(err)
			return
		}

		if matches != set.matches {
			expr, err := cmd.Expression()
			if err != nil {
				t.Error(err)
			}

			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"", matches, set.matches, expr.String(), set.request)
		}
	}
}

func TestPosition(t *testing.T) {
	var data = []struct {
		command  string
		param    Parameter
		position int
	}{
		{"command <lorem>", NewParameterWithType("lorem", "string"), 0},
		{"command <lorem> <ipsum> <dolor> <sit> <amet>", NewParameterWithType("dolor", "string"), 2},
		{"command <lorem> <ipsum> <dolor:integer> <sit> <amet>", NewParameterWithType("dolor", "string"), -1},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", NewParameterWithType("lorem", "string"), -1},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", NewParameterWithType("lorem", "integer"), 0},
		{"command <lorem:integer> <ipsum> <lorem:string> <sit> <amet>", NewParameterWithType("lorem", "integer"), 0},
		{"command <lorem:integer> <ipsum> <lorem:string> <sit> <amet>", NewParameterWithType("lorem", "string"), 2},
	}

	var cmd Command
	for _, set := range data {
		cmd = New(set.command)

		pos, err := cmd.Position(set.param)
		if err != nil {
			t.Error(err)
		}

		if pos != set.position {
			t.Errorf("Position() should be \"%d\", but is \"%d\"", set.position, pos)
		}
	}
}

func TestHas(t *testing.T) {
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
		cmd = New(set.command)
		has, err := cmd.Has(set.parameter)
		if err != nil {
			t.Error(err)
		}

		if has != set.has {
			t.Errorf("HasParameter is \"%v\", expected \"%v\"", has, set.has)
		}
	}
}

func TestParameters(t *testing.T) {
	var data = []struct {
		command    string
		parameters []Parameter
	}{
		{"command <lorem>", []Parameter{NewParameterWithType("lorem", "string")}},
		{"cmd <lorem:string>", []Parameter{NewParameterWithType("lorem", "string")}},
		{"command <lorem:integer>", []Parameter{NewParameterWithType("lorem", "integer")}},
		{"example <lorem> <ipsum> <dolor>", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem> <ipsum> <dolor:string>", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem> <ipsum:string> <dolor>", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor>", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor:string>", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor:integer>", []Parameter{NewParameterWithType("lorem", "string"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "integer")}},
		{"command <lorem:integer> <ipsum:string> <dolor:integer>", []Parameter{NewParameterWithType("lorem", "integer"), NewParameterWithType("ipsum", "string"), NewParameterWithType("dolor", "integer")}},
	}

	var cmd Command
	for _, set := range data {
		cmd = New(set.command)

		if cmd.Text() != set.command {
			t.Errorf("cmd.Text() must be \"%s\", but is \"%s\"", set.command, cmd.Text())
		}

		for _, param := range set.parameters {
			has, err := cmd.Has(param)
			if err != nil {
				t.Error(err)
			}

			if !has {
				t.Errorf("\"%s\" missing parameter.Item \"%+v\"", cmd.Text(), param)
			}
		}
	}
}
