# Golang Reverse Proxy Server

This project implements a simple yet flexible reverse proxy server using Go's `httputil.ReverseProxy`. It's designed to forward incoming HTTP requests to a specified target server, making it useful for various scenarios such as load balancing, API gateway implementations, or adding a layer of abstraction between clients and backend services.

## Features

- Configurable next hop (target) server
- Customizable listen address
- Optional request logging
- Command-line flags for easy configuration
- Cross-platform daemonization support (Windows, Linux, macOS)

## Prerequisites

- Go 1.13 or higher

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/presbrey/shrp.git
   ```
2. Navigate to the project directory:
   ```
   cd shrp
   ```

## Usage

To run the reverse proxy server, use the following command:

```
go run main.go [flags]
```

### Use Cases

1. **Testing DNS Migrations**: When migrating services between different DNS providers or IP addresses, this reverse proxy can be used to verify the new configuration before making the final switch. By setting up the proxy to forward requests to the new IP or DNS, you can test the new setup without changing the actual DNS records.

   Example:
   ```
   go run main.go -nexthop https://new-ip-or-dns.com -host original-domain.com
   ```
   This allows you to send requests to the proxy (which will appear to come from `original-domain.com`) and have them forwarded to the new IP or DNS, simulating the post-migration behavior.

2. **Load Balancing**: Although not a full-featured load balancer, this proxy can be used for simple round-robin load balancing by running multiple instances with different `-nexthop` targets.

3. **API Gateway**: Act as a simple API gateway, potentially adding authentication or request modification before forwarding to backend services.

4. **SSL Termination**: By running the proxy with SSL and forwarding to non-SSL backends, you can offload SSL processing.

### Available Flags

- `-host`: Value of the Host header to send to the next hop server
- `-nexthop`: URL of the next hop (target) server (default: "https://httpbin.org/")
- `-listen`: Address to listen on (default: ":8000")
- `-log`: Enable request logging (default: false)
- `-daemon`: Run as a daemon (default: false)
- `-insecure`: Ignore SSL certificate errors (default: false)

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

4. Run as a daemon:
   ```
   go run main.go -daemon
   ```

## Daemonization

The `-daemon` flag allows you to run the reverse proxy as a background process. This is useful for long-running services that need to persist even after the user logs out. The daemonization process works on Windows, Linux, and macOS.

When run as a daemon:
1. The program will start a new process in the background.
2. The original process will exit immediately.
3. The new process will continue running in the background, serving requests.

On Unix-like systems (Linux, macOS), the daemon process will be detached from the terminal session.

The daemon process sets the environment variable `DAEMON=1` to indicate that it's running in daemon mode.

To stop the daemon:
- On Windows: Use the Task Manager or command line tools to terminate the process.
- On Unix-like systems: Use the `ps` command to find the process ID, then use `kill` to terminate it.

## How It Works

1. The program parses command-line flags to configure the proxy server.
2. If the `-daemon` flag is set, it starts a new background process.
3. It creates a `httputil.ReverseProxy` instance to handle the proxying.
4. If logging is enabled, it uses a custom `loggingRoundTripper` to log each request.
5. The server runs indefinitely, handling incoming requests and forwarding them to the specified next hop.

## Customization

You can extend this reverse proxy by modifying the `main.go` file. Some possible enhancements include:

- Adding authentication
- Implementing custom load balancing logic
- Adding request/response modification capabilities
- Implementing retry logic for failed requests
- Adding support for HTTPS

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
