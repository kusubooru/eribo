package rp

import (
	"reflect"
	"testing"
)

func TestTietoolMarshalText(t *testing.T) {
	tool := Tietool{
		Name:    "foo",
		Quality: Rare,
		Desc: tmplMust(`/me gives a rare, amazing, fantastic {{.Tool}} to
		{{.User}}`),
	}
	got, err := tool.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}
	want := `foo
rare
/me gives a rare, amazing, fantastic {{.Tool}} to {{.User}}
`
	if got != want {
		t.Errorf("MarshalText \nhave: %q\nwant: %q", got, want)
	}
}

func TestTietoolUnmarshalText(t *testing.T) {
	in := `foo
rare
/me gives {{.Tool}} to {{.User}}
`
	tool := new(Tietool)
	err := tool.UnmarshalText(in)
	if err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}

	want := &Tietool{
		Name:    "foo",
		Quality: Rare,
		Desc:    tmplMust(`/me gives {{.Tool}} to {{.User}}`),
	}
	if got := tool; !reflect.DeepEqual(got, want) {
		t.Errorf("UnmarshalText \nhave: %#v\nwant: %#v", got, want)
	}
}
