package allot

// MatchInterface is
type MatchInterface interface {
	String(name string) (string, error)
	Integer(name string) (int, error)
}

// Match is
type Match struct {
	Command CommandInterface
	Request string
}

// String returns strin parameter
func (m Match) String(name string) (string, error) {
	return m.Command.GetString(m.Request, name)
}

// Integer returns integer parameter
func (m Match) Integer(name string) (int, error) {
	return m.Command.GetInteger(m.Request, name)
}
