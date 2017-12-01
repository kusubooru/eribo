package mysql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kusubooru/eribo/eribo"
	"github.com/kusubooru/eribo/flist"
)

func TestLog(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	e := &eribo.Event{
		Command: eribo.CmdTomato,
		Player:  "foo",
	}
	if err := s.Log(e); err != nil {
		t.Fatal("Log failed:", err)
	}

	have, err := s.GetLog(1)
	if err != nil {
		t.Fatal("GetLog failed:", err)
	}
	want := &eribo.Event{ID: 1, Player: "foo", Command: eribo.CmdTomato}

	// ignore created
	have.Created = time.Time{}

	if !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("Log = \nhave: %#v\nwant: %#v", have, want)
	}
}

func TestGetRecentLogs(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	e := &eribo.Event{
		Command: eribo.CmdTomato,
		Player:  "foo",
	}
	for i := 0; i < 3; i++ {
		if err := s.Log(e); err != nil {
			t.Fatal("Log failed:", err)
		}
	}

	have, err := s.GetRecentLogs(2, 0)
	if err != nil {
		t.Fatal("GetRecentLogs failed:", err)
	}
	want := []*eribo.Event{
		{ID: 3, Player: "foo", Command: eribo.CmdTomato},
		{ID: 2, Player: "foo", Command: eribo.CmdTomato},
	}

	// ignore created
	for _, h := range have {
		h.Created = time.Time{}
	}

	if !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("GetRecentLogs = \nhave: %#v\nwant: %#v", have, want)
	}
}

func TestGetRecentLothLogs(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	issuer := "jin"
	channel := "2ch"
	loth := eribo.NewLoth(&eribo.Player{Name: "foo", Role: flist.RoleSwitch, Status: flist.StatusOnline}, 1*time.Hour)
	isNew := true
	targets := []*eribo.Player{{Name: "bar"}, {Name: "baz"}}
	expires := loth.Expires
	created1 := time.Now().UTC().Add(1 * time.Second).Truncate(1 * time.Microsecond)
	created2 := time.Now().UTC().Add(2 * time.Second).Truncate(1 * time.Microsecond)
	created3 := time.Now().UTC().Add(3 * time.Second).Truncate(1 * time.Microsecond)

	logs := []*eribo.LothLog{
		{Issuer: issuer, Channel: channel, Loth: loth, IsNew: isNew, Targets: targets, Created: created1},
		{Issuer: issuer, Channel: channel, Loth: loth, IsNew: isNew, Targets: targets, Created: created2},
		{Issuer: issuer, Channel: channel, Loth: loth, IsNew: isNew, Targets: targets, Created: created3},
	}

	for _, lothLog := range logs {
		if err := s.LogLoth(lothLog); err != nil {
			t.Fatal("LogLoth failed:", err)
		}
	}

	have, err := s.GetRecentLothLogs(2, 0)
	if err != nil {
		t.Fatal("GetRecentLothLogs failed:", err)
	}
	want := []*eribo.LothLog{
		{ID: 3, Issuer: "jin", Channel: "2ch", Created: created3,
			Loth: &eribo.Loth{
				Player:  &eribo.Player{Name: "foo", Role: flist.RoleSwitch, Status: flist.StatusOnline},
				Expires: expires,
			},
			IsNew: true, Targets: targets},
		{ID: 2, Issuer: "jin", Channel: "2ch", Created: created2,
			Loth: &eribo.Loth{
				Player:  &eribo.Player{Name: "foo", Role: flist.RoleSwitch, Status: flist.StatusOnline},
				Expires: expires,
			},
			IsNew: true, Targets: targets},
	}

	if !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("GetRecentLothLogs = \nhave: %#v\nwant: %#v", have, want)
	}
}

func TestLogLoth_unableToFindEligibleTarget(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	issuer := "jin"
	channel := "2ch"
	isNew := false
	targets := []*eribo.Player{{Name: "bar"}, {Name: "baz"}}
	created1 := time.Now().UTC().Add(1 * time.Second).Truncate(1 * time.Microsecond)
	created2 := time.Now().UTC().Add(2 * time.Second).Truncate(1 * time.Microsecond)
	created3 := time.Now().UTC().Add(3 * time.Second).Truncate(1 * time.Microsecond)

	logs := []*eribo.LothLog{
		{Issuer: issuer, Channel: channel, IsNew: isNew, Targets: targets, Created: created1},
		{Issuer: issuer, Channel: channel, IsNew: isNew, Targets: targets, Created: created2},
		{Issuer: issuer, Channel: channel, IsNew: isNew, Targets: targets, Created: created3},
	}

	for _, lothLog := range logs {
		if err := s.LogLoth(lothLog); err != nil {
			t.Fatal("LogLoth failed:", err)
		}
	}

	have, err := s.GetRecentLothLogs(2, 0)
	if err != nil {
		t.Fatal("GetRecentLothLogs failed:", err)
	}
	want := []*eribo.LothLog{
		{ID: 3, Issuer: "jin", Channel: "2ch", Created: created3,
			Loth:  &eribo.Loth{Player: &eribo.Player{}, Expires: created3},
			IsNew: false, Targets: targets},
		{ID: 2, Issuer: "jin", Channel: "2ch", Created: created2,
			Loth:  &eribo.Loth{Player: &eribo.Player{}, Expires: created2},
			IsNew: false, Targets: targets},
	}

	if !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("GetRecentLothLogs = \nhave: %#v\nwant: %#v", have, want)
	}
}
