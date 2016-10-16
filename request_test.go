package allot

import "testing"

var result Request

func TestRequest(t *testing.T) {
	cmd := "command first second thrid"
	req := NewRequest(cmd)

	if req.Text() != cmd {
		t.Errorf("Text() must return \"%s\", but got \"%s\"", req.Text(), cmd)
	}

	if req.Command() != "command" {
		t.Errorf("Command() must return \"%s\", but got \"%s\"", req.Command(), "command")
	}
}

func BenchmarkRequest(b *testing.B) {
	var r Request

	for n := 0; n < b.N; n++ {
		r = NewRequest("command first second thrid")
	}

	result = r
}
