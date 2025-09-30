<div align="center">
  <a href="https://github.com/filippofinke/docker-events">
    <img width="200px" src="https://github.com/user-attachments/assets/e9712d6a-32e9-4e9b-a545-11aa88dadf65" alt="Docker Events" />
  </a>
  <h3 align="center">Docker Events</h3>
</div>

> Watch Docker events in real-time and dispatch rich notifications with a lightweight Go service.

This project uses docker system events and forwards meaningful summaries through configurable notification channels powered by [`nikoksr/notify`](https://github.com/nikoksr/notify). It is designed to be small, dependable, and easy to extend with new transports or event processing rules.

## Features

- [x] Real-time streaming of `docker system events` via a managed watcher
- [x] Human-friendly notifications enriched with context (timestamp, actor attributes, status)
- [x] Multi-channel delivery (Slack, Telegram, Discord) powered by `github.com/nikoksr/notify`
- [x] Opt-in filtering by Docker event types and CLI filters
- [x] Environment variable driven configuration with validation at startup
- [x] Composable code structure (config, watcher, notifier) and unit tests for formatting helpers

## Quick Start

Prerequisites

- Go 1.24+
- Docker CLI with access to the target daemon
- At least one notifications provider configured (Slack bot token + channel IDs, Telegram bot token + chat IDs, or Discord bot token + channel IDs)

Clone & install dependencies

```bash
git clone https://github.com/filippofinke/docker-events.git
cd docker-events
go mod tidy
```

Configure environment

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Update `.env` with your Slack token, target channels, and any optional Docker filters.

Run locally

```bash
# start the watcher (module-aware path)
go run ./cmd
```

The service will stream logs to stdout and post notifications for matching Docker events. Stop with `Ctrl+C`.

## Configuration

All settings are provided via environment variables (see `.env.example`). Key options:

- `SLACK_BOT_TOKEN`: Slack bot token used to authenticate the notifier.
- `SLACK_CHANNEL_IDS`: Comma-separated list of Slack channel IDs (e.g. `C0123456,C0ABCDEF`).
- `TELEGRAM_BOT_TOKEN`: Telegram bot token created with [BotFather](https://core.telegram.org/bots#6-botfather).
- `TELEGRAM_CHAT_IDS`: Comma-separated list of chat IDs (negative values for group chats are supported).
- `DISCORD_BOT_TOKEN`: Discord bot token generated from the Developer Portal.
- `DISCORD_CHANNEL_IDS`: Comma-separated list of Discord channel IDs to notify.
- `NOTIFY_SUBJECT_PREFIX`: Prefix for notification subjects (defaults to `Docker event`).
- `DOCKER_CLI_PATH`: Path to the Docker CLI binary (defaults to `docker`).
- `DOCKER_EVENT_FILTERS`: Comma-separated filters passed to `docker system events` (same syntax as the CLI `--filter` flag, e.g. `status=start,type=container`).
- `DOCKER_EVENT_TYPES`: Comma-separated list of Docker event types to keep (e.g. `container,image,volume`).

> **Security note:** Do not commit an `.env` file containing real tokens. Use `.env` locally or provide the variables through your orchestrator of choice.

## Docker Event Filters

The Docker CLI supports a rich set of filters that can be combined in `DOCKER_EVENT_FILTERS`. Supported filter keys include:

- `config=<name or id>`
- `container=<name or id>`
- `daemon=<name or id>`
- `event=<event action>`
- `image=<repository or tag>`
- `label=<key>` or `label=<key>=<value>`
- `network=<name or id>`
- `node=<id>`
- `plugin=<name or id>`
- `scope=<local or swarm>`
- `secret=<name or id>`
- `service=<name or id>`
- `type=<container|image|volume|network|daemon|plugin|service|node|secret|config>`
- `volume=<name>`

Provide multiple filters by comma-separating entries (e.g. `DOCKER_EVENT_FILTERS=event=start,scope=swarm`); the service will translate each entry into an individual `--filter` flag for `docker system events`.

More details in the [Docker documentation](https://docs.docker.com/reference/cli/docker/system/events/#filter).

## Docker Event Types

`DOCKER_EVENT_TYPES` narrows processing to one or more top-level Docker object kinds. Valid values:

- `container`
- `image`
- `plugin`
- `volume`
- `network`
- `daemon`
- `service`
- `node`
- `secret`
- `config`

Leave the variable empty to accept every event type from the stream.

More details in the [Docker documentation](https://docs.docker.com/engine/reference/commandline/system_events/#object-types).

## Extending Notifications

`internal/notifier` wraps `github.com/nikoksr/notify`, so adding more destinations is straightforward:

1. Import the desired service package (e.g. `github.com/nikoksr/notify/service/telegram`).
2. Create a service instance in `Setup` based on new configuration.
3. Register it with the shared notifier (`n.client.UseServices(service)`).

## Running Tests

```bash
go test ./...
```

## Docker Usage

A minimal container image can be built with:

```bash
# build the Go binary locally and package it into a container image
docker build -t docker-events:latest .
```

The included `Dockerfile` uses a multi-stage build: it compiles a static Go binary in a Go builder image and copies it into the official Docker CLI image so the binary can call `docker` if needed.

Important runtime considerations:

- The service talks to the Docker daemon. In most deployments you should mount the host Docker socket into the container so the service can observe events:

  - `/var/run/docker.sock:/var/run/docker.sock:ro` (read-only mount used in the example compose file)

- Environment variables are used for configuration. The repository contains a `.env.example`—copy it to `.env` and set your provider tokens and channels. Do not commit `.env` with real secrets.

Compose example (loads `.env` automatically)

```yaml
services:
  docker-events:
    build: .
    env_file:
      - .env
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    restart: unless-stopped
```

Start with docker-compose:

```bash
# ensure .env exists in the project root (copy from .env.example)
cp .env.example .env
# build & start in background
docker compose up -d --build
# view logs
docker compose logs -f docker-events
```

If you prefer to run the image directly:

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock:ro \
  --env-file .env filippofinke/docker-events:latest
```

Remember: by default Docker Compose loads a top-level `.env` file. The `env_file` entry above is explicit and can be used by other tools that also support `env_file`.

## Author

👤 **Filippo Finke**

- Website: [https://filippofinke.ch](https://filippofinke.ch)
- Twitter: [@filippofinke](https://twitter.com/filippofinke)
- GitHub: [@filippofinke](https://github.com/filippofinke)
- LinkedIn: [@filippofinke](https://linkedin.com/in/filippofinke)
