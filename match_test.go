package allot

import "testing"

func TestMatchAndIteger(t *testing.T) {
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
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request does not match Command")
		}

		value, err := match.Integer(set.parameter)

		if err != nil {
			t.Errorf("Parsign parameter returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%d\", expected \"%d\"", value, set.value)
		}
	}
}

func TestMatchAndString(t *testing.T) {
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
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request does not match Command")
		}

		value, err := match.String(set.parameter)

		if err != nil {
			t.Errorf("Parsign parameter returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%s\", expected \"%s\"", value, set.value)
		}
	}
}
