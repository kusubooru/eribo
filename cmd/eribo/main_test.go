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

func TestArgsContain(t *testing.T) {
	var tests = []struct {
		inArgs  []string
		in      string
		outArgs []string
		out     bool
	}{
		{[]string{"10", "0", "desc"}, "desc", []string{"10", "0"}, true},
		{[]string{"10", "0"}, "desc", []string{"10", "0"}, false},
		{[]string{}, "desc", []string{}, false},
		{nil, "desc", nil, false},
	}

	for _, tt := range tests {
		args, ok := argsContain(tt.inArgs, tt.in)
		if got, want := ok, tt.out; got != want {
			t.Errorf("popArgs(%v, %v) = %v, want %v", tt.inArgs, tt.in, got, want)
		}
		if got, want := args, tt.outArgs; !reflect.DeepEqual(got, want) {
			t.Errorf("\nhave: %v\nwant: %v", got, want)
		}
	}
}

func TestArgsPopAtoi(t *testing.T) {
	var tests = []struct {
		inArgs  []string
		n       int
		outArgs []string
		ok      bool
	}{
		{[]string{"10", "0", "desc"}, 10, []string{"0", "desc"}, true},
		{[]string{"desc", "10", "0"}, 10, []string{"desc", "0"}, true},
		{[]string{"desc"}, 0, []string{"desc"}, false},
		{[]string{}, 0, []string{}, false},
		{nil, 0, nil, false},
	}

	for _, tt := range tests {
		n, args, ok := argsPopAtoi(tt.inArgs)
		if got, want := n, tt.n; got != want {
			t.Errorf("popArgs(%v) = %v, want %v", tt.inArgs, got, want)
		}
		if got, want := ok, tt.ok; got != want {
			t.Errorf("popArgs(%v) = %v, want %v", tt.inArgs, got, want)
		}
		if got, want := args, tt.outArgs; !reflect.DeepEqual(got, want) {
			t.Errorf("\nhave: %v\nwant: %v", got, want)
		}
	}
}

func TestArgsPopAtoiDefault(t *testing.T) {
	var tests = []struct {
		inArgs    []string
		inDefault int
		n         int
		outArgs   []string
	}{
		{[]string{"10", "0", "desc"}, 7, 10, []string{"0", "desc"}},
		{[]string{"desc", "10", "0"}, 7, 10, []string{"desc", "0"}},
		{[]string{"desc"}, 7, 7, []string{"desc"}},
		{[]string{}, 0, 0, []string{}},
		{nil, 7, 7, nil},
	}

	for _, tt := range tests {
		n, args := argsPopAtoiDefault(tt.inArgs, tt.inDefault)
		if got, want := n, tt.n; got != want {
			t.Errorf("popArgs(%v) = %v, want %v", tt.inArgs, got, want)
		}
		if got, want := args, tt.outArgs; !reflect.DeepEqual(got, want) {
			t.Errorf("\nhave: %v\nwant: %v", got, want)
		}
	}
}
