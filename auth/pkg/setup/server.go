package setup

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func ServerUp(workDir string, url string, doLog func(...any)) (*exec.Cmd, error) {
	doLog(fmt.Sprintf("Starting server in `%s` dir...", workDir))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", "build/auth", "cmd/authmain.go")
	cmd.Dir = workDir
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Failed to start server building: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Failed to build server: %w", err)
	}

	cmd = exec.Command("./build/auth")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Failed to start server: %w", err)
	}

	if err := WaitForServerReady(url, 5*time.Second); err != nil {
		return nil, err
	}

	doLog("Server started successfully")
	return cmd, nil
}

func WaitForServerReady(url string, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		if _, err := http.Get(url); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("Server not ready after %s", timeout)
}

func ServerDown(cmd *exec.Cmd, doLog func(...any)) {
	if cmd == nil || cmd.Process == nil {
		doLog("Error: Process is nil")
		return
	}

	doLog("Stopping server...")
	doLog(cmd.Args)

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Error sending interrupt: %v\n", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			doLog(fmt.Sprintf("Server stopped: %v\n", err))
		} else {
			doLog("Server stopped gracefully")
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		doLog("Server force killed after timeout")
	}
}