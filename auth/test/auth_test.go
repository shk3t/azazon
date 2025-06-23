package authtest

import (
	"auth/internal/config"
	m "auth/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

var cmd *exec.Cmd
var baseUrl string

func TestMain(m *testing.M) {
	initLogger()
	defer deiniLogger()

	os.Setenv(config.ServiceName+"_TEST", "true")
	if err := config.LoadEnvs("../../.env"); err != nil {
		dLog(err)
		os.Exit(1)
	}
	baseUrl = fmt.Sprintf("http://localhost:%d", config.Env.Port)

	if err := serverUp(); err != nil {
		dLog(err)
		os.Exit(1)
	}
	defer serverDown()

	dLog("Running tests...")
	exitCode := m.Run()
	dLog("Tests running finished with exit code:", exitCode)
	// os.Exit(exitCode)
	dLog("Exit")
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
	dLog("Starting server...")
	ctx := context.Background()

	cmd = exec.CommandContext(ctx, "go", "run", "cmd/authmain.go")
	cmd.Dir = "/home/shket/projects/go/azazon/auth"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	if err := waitForServerReady(baseUrl, 5*time.Second); err != nil {
		return err
	}

	dLog("Server started")
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
		dLog("Error: Process is nil")
		return
	}

	dLog("Stopping server...")

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Error sending interrupt: %v\n", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			dLog(fmt.Sprintf("Server stopped with error: %v\n", err))
		} else {
			dLog("Server stopped gracefully")
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		dLog("Server force killed after timeout")
	}
}

var DebugLogger *log.Logger
var dLog func(...any)

func initLogger() {
	err := os.MkdirAll("../logs", 0755)
	if err != nil {
		panic("Can't create \"logs\" directory")
	}
	debugLogFile, err := os.OpenFile("../logs/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Can't open \"debug.log\" file")
	}

	DebugLogger = log.New(debugLogFile, "", log.LstdFlags|log.Lshortfile)
	dLog = DebugLogger.Println
}

func deiniLogger() {
	writer := DebugLogger.Writer()
	writeCloser, ok := writer.(io.WriteCloser)
	if ok {
		err := writeCloser.Close()
		if err != nil {
			panic("Can't close log file")
		}
	}
}