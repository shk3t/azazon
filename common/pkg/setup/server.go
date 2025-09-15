package setup

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ServerUp(appName string, workDir string, url string, logger *log.Logger) (*exec.Cmd, error) {
	appName = strings.ToLower(appName)
	port := strings.Split(url, ":")[1]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stdoutBuf bytes.Buffer
	cmd := exec.CommandContext(
		ctx, "bash", "-c",
		"ss -tuln | tr -s ' ' | cut -d ' ' -f 5 | grep "+port,
	)
	cmd.Stdout = &stdoutBuf
	cmd.Run()
	if stdoutBuf.String() != "" {
		logger.Printf("%s port is already occupied. Killing previously run server...", port)
		err := exec.CommandContext(
			ctx, "bash", "-c",
			fmt.Sprintf("kill $(lsof -i :%s -t)", port),
		).Run()
		if err != nil {
			return nil, fmt.Errorf("Unable to kill process: %w", err)
		}
	}

	cmd = exec.CommandContext(
		ctx,
		"go", "build", "-o",
		"build/"+appName,
		"cmd/main.go",
	)
	cmd.Dir = workDir
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Failed to start server building: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Failed to build server: %w", err)
	}

	logger.Printf("Starting server in `%s` dir...\n", workDir)

	cmd = exec.Command("./build/" + appName)
	cmd.Dir = workDir
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()

	if err := cmd.Start(); err != nil {
		return cmd, fmt.Errorf("Failed to start server: %w", err)
	}

	logger.Println("Server started successfully")
	return cmd, nil
}

func ServerDown(cmd *exec.Cmd, logger *log.Logger) {
	if cmd == nil || cmd.Process == nil {
		logger.Println("Error: Process is nil")
		return
	}

	logger.Println("Stopping server...")

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Error sending interrupt: %v\n", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			logger.Printf("Server stopped: %v\n", err)
		} else {
			logger.Println("Server stopped gracefully")
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		logger.Println("Server force killed after timeout")
	}
}