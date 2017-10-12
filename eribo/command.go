package eribo

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Command int

const (
	CmdUnknown Command = iota
	CmdTomato
	CmdTieup
	CmdFeedback
	CmdTktool
	CmdVonprove
)

func (c Command) String() string {
	switch c {
	default:
		return ""
	case CmdTomato:
		return "!tomato"
	case CmdTieup:
		return "!tieup"
	case CmdFeedback:
		return "!feedback"
	case CmdTktool:
		return "!tktool"
	case CmdVonprove:
		return "!Vonprove"
	}
}

func makeCommand(s string) Command {
	switch s {
	default:
		return CmdUnknown
	case "!tomato":
		return CmdTomato
	case "!tieup":
		return CmdTieup
	case "!feedback":
		return CmdFeedback
	case "!tktool":
		return CmdTktool
	case "!Vonprove":
		return CmdVonprove
	}
}

func (c Command) HasPrefix(s string) bool {
	return strings.HasPrefix(s, c.String())
}

func (c Command) Value() (driver.Value, error) { return c.String(), nil }
func (c *Command) Scan(value interface{}) error {
	if value == nil {
		*c = CmdUnknown
		return nil
	}
	switch v := value.(type) {
	case string:
		*c = makeCommand(v)
		return nil
	case []byte:
		*c = makeCommand(string(v))
		return nil
	}
	return fmt.Errorf("cannot scan Command value")
}
