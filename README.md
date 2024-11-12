# Cross-Platform Mac/WindowsCommand Execution Application

## Overview
This application is a cross-platform command execution tool capable of executing simple system commands (ping and system info). It listens for HTTP requests and returns JSON responses.

## Installation
For macOS, use the `.pkg` installer. Ensure the app starts on boot.

## Build
Run `go build ./cmd/main.go` to compile the application.

## API Documentation
- **POST /execute**
  - **Type**: Command type ("ping" or "sysinfo")
  - **Payload**: Host for ping command

## Usage
1. Start the server: `./crossplat-app`
2. Send POST requests to `localhost:8080/execute`.

## Testing
Run `go test ./...` to execute tests.