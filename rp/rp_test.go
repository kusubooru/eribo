package rp

import (
	"strings"
	"testing"
)

func TestTieUps(t *testing.T) {
	for _, tt := range tieUps {
		checkEmote(t, tt)
		checkBBCode(t, tt)
		checkVerbCount(t, tt, "%s", 1)
	}
}

func checkEmoteTktool(t *testing.T, tool Tktool) {
	t.Helper()
	checkEmote(t, tool.Poor)
	checkEmote(t, tool.Common)
	checkEmote(t, tool.Uncommon)
}

func checkEmote(t *testing.T, s string) {
	t.Helper()
	if !strings.HasPrefix(s, "/me ") {
		t.Errorf("expected '/me ' prefix, message is: %q", s)
	}
}

func checkBBCode(t *testing.T, s string) {
	t.Helper()
	checkSyntax(t, s, "[", "]")
	checkSyntax(t, s, "[u]", "[/u]")
	checkSyntax(t, s, "[color=", "[/color]")
}

func checkActionsTietool(t *testing.T, s string) {
	t.Helper()
	checkSyntax(t, s, "{", "}")
	checkSyntax(t, s, "{", "}")
}

func checkBBcodeTktool(t *testing.T, tool Tktool) {
	t.Helper()
	checkBBCode(t, tool.Poor)
	checkBBCode(t, tool.Common)
	checkBBCode(t, tool.Uncommon)
}

func checkSyntax(t *testing.T, s, start, end string) {
	t.Helper()
	if got, want := strings.Count(s, start), strings.Count(s, end); got != want {
		t.Errorf("Number of %q = %d, number of %q = %d, message is: %q", start, got, end, want, s)
	}
}

func checkVerbCount(t *testing.T, s, verb string, n int) {
	t.Helper()
	if got, want := strings.Count(s, verb), n; got != want {
		t.Errorf("Number of %q = %d, want %d, message is: %q", verb, got, want, s)
	}
}

func checksSubstringExists(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("Expected %q to be in message %q", substr, s)
	}
}

func checkVerbCountTktool(t *testing.T, tool Tktool) {
	t.Helper()
	if tool.Colors != nil {
		checkVerbCount(t, tool.Poor, "%s", 2)
		checkVerbCount(t, tool.Common, "%s", 2)
		checkVerbCount(t, tool.Uncommon, "%s", 2)
	} else {
		checkVerbCount(t, tool.Poor, "%s", 1)
		checkVerbCount(t, tool.Common, "%s", 1)
		checkVerbCount(t, tool.Uncommon, "%s", 1)
	}
}

func checkQualityColorCode(t *testing.T, tool Tktool) {
	t.Helper()
	checkSyntax(t, tool.Poor, "[color=gray]", "[/color]")
	checkSyntax(t, tool.Common, "[color=white]", "[/color]")
	checkSyntax(t, tool.Uncommon, "[color=green]", "[/color]")
}

func TestTktools(t *testing.T) {
	for _, tt := range tktools {
		checkEmoteTktool(t, tt)
		checkBBcodeTktool(t, tt)
		checkVerbCountTktool(t, tt)
		checkQualityColorCode(t, tt)
	}
}

func TestTietools(t *testing.T) {
	for _, tt := range tietools {
		s, err := tt.Apply("John Doe")
		if err != nil {
			t.Fatal(err)
		}
		checkEmote(t, s)
		checkBBCode(t, s)
		checksSubstringExists(t, s, "John Doe")
	}
}

func TestVonproves(t *testing.T) {
	for _, tt := range vonproves {
		checkEmote(t, tt.Raw)
		checkBBCode(t, tt.Raw)
		if tt.HasDate || tt.HasDuration {
			checkVerbCount(t, tt.Raw, "%v", 1)
		}
		if tt.HasUser {
			checkVerbCount(t, tt.Raw, "%s", 1)
		}
	}
}

func TestTktool_Apply(t *testing.T) {
	var tktoolApplyTests = []struct {
		q    Quality
		user string
		tool Tktool
		out  string
	}{
		{
			q:    Common,
			user: "Bob",
			tool: Tktool{Common: "hand %s a common tool"},
			out:  "hand Bob a common tool",
		},
		{
			q:    Common,
			user: "Bob",
			tool: Tktool{Common: "hand %s a common %s tool", Colors: []Color{Red}},
			out:  "hand Bob a common red tool",
		},
		{
			q:    Poor,
			user: "Bob",
			tool: Tktool{Poor: "hand %s a poor tool", Common: "hand %s a common tool"},
			out:  "hand Bob a poor tool",
		},
		{
			q:    Unknown,
			user: "Bob",
			tool: Tktool{Poor: "hand %s a poor tool", Common: "hand %s a common tool"},
			out:  "hand Bob a common tool",
		},
		{
			q:    Uncommon,
			user: "Bob",
			tool: Tktool{Poor: "hand %s a poor tool", Common: "hand %s a common tool", Uncommon: "hand %s an uncommon tool"},
			out:  "hand Bob an uncommon tool",
		},
	}
	for _, tt := range tktoolApplyTests {
		got, want := tt.tool.Apply(tt.user, tt.q), tt.out
		if got != want {
			t.Errorf("tktool.Apply(%q, %v) = %q, want %q", tt.user, tt.q, got, want)
		}
	}
}
