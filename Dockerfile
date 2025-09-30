FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/docker-events ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/docker-events /bin/docker-events
ENTRYPOINT ["/bin/docker-events"]
