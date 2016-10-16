package allot

import "strings"

// RequestInterface is
type RequestInterface interface {
	Text() string
	Command() string
}

// Request is
type Request struct {
	text string
}

// Text returns
func (r Request) Text() string {
	return r.text
}

// Command is
func (r Request) Command() string {
	return strings.Split(r.Text(), " ")[0]
}

// NewRequest returns
func NewRequest(req string) Request {
	return Request{req}
}
