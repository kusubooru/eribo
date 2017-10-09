package rp

import (
	"strings"
	"testing"
)

func TestTieUps(t *testing.T) {
	for _, tt := range tieUps {
		if !strings.HasPrefix(tt, "/me ") {
			t.Errorf("expected '/me ' prefix on %q", tt)
		}
		if got, want := strings.Count(tt, "%s"), 1; got != want {
			t.Errorf("Number of "+`'%s'`+" in phrase = %d, want %d, on %q", got, want, tt)
		}
		if got, want := strings.Count(tt, "[u]"), strings.Count(tt, "[/u]"); got != want {
			t.Errorf("Number of '[u]' = %d, number of '[/u]' = %d, on %q", got, want, tt)
		}
	}
}
