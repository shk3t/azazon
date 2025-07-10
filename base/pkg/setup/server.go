package setup

import (
	api "base/api/go"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetClient(url string) (
	client api.AuthServiceClient,
	closeConn func() error,
	err error,
) {
	conn, err := grpc.NewClient(
		url, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client = api.NewAuthServiceClient(conn)
	return client, conn.Close, nil
}

func ServerUp(workDir string, url string, doLog func(...any)) (*exec.Cmd, error) {
	doLog(fmt.Sprintf("Starting server in `%s` dir...", workDir))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", "build/auth", "cmd/authmain.go")
	cmd.Dir = workDir
	if err := cmd.Start(); err != nil {
		return cmd, fmt.Errorf("Failed to start server building: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return cmd, fmt.Errorf("Failed to build server: %w", err)
	}

	cmd = exec.Command("./build/auth")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return cmd, fmt.Errorf("Failed to start server: %w", err)
	}

	if err := WaitForServerReady(url, 5*time.Second); err != nil {
		return cmd, err
	}

	doLog("Server started successfully")
	return cmd, nil
}

func WaitForServerReady(url string, timeout time.Duration) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		time.Sleep(timeout)
		ticker.Stop()
		done <- struct{}{}
	}()

	for {
		select {
		case <-ticker.C:
			if _, _, err := GetClient(url); err == nil {
				return nil
			}
		case <-done:
			return fmt.Errorf("Server not ready after %s", timeout)
		}
	}
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