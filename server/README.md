# Log Generator Server

This service generates dummy logs after x seconds, logs can be divided into multiple levels.

## Structure

- `cmd/server/` - Contains the main application entry point
- `internal/` - Contains internal packages:
  - `config/` - Configuration management
  - `logger/` - Logging functionality
  - `websocket/` - WebSocket handling
  - `loggenerator/` - Code that generates mock logs
  - `model/` - Data models

## Running the Server

To run the server:

```bash
cd cmd/server
go run main.go
```

## Dependencies

This project uses Go modules for dependency management.
