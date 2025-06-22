package authtest

import (
	"auth/internal/config"
	m "auth/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

var cmd *exec.Cmd
var baseUrl string

func TestMain(m *testing.M) {
	if err := serverUp(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer serverDown()

	if err := config.LoadEnvs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	baseUrl = fmt.Sprintf("http://localhost:%d", config.Env.Port)

	// exitCode := m.Run()
	fmt.Println("IT WORKS!")

	// os.Exit(exitCode)
}

func TestRegister(t *testing.T) {
	for _, testCase := range registerTestCases {
		jsonData, err := json.Marshal(testCase.payload)
		if err != nil {
			t.Error(err)
		}

		resp, err := http.Post(
			baseUrl+"/auth/register",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != testCase.statusCode {
			t.Errorf(
				"Unexpected statusCode: %d\nExpected: %d",
				resp.StatusCode,
				testCase.statusCode,
			)
		}

		if resp.StatusCode >= 400 {
			continue
		}

		body := m.User{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			t.Error(err)
		}

		if body.Login != testCase.response.User.Login {
			t.Errorf(
				"Unexpected user login: %s\nExpected: %s",
				body.Login,
				testCase.response.User.Login,
			)
		}
	}
}

func serverUp() error {
	fmt.Println("Starting server...")
	ctx := context.Background()

	cmd = exec.CommandContext(ctx, "bash", "-c", "go run cmd/authmain.go")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	if err := waitForServerReady(baseUrl, 5*time.Second); err != nil {
		return err
	}

	fmt.Println("Server started")
	return nil
}

func waitForServerReady(url string, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		if _, err := http.Get(url); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("Server not ready after %s", timeout)
}

func serverDown() {
	if cmd == nil || cmd.Process == nil {
		return
	}

	fmt.Println("Stopping server...")

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Error sending interrupt: %v\n", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("Server stopped with error: %v\n", err)
		} else {
			fmt.Println("Server stopped gracefully")
		}
	case <-time.After(5 * time.Second):
		fmt.Println("Force killing server after timeout")
		cmd.Process.Kill()
	}
}