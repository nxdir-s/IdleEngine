[![Go Report Card](https://goreportcard.com/badge/github.com/nxdir-s/IdleEngine)](https://goreportcard.com/report/github.com/nxdir-s/IdleEngine)

# IdleEngine

IdleEngine is a game engine written in Go for idle/incremental type games. Players connect to the game server via websockets and have their actions processed at set intervals

## Getting Started

### Prerequisites

- Go (version 1.24 or higher recommended)
- Docker/Docker Compose

### Running IdleEngine

1. **Build Docker Images**:

   ```bash
   docker compose build
   ```

2. **Run Infrastructure**:

   ```bash
   docker compose up
   ```

### Connect to Game Server

Once the game server is up and running you can connect to it using the following

```bash
go run cmd/client/main.go ws://127.0.0.1:80/ws
```

The client defaults to 50 concurrent connections but can be configured by adding an additional argument with your desired count

```bash
go run cmd/client/main.go ws://127.0.0.1:80/ws 100
```

### Kafka UI

A Kafka UI for managing clusters can be accessed at the following address `http://127.0.0.1:80`

> #### [Project Structure Docs](https://github.com/nxdir-s/go-hexarch)
