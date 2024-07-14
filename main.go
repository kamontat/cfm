package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
)

func main() {
	// Define command-line flags
	nextHop := flag.String("nexthop", "http://localhost:8080", "URL of the next hop (target) server")
	listenAddr := flag.String("listen", ":8000", "Address to listen on")
	logRequests := flag.Bool("log", false, "Enable request logging")
	daemonize := flag.Bool("daemon", false, "Run as a daemon")

	// Parse the flags
	flag.Parse()

	// Handle daemonization
	if *daemonize {
		if !runningAsDaemon() {
			daemonizeProcess()
			return
		}
	}

	// Parse the next hop URL
	target, err := url.Parse(*nextHop)
	if err != nil {
		log.Fatalf("Invalid next hop URL: %v", err)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Optionally add request logging
	if *logRequests {
		proxy.Transport = &loggingRoundTripper{http.DefaultTransport}
	}

	// Create a handler that will be used to serve all requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	// Start the server
	fmt.Printf("Starting reverse proxy server on %s, forwarding to %s\n", *listenAddr, *nextHop)
	log.Fatal(http.ListenAndServe(*listenAddr, handler))
}

// loggingRoundTripper is a custom RoundTripper that logs requests
type loggingRoundTripper struct {
	wrapped http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("Proxying request: %s %s", req.Method, req.URL)
	return l.wrapped.RoundTrip(req)
}

// runningAsDaemon checks if the current process is running as a daemon
func runningAsDaemon() bool {
	return os.Getenv("FORKED") == "1"
}

// daemonizeProcess starts a new process as a daemon
func daemonizeProcess() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Env = append(os.Environ(), "FORKED=1")

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start daemon: %v", err)
	}

	fmt.Println("Daemon started successfully.")
	os.Exit(0)
}
