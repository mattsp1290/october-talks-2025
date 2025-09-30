# Birb-Client - UUID Writer for Birb-Nest

A Go service that continuously writes random UUID4 key-value pairs to the [birb-nest](https://github.com/birbparty/birb-nest) persistent cache service.

## Features

- ✅ Continuous UUID generation and writing
- ✅ Configurable write interval
- ✅ Graceful shutdown on SIGINT/SIGTERM
- ✅ Simple HTTP client (no external dependencies except uuid)
- ✅ Connection pooling via standard library
- ✅ Structured logging with timestamps

## Prerequisites

- Go 1.24.4 or higher
- Running birb-nest instance (see [birb-nest setup](#birb-nest-setup))

## Installation

```bash
cd example/birb-client
go mod download
```

## Usage

### Basic Usage (with defaults)

```bash
cd example/birb-client
go run ./cmd/main.go
```

This will:
- Connect to birb-nest at `http://localhost:8080`
- Write a new UUID key-value pair every 3 seconds
- Log each operation to stdout

### Custom Configuration

Configure the service using environment variables:

```bash
# Custom birb-nest URL and write interval
BIRB_NEST_URL=http://localhost:8080 \
WRITE_INTERVAL=5 \
LOG_LEVEL=debug \
go run ./cmd/main.go
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BIRB_NEST_URL` | Base URL of the birb-nest API | `http://localhost:8080` |
| `WRITE_INTERVAL` | Seconds between writes | `3` |
| `LOG_LEVEL` | Logging verbosity (debug, info, warn, error) | `info` |

## Birb-Nest Setup

If you don't have birb-nest running, start it with Docker:

```bash
cd /Users/punk1290/git/birb-nest
docker-compose up -d
```

Or use the development mode:

```bash
cd /Users/punk1290/git/birb-nest
make dev
```

Verify birb-nest is running:

```bash
curl http://localhost:8080/health
```

## Example Output

```
2025/09/29 21:00:00 Birb-Client starting...
2025/09/29 21:00:00 Configuration:
2025/09/29 21:00:00   - Birb Nest URL: http://localhost:8080
2025/09/29 21:00:00   - Write Interval: 3s
2025/09/29 21:00:00   - Log Level: info
2025/09/29 21:00:00 Successfully connected to birb-nest
2025/09/29 21:00:00 Starting UUID writer with 3s interval
2025/09/29 21:00:00 Writing: key=a3bb189e-8bf9-4fa2-b23d-0d7f8e8d6c7a value=f7e5d3c1-9a2b-4e1f-8c7d-6b5a4e3d2c1b
2025/09/29 21:00:00 Successfully wrote key=a3bb189e-8bf9-4fa2-b23d-0d7f8e8d6c7a
2025/09/29 21:00:03 Writing: key=b4cc298f-9cg0-5gb3-c34e-1e8g9f9e7d8b value=g8f6e4d2-0b3c-5f2g-9d8e-7c6b5f4e3d2c
2025/09/29 21:00:03 Successfully wrote key=b4cc298f-9cg0-5gb3-c34e-1e8g9f9e7d8b
...
```

## Graceful Shutdown

Press `Ctrl+C` or send `SIGTERM` to gracefully shutdown:

```
^C2025/09/29 21:00:15 Received signal: interrupt
2025/09/29 21:00:15 Initiating graceful shutdown...
2025/09/29 21:00:15 Writer stopped by context cancellation
2025/09/29 21:00:15 Closing writer
2025/09/29 21:00:15 Birb-Client shutdown complete
```

## Architecture

```
┌─────────────────┐
│  birb-client    │
│                 │
│  ┌───────────┐  │      HTTP POST /v1/cache/{key}
│  │  Writer   │──┼──────────────────────────────────┐
│  └───────────┘  │                                  │
│       │         │                                  ▼
│       │         │                            ┌──────────┐
│  ┌────▼─────┐   │                            │  Birb    │
│  │ UUID Gen │   │                            │  Nest    │
│  └──────────┘   │                            │  API     │
└─────────────────┘                            └──────────┘
```

## Project Structure

```
birb-client/
├── cmd/
│   └── main.go              # Entry point with signal handling
├── internal/
│   ├── config/
│   │   └── config.go        # Environment configuration
│   └── writer/
│       └── writer.go        # UUID writing logic
├── go.mod                   # Module definition
├── go.sum                   # Dependency checksums
└── README.md                # This file
```

## Development

### Build

```bash
go build -o birb-client ./cmd/main.go
./birb-client
```

### Test Connectivity

```bash
# Check if birb-nest is accessible
curl http://localhost:8080/health

# Verify written data
curl http://localhost:8080/v1/cache/YOUR_UUID_KEY
```

## Troubleshooting

### Cannot connect to birb-nest

**Error**: `Failed to ping birb-nest: connection refused`

**Solution**: Ensure birb-nest is running:
```bash
curl http://localhost:8080/health
```

### Write failures

**Error**: `Warning: write failed: context deadline exceeded`

**Solution**: Check birb-nest logs for errors:
```bash
cd /Users/punk1290/git/birb-nest
docker-compose logs -f
```

## License

Same as parent project.