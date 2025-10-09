package docker

import (
	"time"

	"github.com/docker/docker/api/types/events"
)

type Event struct {
	ID        string
	Status    string
	From      string
	Type      string
	Action    string
	Scope     string
	Actor     Actor
	Timestamp time.Time
}

func (e Event) Time() string {
	return e.Timestamp.Format(time.RFC3339)
}

type Actor struct {
	ID         string
	Attributes map[string]string
}

func convertMessage(msg events.Message) Event {
	timestamp := time.Unix(0, msg.TimeNano)
	if msg.TimeNano == 0 && msg.Time != 0 {
		timestamp = time.Unix(msg.Time, 0)
	}

	return Event{
		ID:     msg.ID,
		Status: string(msg.Status),
		From:   msg.From,
		Type:   string(msg.Type),
		Action: string(msg.Action),
		Scope:  string(msg.Scope),
		Actor: Actor{
			ID:         msg.Actor.ID,
			Attributes: msg.Actor.Attributes,
		},
		Timestamp: timestamp,
	}
}
