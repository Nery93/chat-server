# chat-server

A real-time, multi-room chat application built with Go and WebSockets. Users can join a room, exchange messages instantly, see who is online, and ask an NVIDIA Nemotron AI assistant for help with the `/ai` command.

The project was built to practice Go concurrency patterns in a real application: goroutines, channels, mutexes, contexts, and graceful shutdown.

## Features

- Multiple rooms created dynamically from the WebSocket URL
- Real-time WebSocket messaging with no polling or repeated HTTP requests
- Structured JSON messages for chat, join, leave, typing, and online-user updates
- Server-enforced user identity and message types, preventing clients from impersonating users or sending system messages
- Online user list for each room
- Typing indicator
- Join and leave announcements
- Ping/pong keepalive with read and write deadlines
- Message size limit of 4 KB
- AI assistant powered by NVIDIA Nemotron, triggered only with `/ai <question>`
- Token usage logging in the server terminal for every AI response
- Graceful shutdown on `Ctrl+C` or `SIGTERM`
- Automated tests for room, client, and handler behavior

## Tech Stack

- Go 1.25+ with native `net/http`
- [`gorilla/websocket`](https://github.com/gorilla/websocket)
- React 19, TypeScript, and Vite
- [`react-markdown`](https://github.com/remarkjs/react-markdown) with GitHub Flavored Markdown support
- NVIDIA NIM API using `nvidia/nemotron-3-ultra-550b-a55b`
- Docker and Docker Compose

## Architecture

```text
Browser
    | WebSocket
    v
Go HTTP server
    |-- handler: upgrades the connection and assigns a room
    |-- client: reads and writes WebSocket messages
    |-- room: manages connected clients and broadcasts messages
    |-- ai: calls NVIDIA Nemotron when a user sends /ai
    v
React chat console
```

Each connected client has its own goroutines for reading and writing. Rooms protect their shared client list with a mutex, while channels carry broadcast messages to connected clients.

## Project Structure

```text
.
├── main.go                 # Server startup and graceful shutdown
├── internal/
│   ├── ai/                 # NVIDIA Nemotron integration
│   ├── client/             # Connected WebSocket client read/write loops
│   ├── handler/            # HTTP routes and WebSocket upgrade
│   ├── message/            # WebSocket JSON message structure
│   └── room/               # Room membership and broadcast logic
├── frontend/               # React/TypeScript chat console
└── web/test-client.html    # Standalone browser test client
```

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js and npm, for the React client
- An NVIDIA API key, only when using the `/ai` command

### Run the server locally

```bash
go run .
```

The server listens on port `8080` by default. Set `PORT` to use a different port:

```bash
PORT=9000 go run .
```

### Configure AI

The AI integration reads its key from `NVIDIA_API_KEY`. Export it in the same terminal before starting the server:

```bash
export NVIDIA_API_KEY="your-key-here"
go run .
```

Never commit API keys. The root `.env` file is ignored by Git, but Go does not load it automatically.

### Run with Docker

```bash
docker compose up --build
```

The API is available at `http://localhost:8080`. To use AI in Docker, pass `NVIDIA_API_KEY` into the `chat-server` service environment in `docker-compose.yml`.

### Run the React client

In a second terminal:

```bash
cd frontend
npm install
npm run dev
```

Open the local URL printed by Vite. Connect from a second browser tab with another username to test a real-time conversation.

The client connects to `ws://localhost:8080` by default. Set `VITE_WS_URL` in `frontend/.env` to point it at another server.

### Use the standalone client

Open `web/test-client.html` directly in a browser. It does not require npm or a frontend development server.

## Routes

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/ws/{room}?user={name}` | Upgrades the connection to WebSocket and joins the selected room. The default user is `Anonymous`. |
| `GET` | `/rooms` | Returns active rooms and their current client counts as JSON. |

## AI Assistant

Send a message beginning with `/ai ` followed by a question:

```text
/ai Explain how WebSockets work
```

The command itself is not broadcast to the room. The server sends the assistant response as a normal chat message from `AI` and logs the token usage in the server terminal:

```text
Uso NVIDIA: prompt=22 resposta=800 total=822
```

The current response limit is `1500` completion tokens. It can be adjusted in `internal/ai/ai.go` through `MaxTokens`.

## Testing

Run all Go tests:

```bash
go test ./...
```

Validate the frontend:

```bash
cd frontend
npm run build
npm run lint
```

## Current Limitations

- AI responses use a non-streaming request, so the chat receives the response only after NVIDIA finishes generating it.
- AI calls currently run synchronously inside the client's read loop.
- Messages and rooms live only in memory and disappear when the server restarts.
- No rate limiting is applied to `/ai` requests yet.

## Next Steps

- Stream AI responses token by token
- Run AI requests independently so a client can keep sending messages while the model responds
- Add timeouts and user-visible AI error messages
- Add rate limiting for AI requests
- Persist chat history

## License

This repository is intended as a learning project.