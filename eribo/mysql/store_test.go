package mysql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kusubooru/eribo/eribo"
)

func TestAddMessageWithURLs(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	created := time.Now().UTC().Add(time.Second).Truncate(timeTruncate)
	m := &eribo.Message{
		Channel: "foo",
		Player:  "bar",
		Message: "baz",
		Created: created,
	}
	urls := []string{"http://url1", "http://url2"}
	if err := s.AddMessageWithURLs(m, urls); err != nil {
		t.Fatal("AddMessageWithURLs failed:", err)
	}

	images, err := s.GetImages()
	if err != nil {
		t.Fatal("GetImages failed:", err)
	}
	if got, want := len(images), 2; got != want {
		t.Fatalf("GetImages produced %d results, want %d", got, want)
	}
	want := []*eribo.Image{
		{ID: 1, URL: "http://url1", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
		{ID: 2, URL: "http://url2", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
	}

	if have := images; !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("AddImages = \nhave: %#v\nwant: %#v", have, want)
	}
}

func TestAddFeedback(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	created := time.Now().UTC().Add(1 * time.Second).Truncate(timeTruncate)
	f := &eribo.Feedback{
		Player:  "foo",
		Message: "bar",
		Created: created,
	}
	if err := s.AddFeedback(f); err != nil {
		t.Fatal("AddFeedback failed:", err)
	}

	feedback, err := s.GetAllFeedback(10, 0)
	if err != nil {
		t.Fatal("GetFeedback failed:", err)
	}
	if got, want := len(feedback), 1; got != want {
		t.Fatalf("GetFeedback produced %d results, want %d", got, want)
	}
	want := []*eribo.Feedback{
		{ID: 1, Player: "foo", Message: "bar", Created: created},
	}

	if have := feedback; !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("AddFeedback = \nhave: %#v\nwant: %#v", have, want)
	}
}

func TestGetRecentFeedback(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	created1 := time.Now().UTC().Add(1 * time.Second).Truncate(timeTruncate)
	created2 := time.Now().UTC().Add(2 * time.Second).Truncate(timeTruncate)
	created3 := time.Now().UTC().Add(3 * time.Second).Truncate(timeTruncate)

	fb := []*eribo.Feedback{
		{Player: "foo", Message: "bar", Created: created1},
		{Player: "foo", Message: "bar", Created: created2},
		{Player: "foo", Message: "bar", Created: created3},
	}
	for _, f := range fb {
		if err := s.AddFeedback(f); err != nil {
			t.Fatal("AddFeedback failed:", err)
		}
	}

	feedback, err := s.GetRecentFeedback(2, 0)
	if err != nil {
		t.Fatal("GetFeedback failed:", err)
	}

	want := []*eribo.Feedback{
		{ID: 3, Player: "foo", Message: "bar", Created: created3},
		{ID: 2, Player: "foo", Message: "bar", Created: created2},
	}

	if have := feedback; !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("AddFeedback = \nhave: %#v\nwant: %#v", have, want)
	}
}
