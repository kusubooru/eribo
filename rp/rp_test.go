package rp

import (
	"strings"
	"testing"
)

func TestTieUps(t *testing.T) {
	for _, tt := range tieUps {
		checkMePrefix(t, tt)
		checkBBCode(t, tt)
		checkVerbCount(t, tt, "%s", 1)
	}
}

func checkMePrefix(t *testing.T, s string) {
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

func TestTietools(t *testing.T) {
	for _, tt := range tietools {
		s, err := tt.Apply("John Doe")
		if err != nil {
			t.Fatal(err)
		}
		checkMePrefix(t, s)
		checkBBCode(t, s)
		checksSubstringExists(t, s, "John Doe")
	}
}

func TestVonproves(t *testing.T) {
	for _, tt := range vonproves {
		checkMePrefix(t, tt.Raw)
		checkBBCode(t, tt.Raw)
		if tt.HasDate || tt.HasDuration {
			checkVerbCount(t, tt.Raw, "%v", 1)
		}
		if tt.HasUser {
			checkVerbCount(t, tt.Raw, "%s", 1)
		}
	}
}
