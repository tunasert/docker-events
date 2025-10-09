package notifier

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/filippofinke/docker-events/internal/docker"
)

type templateData struct {
	Event       docker.Event
	Events      []docker.Event
	ShortID     string
	Name        string
	Logs        string
	Time        string
	containerID string
	logLines    int64
	dockerCli   *client.Client
	isGrouped   bool
}

// EventCount returns the number of events (for grouped events)
func (t *templateData) EventCount() int {
	if t.isGrouped {
		return len(t.Events)
	}
	return 1
}

// Type returns the event type
func (t *templateData) Type() string {
	return t.Event.Type
}

// Action returns the event action
func (t *templateData) Action() string {
	return t.Event.Action
}

// ID returns the full ID
func (t *templateData) ID() string {
	return t.Event.ID
}

// Status returns the event status
func (t *templateData) Status() string {
	return t.Event.Status
}

// From returns the event from field
func (t *templateData) From() string {
	return t.Event.From
}

// Scope returns the event scope
func (t *templateData) Scope() string {
	return t.Event.Scope
}

// Actor returns the actor
func (t *templateData) Actor() docker.Actor {
	return t.Event.Actor
}

// Attribute returns a specific attribute value
func (t *templateData) Attribute(key string) string {
	if t.Event.Actor.Attributes != nil {
		return t.Event.Actor.Attributes[key]
	}
	return ""
}

// GetLogs fetches container logs if available
func (t *templateData) GetLogs() string {
	if t.Logs != "" {
		return t.Logs
	}

	if t.dockerCli == nil || t.containerID == "" || t.logLines <= 0 {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", t.logLines),
	}

	logs, err := t.dockerCli.ContainerLogs(ctx, t.containerID, options)
	if err != nil {
		return fmt.Sprintf("[error fetching logs: %v]", err)
	}
	defer logs.Close()

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, logs)

	t.Logs = strings.TrimSpace(buf.String())
	return t.Logs
}

func formatEventWithTemplate(tmplStr string, event docker.Event, dockerCli *client.Client, logLines int) (string, string, error) {
	// Extract container name from attributes
	containerName := ""
	containerID := ""

	if event.Type == "container" {
		if name, ok := event.Actor.Attributes["name"]; ok {
			containerName = name
		}
		containerID = event.Actor.ID
	}

	// Create short ID (first 12 characters)
	shortID := event.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	data := &templateData{
		Event:       event,
		ShortID:     shortID,
		Name:        containerName,
		Time:        event.Timestamp.Format(time.RFC3339),
		containerID: containerID,
		logLines:    int64(logLines),
		dockerCli:   dockerCli,
	}

	// Parse template
	tmpl, err := template.New("message").Parse(tmplStr)
	if err != nil {
		return "", "", fmt.Errorf("parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("execute template: %w", err)
	}

	return strings.TrimSpace(buf.String()), "", nil
}

func formatGroupedEventsWithTemplate(tmplStr string, events []docker.Event, dockerCli *client.Client, logLines int) (string, string, error) {
	if len(events) == 0 {
		return "", "", fmt.Errorf("no events to format")
	}

	if len(events) == 1 {
		return formatEventWithTemplate(tmplStr, events[0], dockerCli, logLines)
	}

	// Use the first event as the primary event for template data
	firstEvent := events[0]

	// Extract container name from attributes
	containerName := ""
	containerID := ""

	if firstEvent.Type == "container" {
		if name, ok := firstEvent.Actor.Attributes["name"]; ok {
			containerName = name
		}
		containerID = firstEvent.Actor.ID
	}

	// Create short ID (first 12 characters)
	shortID := firstEvent.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	data := &templateData{
		Event:       firstEvent,
		Events:      events,
		ShortID:     shortID,
		Name:        containerName,
		Time:        firstEvent.Timestamp.Format(time.RFC3339),
		containerID: containerID,
		logLines:    int64(logLines),
		dockerCli:   dockerCli,
		isGrouped:   true,
	}

	// Parse template
	tmpl, err := template.New("message").Parse(tmplStr)
	if err != nil {
		return "", "", fmt.Errorf("parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("execute template: %w", err)
	}

	return strings.TrimSpace(buf.String()), "", nil
}
