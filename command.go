package allot

import (
	"regexp"
	"strings"
)

// CommandInterface is the interface
type CommandInterface interface {
	Expression() *regexp.Regexp
	GetParameter(req RequestInterface, param ParameterInterface) string
	HasParameter(name ParameterInterface) bool
	Matches(req RequestInterface) bool
	Name() string
	Parameters() []Parameter
	Position(param ParameterInterface) int
	Text() string
}

// Command is the struct
type Command struct {
	text string
}

// Name returns the command name
func (c Command) Name() string {
	return strings.Split(c.Text(), " ")[0]
}

// Text returns the text
func (c Command) Text() string {
	return c.text
}

// Expression returns the regex for the command
func (c Command) Expression() *regexp.Regexp {
	var params []string
	expr := c.Name()

	for _, param := range c.Parameters() {
		params = append(params, "("+param.Expression().String()+")")
	}

	if len(params) > 0 {
		expr = expr + " " + strings.Join(params, " ")
	}

	return regexp.MustCompile("^" + expr + "$")
}

// Parameters returns the list of parameters
func (c Command) Parameters() []Parameter {
	var list []Parameter
	splits := strings.Split(c.Text(), " ")

	for index, item := range splits {
		if index > 0 {
			list = append(list, Parse(item))
		}
	}

	return list
}

// HasParameter checks if parameter by name is available
func (c Command) HasParameter(param ParameterInterface) bool {
	for _, item := range c.Parameters() {
		if item.Equals(param) {
			return true
		}
	}

	return false
}

// Position returns the position of a paramter
func (c Command) Position(param ParameterInterface) int {
	for index, item := range c.Parameters() {
		if item.Equals(param) {
			return index
		}
	}

	return -1
}

// GetParameter gets value for parameter
func (c Command) GetParameter(req RequestInterface, param ParameterInterface) string {
	matches := c.Expression().FindAllStringSubmatch(req.Text(), -1)[0][1:]
	return matches[c.Position(param)]
}

// Matches checks if a comand definition matches a request
func (c Command) Matches(req RequestInterface) bool {
	return c.Name() == req.Command() && c.Expression().MatchString(req.Text())
}

// NewCommand returns a new command
func NewCommand(command string) Command {
	return Command{command}
}
