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

	have, err := s.GetImages(5, 0, false, false)
	if err != nil {
		t.Fatal("GetImages failed:", err)
	}
	want := []*eribo.Image{
		{ID: 1, URL: "http://url1", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
		{ID: 2, URL: "http://url2", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
	}
	deepEqual(t, have, want, "AddImages = ")
}

func TestGetImages(t *testing.T) {
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

	created2 := created.Add(time.Second)
	m.Created = created2
	urls = []string{"http://url3"}
	if err := s.AddMessageWithURLs(m, urls); err != nil {
		t.Fatal("AddMessageWithURLs second failed:", err)
	}

	have, err := s.GetImages(5, 0, true, false)
	if err != nil {
		t.Fatal("GetImages failed:", err)
	}
	want := []*eribo.Image{
		{ID: 3, URL: "http://url3", MessageID: 2, Created: created2,
			Message: &eribo.Message{ID: 2, Channel: "foo", Player: "bar", Message: "baz", Created: created2},
		},
		{ID: 1, URL: "http://url1", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
		{ID: 2, URL: "http://url2", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
	}
	deepEqual(t, have, want, "GetImages = ")
}

func TestSetImageKuid(t *testing.T) {
	s := setup(t)
	defer teardown(t, s)

	created := time.Now().UTC().Truncate(timeTruncate)
	m := &eribo.Message{Channel: "foo", Player: "bar", Message: "baz", Created: created}
	urls := []string{"http://url"}
	if err := s.AddMessageWithURLs(m, urls); err != nil {
		t.Fatal("AddMessageWithURLs failed:", err)
	}

	if err := s.SetImageKuid(1, 1337); err != nil {
		t.Fatal("setting image kuid failed:", err)
	}

	got, err := s.GetImage(1)
	if err != nil {
		t.Fatal("getting image failed:", err)
	}

	want := &eribo.Image{
		ID:        1,
		URL:       "http://url",
		Done:      true,
		Kuid:      1337,
		Created:   created,
		MessageID: 1,
		Message:   &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
	}
	deepEqual(t, got, want, "setting image kuid")
}

func TestToggleImageDone(t *testing.T) {
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

	if err := s.ToggleImageDone(2); err != nil {
		t.Fatal("toggling image done failed:", err)
	}

	have, err := s.GetImages(5, 0, false, false)
	if err != nil {
		t.Fatal("GetImages failed:", err)
	}
	want := []*eribo.Image{
		{ID: 1, URL: "http://url1", MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
		{ID: 2, URL: "http://url2", Done: true, MessageID: 1, Created: created,
			Message: &eribo.Message{ID: 1, Channel: "foo", Player: "bar", Message: "baz", Created: created},
		},
	}
	deepEqual(t, have, want, "toggling image done")
}

func deepEqual(t *testing.T, have, want interface{}, message string) {
	t.Helper()
	if !reflect.DeepEqual(have, want) {
		data, _ := json.Marshal(have)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
		t.Fatalf("%s\nhave: %#v\nwant: %#v", message, have, want)
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

	have, err := s.GetAllFeedback(10, 0)
	if err != nil {
		t.Fatal("GetFeedback failed:", err)
	}
	want := []*eribo.Feedback{
		{ID: 1, Player: "foo", Message: "bar", Created: created},
	}
	deepEqual(t, have, want, "AddFeedback = ")
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

	have, err := s.GetRecentFeedback(2, 0)
	if err != nil {
		t.Fatal("GetFeedback failed:", err)
	}

	want := []*eribo.Feedback{
		{ID: 3, Player: "foo", Message: "bar", Created: created3},
		{ID: 2, Player: "foo", Message: "bar", Created: created2},
	}
	deepEqual(t, have, want, "getting recent feedback")
}
