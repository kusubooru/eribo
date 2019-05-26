package rp

import (
	"reflect"
	"strings"
	"testing"
)

func TestTietoolMarshalText(t *testing.T) {
	tool := Tietool{
		name:    "foo",
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
		name:    "foo",
		Quality: Rare,
		Desc:    tmplMust(`/me gives {{.Tool}} to {{.User}}`),
	}
	if got := tool; !reflect.DeepEqual(got, want) {
		t.Errorf("UnmarshalText \nhave: %#v\nwant: %#v", got, want)
	}
}

func TestTietoolsApply(t *testing.T) {
	user := "bob"
	for _, tool := range tietools {
		msg, err := tool.Apply(user)
		if err != nil {
			t.Fatalf("applying tool %v, returned err: %v", tool, err)
		}
		if !strings.Contains(msg, user) {
			t.Errorf("applying user %q on tool %+v, message = %q, want user in message", user, tool, msg)
		}
		if !strings.Contains(msg, tool.Name()) {
			t.Errorf("applying user %q on tool %+v, message = %q, want %s in message", user, tool, msg, tool.Name())
		}
	}
}
