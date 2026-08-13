# chat-server

A real-time chat server written in Go, using plain WebSocket (no HTTP framework). Users join a room through a URL, and every message is broadcast in real time to everyone in that room — no polling, no repeated HTTP requests.

This project was built as a learning exercise in Go concurrency (goroutines, channels, mutexes, context, graceful shutdown), applied to a real use case: a WebSocket server handling multiple rooms and multiple concurrent clients.

## Features

- Multiple rooms, created dynamically from the name used in the URL
- Broadcast of messages to every client in a room, protected by a mutex
- Structured JSON messages (`chat`, `join`, `leave`), with the type and sender always validated/enforced server-side — a client can't impersonate another user or fake a system message
- User identification via query string (`?user=name`)
- Automatic announcements when someone joins or leaves a room
- Keepalive via ping/pong and read/write deadlines, to detect dead connections
- Graceful shutdown (`Ctrl+C` or `SIGTERM` let in-flight requests finish, with a timeout)
- Standalone HTML/JS test client (`web/test-client.html`), no dependencies

## Stack

- Go — native `net/http`, no routing framework
- [`gorilla/websocket`](https://github.com/gorilla/websocket) for the WebSocket protocol
- Docker + Docker Compose
- No database — in-memory state

## Project structure

```
/
├── main.go              # server startup, graceful shutdown
├── internal/
│   ├── client/           # a connected WebSocket client (read/write)
│   ├── room/              # a room's client list, broadcast
│   ├── handler/          # HTTP routes, WebSocket upgrade, room/client wiring
│   └── message/          # message format (JSON)
└── web/
    └── test-client.html   # standalone test client
```

## Getting started

### With Go, locally

Requires Go 1.25+.

```bash
go run .
```

Runs on port `8080` by default. To use a different port:

```bash
PORT=9000 go run .
```

### With Docker

```bash
docker compose up --build
```

Available at `http://localhost:8080` (configurable in `docker-compose.yml`).

## Trying it out

1. Start the server (either method above).
2. Open `web/test-client.html` directly in your browser — no need to serve it from anywhere.
3. Enter a name and a room, and connect.
4. Open the same file in another tab, with a different name, in the same room, to simulate a conversation between two users.

### Routes

```
GET /ws/{room}?user={name}   → upgrades to WebSocket, joins "room" identified as "name"
```

If `user` is not provided, the server defaults to `"Anonymous"`.

## Roadmap

This project's history and next steps are documented in [`CLAUDE.md`](./CLAUDE.md).