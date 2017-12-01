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

	m := &eribo.Message{
		Channel: "foo",
		Player:  "bar",
		Message: "baz",
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
		{ID: 1, URL: "http://url1", MessageID: 1, Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz"}},
		{ID: 2, URL: "http://url2", MessageID: 1, Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz"}},
	}
	// ignore created
	for _, img := range images {
		img.Created = time.Time{}
		img.Message.Created = time.Time{}
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

	f := &eribo.Feedback{
		Player:  "foo",
		Message: "bar",
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
		{ID: 1, Player: "foo", Message: "bar"},
	}
	// ignore created
	for _, f := range feedback {
		f.Created = time.Time{}
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

	f := &eribo.Feedback{
		Player:  "foo",
		Message: "bar",
	}
	for i := 0; i < 3; i++ {
		if err := s.AddFeedback(f); err != nil {
			t.Fatal("AddFeedback failed:", err)
		}
	}

	feedback, err := s.GetRecentFeedback(2, 0)
	if err != nil {
		t.Fatal("GetFeedback failed:", err)
	}

	want := []*eribo.Feedback{
		{ID: 3, Player: "foo", Message: "bar"},
		{ID: 2, Player: "foo", Message: "bar"},
	}
	// ignore created
	for _, f := range feedback {
		f.Created = time.Time{}
	}
	if have := feedback; !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("AddFeedback = \nhave: %#v\nwant: %#v", have, want)
	}
}
