package notifier

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/filippofinke/docker-events/internal/docker"
)

func formatEvent(subjectPrefix string, event docker.Event) (string, string) {
	prefix := strings.TrimSpace(subjectPrefix)
	if prefix == "" {
		prefix = "Docker event"
	}

	subject := fmt.Sprintf("%s: %s %s", prefix, event.Type, event.Action)
	if event.Actor.ID != "" {
		subject = fmt.Sprintf("%s (%s)", subject, event.Actor.ID)
	}

	var body strings.Builder
	body.WriteString(fmt.Sprintf("Time: %s\n", event.Timestamp.Format(time.RFC3339)))
	if event.Status != "" {
		body.WriteString(fmt.Sprintf("Status: %s\n", event.Status))
	}
	if event.From != "" {
		body.WriteString(fmt.Sprintf("From: %s\n", event.From))
	}
	if event.Scope != "" {
		body.WriteString(fmt.Sprintf("Scope: %s\n", event.Scope))
	}
	if event.ID != "" {
		body.WriteString(fmt.Sprintf("ID: %s\n", event.ID))
	}
	if event.Actor.ID != "" {
		body.WriteString(fmt.Sprintf("Actor: %s\n", event.Actor.ID))
	}

	if len(event.Actor.Attributes) > 0 {
		body.WriteString("Attributes:\n")
		keys := make([]string, 0, len(event.Actor.Attributes))
		for key := range event.Actor.Attributes {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			body.WriteString(fmt.Sprintf("  - %s=%s\n", key, event.Actor.Attributes[key]))
		}
	}

	return subject, strings.TrimSpace(body.String())
}
