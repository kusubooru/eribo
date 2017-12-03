package flist

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestCharacterData_HasFaveKink(t *testing.T) {
	var tests = []struct {
		kinksMap map[string]string
		kinks    map[string]string
		search   string
		out      bool
	}{
		{kinksMap: map[string]string{"7": "bar"}, kinks: map[string]string{"7": "fave"}, search: "bar", out: true},
		{kinksMap: map[string]string{"7": "bar"}, kinks: map[string]string{"7": "yes"}, search: "bar", out: false},
		{kinksMap: map[string]string{}, kinks: map[string]string{"7": "yes"}, search: "bar", out: false},
		{kinksMap: nil, kinks: map[string]string{"7": "yes"}, search: "bar", out: false},
	}
	for _, tt := range tests {
		char := &CharacterData{Kinks: tt.kinks}
		got := char.HasFaveKink(tt.kinksMap, tt.search)
		want := tt.out
		if got != want {
			t.Errorf("search %s in %v (%v) = %v, want %v", tt.search, tt.kinks, tt.kinksMap, got, want)
		}
	}
}

func TestCharacterData_HasFaveCustomKink(t *testing.T) {
	var tests = []struct {
		kinks  []*CustomKink
		search string
		out    bool
	}{
		{kinks: []*CustomKink{{Name: "foo"}, {Name: "bar", Choice: "fave"}}, search: "ar", out: true},
		{kinks: []*CustomKink{{Name: "foo"}, {Name: "bar", Choice: "yes"}}, search: "ar", out: false},
		{kinks: []*CustomKink{{Name: "foo"}, {Name: "bar"}}, search: "baz", out: false},
	}
	for _, tt := range tests {
		char := &CharacterData{CustomKinks: tt.kinks}
		got := char.HasFaveCustomKink(tt.search)
		want := tt.out
		if got != want {
			t.Errorf("search %s in %v = %v, want %v", tt.search, tt.kinks, got, want)
		}
	}
}

func TestCharacterData_UnmarshalJSON(t *testing.T) {
	const input = `
{
    "id": 1337,
    "name": "John Doe",
    "description": "this is the description",
    "views": 82,
    "customs_first": false,
    "custom_title": "",
    "created_at": 1505346050,
    "updated_at": 1512297470,
    "kinks": {
        "79": "fave"
    },
    "custom_kinks": {
        "18651947": {
            "name": "ck 1",
            "description": "ck 1 description",
            "choice": "fave",
            "children": [
                79
            ]
        },
        "18651948": {
            "name": "ck 2",
            "description": "ck 2 description",
            "choice": "fave",
            "children": []
        },
        "19730982": {
            "name": "ck 1",
            "description": "different ck 1 description",
            "choice": "fave",
            "children": []
        }
    },
    "infotags": {
        "2": "4",
        "3": "1"
    },
    "error": ""
}
`
	r := bytes.NewReader([]byte(input))
	have := new(CharacterData)
	err := json.NewDecoder(r).Decode(have)
	if err != nil {
		t.Error("Unmarshal CharacterData error:", err)
	}
	want := &CharacterData{
		ID:          1337,
		Name:        "John Doe",
		Description: "this is the description",
		Views:       82,
		CreatedAt:   1505346050,
		UpdatedAt:   1512297470,
		Kinks: Kinks{
			"79": "fave",
		},
		Infotags: Infotags{
			"2": "4",
			"3": "1",
		},
		CustomKinks: []*CustomKink{
			{ID: "18651947", Name: "ck 1", Description: "ck 1 description", Choice: "fave"},
			{ID: "18651948", Name: "ck 2", Description: "ck 2 description", Choice: "fave"},
			{ID: "19730982", Name: "ck 1", Description: "different ck 1 description", Choice: "fave"},
		},
	}
	haveCK := []*CustomKink(have.CustomKinks)
	wantCK := []*CustomKink(want.CustomKinks)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Unmarshal CharacterData = \nhave: %#v\nwant: %#v", have, want)
	}
	for i := range haveCK {
		have, want := *haveCK[i], *wantCK[i]
		if have != want {
			t.Errorf("Unmarshal CharacterData CK = \nhave: %#v\nwant: %#v", have, want)
		}
	}
}
