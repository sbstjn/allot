package allot

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// CommandInterface is the interface
type CommandInterface interface {
	Expression() *regexp.Regexp
	GetInteger(req string, param string) (int, error)
	GetParameter(req string, param ParameterInterface) (string, error)
	GetString(req string, param string) (string, error)
	HasParameter(name ParameterInterface) bool
	Match(req string) (MatchInterface, error)
	Matches(req string) bool
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

// Position returns the position of a parameter
func (c Command) Position(param ParameterInterface) int {
	for index, item := range c.Parameters() {
		if item.Equals(param) {
			return index
		}
	}

	return -1
}

// GetParameter gets value for parameter
func (c Command) GetParameter(req string, param ParameterInterface) (string, error) {
	pos := c.Position(param)

	if pos == -1 {
		return "", errors.New("Unknonw parameter for string.")
	}

	matches := c.Expression().FindAllStringSubmatch(req, -1)[0][1:]
	return matches[c.Position(param)], nil
}

// GetString returns a string parameter
func (c Command) GetString(req string, param string) (string, error) {
	return c.GetParameter(req, NewParameterWithType(param, "string"))
}

// GetInteger returns an integer parameter
func (c Command) GetInteger(req string, param string) (int, error) {
	str, err := c.GetParameter(req, NewParameterWithType(param, "integer"))

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

// Match returns matches command
func (c Command) Match(req string) (MatchInterface, error) {
	if c.Matches(req) {
		return Match{c, req}, nil
	}

	return nil, errors.New("Request does not match Command.")
}

// Matches checks if a comand definition matches a request
func (c Command) Matches(req string) bool {
	return c.Name() == strings.Split(req, " ")[0] && c.Expression().MatchString(req)
}

// NewCommand returns a new command
func NewCommand(command string) Command {
	return Command{command}
}
