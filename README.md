# Cross-Platform Mac/WindowsCommand Execution Application

## Overview
This application is a cross-platform command execution tool capable of executing simple system commands (ping and system info). It listens for HTTP requests and returns JSON responses.

## Development

This application uses the make utility for building, testing, linting and packaging the application.

To install make on macOS, run `brew install make`.

To install make on Windows either:
 - Using Chocolatey, run `choco install make` 
 - Directly download and install binary from [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm)

To lint the application, run
```
make lint
```

To test the application, run
```
make test
```

To build the application, run 
```
make build
```

To package the application, run
```
make package
```

To package the application with a custom name run
```
make package APP_NAME=<name>
```

To package the application for a specific platform run
```
make package PLATFORM=<platform>
```
Where `<platform>` is one of `darwin`, `windows` or `linux` (linux is not yet supported).

## Installation

To install the application on macOS, double click the dmg file created in the `dist` directory by the package target.

Then copy the application to the Applications directory, and follow the instructions in the README.txt to install the launch agent


## API Documentation

This application exposes a simple HTTP API with a single endpoint accepting POST requests for executing commands.

- **POST /execute**
  - **Type**: Command type ("ping" or "sysinfo")
  - **Payload**: Host for ping command

All requests should be sent as JSON with the following structure:

```
{
    "type": string,
    "payload": string
}
```

An example request to execute the `ping` command with the payload `127.0.0.1` would be:

```
{
    "type": "ping",
    "payload": "127.0.0.1"
}
```

## Running the application

If the application has been installed using the dmg and following the instructions in the README.txt file, the application will be available in the Applications directory and will be launched automatically at login. 

If the application has been built from source, you can run the application manually using the following MAKE target. 

The application uses an environment variable, COMMANDER_APP_PORT, to set the port the application listens on. If no environment variable is set, the application will default to port 8080.

1. Set the environment variable, COMMANDER_APP_PORT, to the desired port.
1. Start the server: `make run`
2. Send POST requests to `localhost:<port>/execute`.

## Testing
Run `make test` to execute tests.
