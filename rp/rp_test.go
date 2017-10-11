package rp

import (
	"strings"
	"testing"
)

func TestTieUps(t *testing.T) {
	for _, tt := range tieUps {
		checkEmote(t, tt)
		checkBBcode(t, tt)
		checkVerbCount(t, tt, "%s", 1)
	}
}

func checkEmote(t *testing.T, s string) {
	t.Helper()
	if !strings.HasPrefix(s, "/me ") {
		t.Errorf("expected '/me ' prefix, message is: %q", s)
	}
}

func checkBBcode(t *testing.T, s string) {
	t.Helper()
	if got, want := strings.Count(s, "[u]"), strings.Count(s, "[/u]"); got != want {
		t.Errorf("Number of '[u]' = %d, number of '[/u]' = %d, message is: %q", got, want, s)
	}
}

func checkVerbCount(t *testing.T, s, verb string, n int) {
	t.Helper()
	if got, want := strings.Count(s, verb), n; got != want {
		t.Errorf("Number of %q = %d, want %d, message is: %q", verb, got, want, s)
	}
}

func TestTktools(t *testing.T) {
	for _, tt := range tktools {
		checkEmote(t, tt.Raw)
		checkBBcode(t, tt.Raw)
		if tt.Colors != nil {
			checkVerbCount(t, tt.Raw, "%s", 2)
		} else {
			checkVerbCount(t, tt.Raw, "%s", 1)
		}
	}
}

func TestVonproves(t *testing.T) {
	for _, tt := range vonproves {
		checkEmote(t, tt.Raw)
		checkBBcode(t, tt.Raw)
		if tt.HasDate || tt.HasDuration {
			checkVerbCount(t, tt.Raw, "%v", 1)
		}
		if tt.HasUser {
			checkVerbCount(t, tt.Raw, "%s", 1)
		}
	}
}
