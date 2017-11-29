package eribo

import (
	"reflect"
	"testing"
)

var parseCommandTests = []struct {
	in  string
	out Command
}{
	{"", CmdUnknown},
	{"!tomato", CmdTomato},
	{"!tomato ", CmdTomato},
	{"!tomato    ", CmdTomato},
	{"!tomato			", CmdTomato},
	{"!tomato			1 2		3", CmdTomato},
	{" !tomato", CmdTomato},
	{"foo !tomato", CmdUnknown},
}

func TestParseCommand(t *testing.T) {
	for _, tt := range parseCommandTests {
		if want, got := tt.out, ParseCommand(tt.in); got != want {
			t.Errorf("ParseCommand(%q) = %q, want %q", tt.in, got, want)
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
