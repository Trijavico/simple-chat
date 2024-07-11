# Simple Real-Time Chat

A simple real-time chat application built with Go and HTMX.

## How It Works

- The server uses a hub to manage WebSocket connections
- Clients connect to the server via WebSocket
- Messages are broadcast to all connected clients
- HTMX handles real-time updates on the frontend

## Getting Started

1. Clone the repository
2. Ensure Go is installed on your system
3. Run the server: `make run` or `go run main.go`
4. Open a web browser and navigate to `http://localhost:3000`
