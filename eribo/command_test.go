package eribo

import (
	"reflect"
	"testing"
)

var parseCommandTests = []struct {
	in   string
	cmd  Command
	args []string
}{
	{"", CmdUnknown, []string{}},
	{"!tomato", CmdTomato, []string{}},
	{"!tomato ", CmdTomato, []string{}},
	{"!tomato    ", CmdTomato, []string{}},
	{"!tomato			", CmdTomato, []string{}},
	{"!tomato			1 2		3", CmdTomato, []string{"1", "2", "3"}},
	{" !tomato", CmdTomato, []string{}},
	{"foo !tomato", CmdUnknown, []string{"!tomato"}},
}

func TestParseCommand(t *testing.T) {
	for _, tt := range parseCommandTests {
		cmd, args := ParseCommand(tt.in)
		if got, want := cmd, tt.cmd; got != want {
			t.Errorf("ParseCommand(%q) cmd = %q, want %q", tt.in, got, want)
		}
		if got, want := args, tt.args; !reflect.DeepEqual(got, want) {
			t.Errorf("ParseCommand(%q) args = %q, want %q", tt.in, got, want)
		}
	}
}

var parseCustomCommandTests = []struct {
	in   string
	cmd  string
	args []string
}{
	{"", "", []string{}},
	{"!tomato", "!tomato", []string{}},
	{"!tomato ", "!tomato", []string{}},
	{"!tomato    ", "!tomato", []string{}},
	{"!tomato			", "!tomato", []string{}},
	{"!tomato			1 2		3", "!tomato", []string{"1", "2", "3"}},
	{" !tomato", "!tomato", []string{}},
	{"foo !tomato", "foo", []string{"!tomato"}},
}

func TestParseCustomCommand(t *testing.T) {
	for _, tt := range parseCustomCommandTests {
		gotCmd, gotArgs := ParseCustomCommand(tt.in)
		if wantCmd, wantArgs := tt.cmd, tt.args; gotCmd != wantCmd || !reflect.DeepEqual(gotArgs, wantArgs) {
			t.Errorf("ParseCustomCommand(%q) = %q, %#v, want %q, %#v", tt.in, gotCmd, gotArgs, wantCmd, wantArgs)
		}
	}
}
