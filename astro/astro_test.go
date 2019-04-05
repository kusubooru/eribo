package astro

import (
	"database/sql/driver"
	"testing"
)

func TestSignValue(t *testing.T) {
	tests := []struct {
		sign Sign
		want driver.Value
	}{
		{Aries, "aries"},
		{"invalid sign", "invalid sign"},
	}

	for _, tt := range tests {
		v, err := tt.sign.Value()
		if err != nil {
			t.Errorf("%q.Value() returned error: %v", tt.sign, err)
		}
		if got, want := v, tt.want; got != want {
			t.Errorf("%q.Value() => %q, want %q", tt.sign, got, want)
		}
	}
}

func TestSignScan(t *testing.T) {
	tests := []struct {
		in   interface{}
		want Sign
	}{
		{"aries", Aries},
		{"invalid sign", "invalid sign"},
	}

	for _, tt := range tests {
		var s Sign
		if err := s.Scan(tt.in); err != nil {
			t.Errorf("Scan(%v) returned error: %v", tt.in, err)
		}
		if got, want := s, tt.want; got != want {
			t.Errorf("Scan(%q) => %q, want %q", tt.in, got, want)
		}
	}
}

func TestSignScanNil(t *testing.T) {
	var s Sign
	if err := s.Scan(nil); err == nil {
		t.Errorf("Scan nil to Sign expected to return error")
	}
}
