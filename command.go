package allot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// CommandInterface describes how to access a Command
type CommandInterface interface {
	Expression() (*regexp.Regexp, error)
	Has(name ParameterInterface) (bool, error)
	Match(req string) (MatchInterface, error)
	Matches(req string) (bool, error)
	Parameters() ([]Parameter, error)
	Position(param ParameterInterface) (int, error)
	Text() string
}

// Command is a Command definition
type Command struct {
	text   string
	escape bool
}

// Text returns the command text
func (c Command) Text() string {
	return c.text
}

// Expression returns the regular expression matching the command text
func (c Command) Expression() (*regexp.Regexp, error) {
	var expr string

	if c.escape {
		expr = regexp.QuoteMeta(c.Text())
	} else {
		expr = c.Text()
	}

	params, err := c.Parameters()
	if err != nil {
		return nil, err
	}

	for _, param := range params  {
		expr = strings.Replace(
			expr,
			fmt.Sprintf("<%s:%s>", param.Name(), param.Data()),
			fmt.Sprintf("(%s)", param.Expression().String()),
			-1,
		)

		expr = strings.Replace(
			expr,
			fmt.Sprintf("<%s>", param.Name()),
			"("+param.Expression().String()+")",
			-1,
		)
	}

	regex, err := regexp.Compile(fmt.Sprintf("^%s$", expr))
	if err != nil {
		return nil, err
	}

	return regex, nil
}

// Parameters returns the list of defined parameters
func (c Command) Parameters() ([]Parameter, error) {
	var list []Parameter
	re, err := regexp.Compile("<(.*?)>")
	if err != nil {
		return nil, err
	}

	result := re.FindAllStringSubmatch(c.Text(), -1)

	for _, p := range result {
		if len(p) != 2 {
			continue
		}

		pType := ""
		if !strings.Contains(p[1], ":") {
			pType = ":string"
		}

		list = append(list, Parse(p[1]+pType))
	}

	return list, nil
}

// Has checks if the parameter is found in the command
func (c Command) Has(param ParameterInterface) (bool, error) {
	pos, err := c.Position(param)
	if err != nil {
		return false, err
	}

	return pos != -1, nil
}

// Position returns the position of a parameter
func (c Command) Position(param ParameterInterface) (int, error) {
	params, err := c.Parameters()
	if err != nil {
		return -1, err
	}

	for index, item := range params {
		if item.Equals(param) {
			return index, nil
		}
	}

	return -1, nil
}

// Match returns the parameter matching the expression at the defined position
func (c Command) Match(req string) (MatchInterface, error) {
	result, err := c.Matches(req)
	if err != nil {
		return nil, err
	}

	if result {
		return Match{c, req}, nil
	}

	return nil, errors.New("Request does not match Command.")
}

// Matches checks if a command definition matches a request
func (c Command) Matches(req string) (bool, error) {
	expr, err := c.Expression()
	if err != nil {
		return false, err
	}

	return expr.MatchString(req), nil
}

// New returns a new command
func New(command string) Command {
	return Command{command, false}
}

// NewWithEscaping returns a new command that escapes regex characters
func NewWithEscaping(command string) Command {
	return Command{command, true}
}
