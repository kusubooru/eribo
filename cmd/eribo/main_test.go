package main

import (
	"reflect"
	"testing"
)

var splitRoomTitlesTests = []struct {
	in  string
	out []string
}{
	{`["Milk, Cookies & Choco", "Milkshakes"]`, []string{"Milk, Cookies & Choco", "Milkshakes"}},
}

func TestSplitRoomTitles(t *testing.T) {
	for _, tt := range splitRoomTitlesTests {
		got, err := splitRoomTitles(tt.in)
		if err != nil {
			t.Errorf("splitRoomTitles(%q) returned err: %v", tt.in, err)
		}

		if want := tt.out; !reflect.DeepEqual(got, want) {
			t.Errorf("splitRoomTitles(%q) = %q, want: %q", tt.in, got, want)
		}
	}
}
