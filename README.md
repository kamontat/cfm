# Golang Reverse Proxy Server

This project implements a simple yet flexible reverse proxy server using Go's `httputil.ReverseProxy`. It's designed to forward incoming HTTP requests to a specified target server, making it useful for various scenarios such as load balancing, API gateway implementations, or adding a layer of abstraction between clients and backend services.

## Features

- Configurable next hop (target) server
- Customizable listen address
- Optional request logging
- Command-line flags for easy configuration

## Prerequisites

- Go 1.13 or higher

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/reverse-proxy-golang.git
   ```
2. Navigate to the project directory:
   ```
   cd reverse-proxy-golang
   ```

## Usage

To run the reverse proxy server, use the following command:

```
go run main.go [flags]
```

### Available Flags

- `-nexthop`: URL of the next hop (target) server (default: "http://localhost:8080")
- `-listen`: Address to listen on (default: ":8000")
- `-log`: Enable request logging (default: false)

### Examples

1. Run with default settings:
   ```
   go run main.go
   ```

2. Specify a different target server and listen address:
   ```
   go run main.go -nexthop http://target-server.com -listen :8080
   ```

3. Enable request logging:
   ```
   go run main.go -log
   ```

## How It Works

1. The program parses command-line flags to configure the proxy server.
2. It creates a `httputil.ReverseProxy` instance to handle the proxying.
3. If logging is enabled, it uses a custom `loggingRoundTripper` to log each request.
4. The server runs indefinitely, handling incoming requests and forwarding them to the specified next hop.

## Customization

You can extend this reverse proxy by modifying the `main.go` file. Some possible enhancements include:

- Adding authentication
- Implementing custom load balancing logic
- Adding request/response modification capabilities
- Implementing retry logic for failed requests

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
