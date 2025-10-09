# Example configurations for MESSAGE_TEMPLATE

# 1. Simple format with container name and action

MESSAGE_TEMPLATE="Container {{.Name}} {{.Action}}"

# Output: Container my-app-1 start

# 2. With timestamp and short ID

MESSAGE_TEMPLATE="[{{.Time}}] {{.Type}} {{.Action}}: {{.Name}} ({{.ShortID}})"

# Output: [2025-10-09T07:37:05Z] container start: my-app-1 (cd280ca744b1)

# 3. Detailed format with image and attributes

MESSAGE_TEMPLATE="Container: {{.Name}}\nAction: {{.Action}}\nImage: {{.From}}\nProject: {{.Attribute \"com.docker.compose.project\"}}\nTime: {{.Time}}"

# Output:

# Container: my-app-1

# Action: start

# Image: my-app-image

# Project: my-project

# Time: 2025-10-09T07:37:05Z

# 4. With container logs (last 20 lines)

MESSAGE_TEMPLATE="Container {{.Name}} {{.Action}}\n\nLogs:\n{{.GetLogs}}"
MESSAGE_LOG_LINES=20

# Output:

# Container my-app-1 start

#

# Logs:

# [2025-10-09 07:37:05] Application starting...

# [2025-10-09 07:37:06] Server listening on port 8080

# ...

# 5. Minimal format

MESSAGE_TEMPLATE="{{.Type}}/{{.Action}}: {{if .Name}}{{.Name}}{{else}}{{.ShortID}}{{end}}"

# Output: container/start: my-app-1

# 6. Slack/Discord friendly format with emojis

MESSAGE_TEMPLATE="🐳 {{.Type}} {{.Action}}\n📦 **Container:** {{.Name}}\n🖼️ **Image:** {{.From}}\n⏰ **Time:** {{.Time}}"

# Output:

# 🐳 container start

# 📦 **Container:** my-app-1

# 🖼️ **Image:** my-app-image

# ⏰ **Time:** 2025-10-09T07:37:05Z

# 7. Conditional format (shows name if available, otherwise short ID)

MESSAGE_TEMPLATE="{{if .Name}}{{.Name}}{{else}}{{.ShortID}}{{end}} - {{.Action}} ({{.Status}})"

# Output: my-app-1 - start (start)

# 8. Full details with logs for troubleshooting

MESSAGE_TEMPLATE="Event: {{.Type}}/{{.Action}}\nContainer: {{.Name}} ({{.ShortID}})\nImage: {{.From}}\nStatus: {{.Status}}\nTime: {{.Time}}\n\nRecent Logs:\n{{.GetLogs}}"
MESSAGE_LOG_LINES=50

# Output:

# Event: container/start

# Container: my-app-1 (cd280ca744b1)

# Image: my-app-image

# Status: start

# Time: 2025-10-09T07:37:05Z

#

# Recent Logs:

# [log lines...]

# 9. Alert style for monitoring

MESSAGE_TEMPLATE="⚠️ ALERT: Container {{.Name}} has {{.Action}}\nTimestamp: {{.Time}}\nID: {{.ShortID}}"

# Output:

# ⚠️ ALERT: Container my-app-1 has stop

# Timestamp: 2025-10-09T07:37:05Z

# ID: cd280ca744b1

# 10. JSON-like format

MESSAGE_TEMPLATE="{\"type\": \"{{.Type}}\", \"action\": \"{{.Action}}\", \"name\": \"{{.Name}}\", \"time\": \"{{.Time}}\"}"

# Output: {"type": "container", "action": "start", "name": "my-app-1", "time": "2025-10-09T07:37:05Z"}
