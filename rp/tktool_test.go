package rp

import (
	"strings"
	"testing"

	"github.com/kusubooru/eribo/loot"
)

func testTktoolApply(t *testing.T, tool Tktool, user string) string {
	t.Helper()
	s, err := tool.Apply(user)
	if err != nil {
		t.Errorf("applying %q to %q returned error: %v", user, tool.Emote.Root, err)
	}
	return s
}

func TestTktools(t *testing.T) {
	for _, tool := range tktools {
		checkSyntax(t, tool.Emote.Root.String(), "{", "}")
		s := testTktoolApply(t, tool, "Bob")
		checkMePrefix(t, s)
		checkBBCode(t, s)
	}
}

func TestTktool_Apply(t *testing.T) {
	tests := []struct {
		tool Tktool
		want string
	}{
		{
			Tktool{name: "tool", Quality: Common, Emote: tmplMust("/me hands {{.User}} a {{.Tool}}.")},
			"/me hands Bob a [color=white]tool[/color].",
		},
		{ // Tool with quality unknown must still produce white color.
			Tktool{name: "tool", Emote: tmplMust("/me hands {{.User}} a {{.Tool}}.")},
			"/me hands Bob a [color=white]tool[/color].",
		},
		{
			Tktool{name: "tool", Quality: Poor, Colors: []Color{Blue}, Emote: tmplMust("/me hands {{.User}} a {{.Color}} {{.Tool}}.")},
			"/me hands Bob a blue [color=gray]tool[/color].",
		},
		{
			Tktool{name: "tool", Quality: Uncommon, Colors: []Color{Blue}, Emote: tmplMust("/me hands {{.User}} a {{.Tool}}.")},
			"/me hands Bob a [color=green]tool[/color].",
		},
	}

	for _, tt := range tests {
		got, err := tt.tool.Apply("Bob")
		if err != nil {
			t.Fatal(err)
		}
		if want := tt.want; got != want {
			t.Errorf("apply %+v\nhave: %v\nwant: %v", tt.tool, got, want)
		}
	}
}

func TestTktoolsApply(t *testing.T) {
	user := "bob"
	for _, tool := range tktools {
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

func TestTktoolsLootTableLegendaries(t *testing.T) {

	tableOneLego := loot.NewTable(
		[]loot.Drop{
			{Item: Tktool{Quality: Legendary}, Weight: 0},
			{Item: Tktool{Quality: Epic}, Weight: 1},
		},
	)
	tests := []struct {
		name string
		t    *TktoolsLootTable
		want int
	}{
		{"1 lego", &TktoolsLootTable{tableOneLego}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Legendaries(); got != tt.want {
				t.Errorf("TktoolsLootTable.Legendaries() = %v, want %v", got, tt.want)
			}
		})
	}
}
