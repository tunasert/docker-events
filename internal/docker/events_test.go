package docker

import (
	"reflect"
	"testing"
	"time"

	"github.com/docker/docker/api/types/events"
)

func TestConvertMessage_TimeFallback(t *testing.T) {
	msg := events.Message{
		ID:       "abc",
		Status:   "start",
		From:     "image:tag",
		Type:     "container",
		Action:   "start",
		Scope:    "local",
		Actor:    events.Actor{ID: "actorid", Attributes: map[string]string{"k": "v"}},
		Time:     123456,
		TimeNano: 0,
	}
	got := convertMessage(msg)
	if got.Timestamp.Before(time.Unix(123456, 0)) || got.Timestamp.After(time.Unix(123456+1, 0)) {
		t.Fatalf("unexpected timestamp: %v", got.Timestamp)
	}
	if got.ID != "abc" {
		t.Fatalf("unexpected id: %s", got.ID)
	}
	if !reflect.DeepEqual(got.Actor.Attributes, map[string]string{"k": "v"}) {
		t.Fatalf("unexpected attributes: %v", got.Actor.Attributes)
	}
}
