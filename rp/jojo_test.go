package rp

import "testing"

func TestStand_apply(t *testing.T) {
	var tests = []struct {
		stand stand
		user  string
		out   string
	}{
		{
			stand{Name: "The World", Type: "Close-Range Stand", Desc: "ZAWARUDO!"},
			"Dio",
			"Dio's new Stand is [u]The World[/u] ([i]Close-Range Stand[/i]): ZAWARUDO!",
		},
	}

	for _, tt := range tests {
		got := tt.stand.apply(tt.user)
		want := tt.out
		if got != want {
			t.Errorf("stand apply on user %q = \nhave: %q\nwant: %q", tt.user, got, want)
		}
	}
}
