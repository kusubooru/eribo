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
	CmdJojo
	CmdLoth
	CmdDadJoke
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
	case CmdJojo:
		return "!jojo"
	case CmdLoth:
		return "!loth"
	case CmdDadJoke:
		return "!dadjoke"
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
	case "!jojo":
		return CmdJojo
	case "!loth":
		return CmdLoth
	case "!dadjoke":
		return CmdDadJoke
	}
}

func ParseCommand(s string) Command {
	cmd := strings.Fields(s)
	if len(cmd) == 0 {
		return CmdUnknown
	}
	return makeCommand(cmd[0])
}

func ParseCustomCommand(s string) (cmd string, args []string) {
	args = []string{}
	f := strings.Fields(s)
	if len(f) == 0 {
		return
	}
	cmd = f[0]
	if len(f) > 1 {
		args = f[1:]
	}
	return
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
