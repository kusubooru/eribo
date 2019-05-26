package astro

import (
	"strings"
	"testing"
)

func TestForInvalidSign(t *testing.T) {
	h, err := For("", "foo")
	if err != nil {
		t.Errorf("For returned err: %v", err)
	}
	if !strings.Contains(h, "signs are") {
		t.Errorf("For with invalid sign should return no error and a message with the valid signs, got %s", h)
	}
}
