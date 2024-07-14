// +build !windows

package main

import (
	"os/exec"
	"syscall"
)

func daemonizeProcess() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Env = append(os.Environ(), "DAEMON=1")

	// For Unix-like systems, detach the process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start daemon: %v", err)
	}

	fmt.Println("Daemon started successfully.")
	os.Exit(0)
}
