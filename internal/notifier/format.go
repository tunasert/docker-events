package notifier

import (
	"fmt"
	"maps"
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

func formatGroupedEvents(subjectPrefix string, events []docker.Event) (string, string) {
	if len(events) == 0 {
		return "", ""
	}

	if len(events) == 1 {
		return formatEvent(subjectPrefix, events[0])
	}

	prefix := strings.TrimSpace(subjectPrefix)
	if prefix == "" {
		prefix = "Docker events"
	}

	// Get container ID/Actor from first event
	containerID := events[0].ID
	if containerID == "" {
		containerID = events[0].Actor.ID
	}

	// Collect unique actions
	actions := make(map[string]bool)
	for _, event := range events {
		actions[event.Action] = true
	}
	actionList := make([]string, 0, len(actions))
	for action := range actions {
		actionList = append(actionList, action)
	}
	sort.Strings(actionList)

	subject := fmt.Sprintf("%s: %d events for container %s (%s)", prefix, len(events), containerID[:12], strings.Join(actionList, ", "))

	var body strings.Builder
	body.WriteString(fmt.Sprintf("Container: %s\n", containerID))
	body.WriteString(fmt.Sprintf("Event count: %d\n", len(events)))
	body.WriteString(fmt.Sprintf("Time range: %s to %s\n\n",
		events[0].Timestamp.Format(time.RFC3339),
		events[len(events)-1].Timestamp.Format(time.RFC3339)))

	// Get common attributes
	commonAttrs := make(map[string]string)
	if len(events[0].Actor.Attributes) > 0 {
		// Start with first event's attributes
		maps.Copy(commonAttrs, events[0].Actor.Attributes)

		// Keep only attributes that are the same across all events
		for _, event := range events[1:] {
			for k, v := range commonAttrs {
				if eventV, ok := event.Actor.Attributes[k]; !ok || eventV != v {
					delete(commonAttrs, k)
				}
			}
		}
	}

	if len(commonAttrs) > 0 {
		body.WriteString("Common attributes:\n")
		keys := make([]string, 0, len(commonAttrs))
		for key := range commonAttrs {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			body.WriteString(fmt.Sprintf("  - %s=%s\n", key, commonAttrs[key]))
		}
		body.WriteString("\n")
	}

	body.WriteString("Events:\n")
	for i, event := range events {
		body.WriteString(fmt.Sprintf("  %d. [%s] %s %s",
			i+1,
			event.Timestamp.Format("15:04:05"),
			event.Type,
			event.Action))

		if event.Status != "" && event.Status != event.Action {
			body.WriteString(fmt.Sprintf(" (status: %s)", event.Status))
		}
		body.WriteString("\n")
	}

	return subject, strings.TrimSpace(body.String())
}
