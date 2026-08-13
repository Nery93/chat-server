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
- Two test clients: a React/TypeScript console (`frontend/`) and a standalone, dependency-free HTML/JS page (`web/test-client.html`)

## Stack

- Go — native `net/http`, no routing framework
- [`gorilla/websocket`](https://github.com/gorilla/websocket) for the WebSocket protocol
- React 19 + TypeScript + Vite for the test client
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
├── frontend/             # React/TypeScript test console (Vite)
└── web/
    └── test-client.html   # standalone, dependency-free test client
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

Start the server first (either method above), then use one of the two clients:

**React console** (`frontend/`):

```bash
cd frontend
npm install
npm run dev
```

Open the printed local URL, enter a name and a room, and connect. Open it again in another tab with a different name, in the same room, to simulate a conversation between two users. The WebSocket address defaults to `ws://localhost:8080` and isn't shown in the UI; override it by copying `frontend/.env.example` to `.env` and setting `VITE_WS_URL`.

**Plain HTML client** (`web/test-client.html`):

Open the file directly in your browser — no install, no server needed to serve it.

### Routes

```
GET /ws/{room}?user={name}   → upgrades to WebSocket, joins "room" identified as "name"
```

If `user` is not provided, the server defaults to `"Anonymous"`.

## Roadmap

This project's history and next steps are documented in [`CLAUDE.md`](./CLAUDE.md).