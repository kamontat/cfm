//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

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
