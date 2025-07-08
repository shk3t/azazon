package authtest

import (
	"auth/internal/model"
	"auth/internal/setup"
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"base/pkg/sugar"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var baseUrl string

func TestMain(m *testing.M) {
	workDir := filepath.Dir(sugar.Default(os.Getwd()))
	os.Setenv(setup.ServiceName+"_TEST", "true")
	os.Setenv(setup.ServiceName+"_PORT", "17071")

	err := setup.InitAll("../../.env", workDir)
	if err != nil {
		if log.TLog != nil {
			log.TLog(err)
		}
		setup.GracefullExit(1)
	}

	baseUrl = fmt.Sprintf("http://localhost:%d", setup.Env.Port)
	cmd, err := baseSetup.ServerUp(workDir, baseUrl, log.TLog)
	if err != nil {
		baseSetup.ServerDown(cmd, log.TLog)
		log.TLog(err)
		setup.GracefullExit(1)
	}

	log.TLog("Running tests...")
	exitCode := m.Run()
	log.TLog("Test run finished")
	baseSetup.ServerDown(cmd, log.TLog)
	setup.GracefullExit(exitCode)
}

func TestRegister(t *testing.T) {
	for _, testCase := range registerTestCases {
		jsonData, err := json.Marshal(testCase.payload)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(
			baseUrl+"/auth/register",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != testCase.statusCode {
			t.Fatalf(
				"Unexpected statusCode: %d\nExpected: %d",
				resp.StatusCode,
				testCase.statusCode,
			)
		}

		if resp.StatusCode >= 400 {
			continue
		}

		body := model.AuthResponse{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			t.Fatal(err)
		}

		if body.User.Login != testCase.response.User.Login {
			t.Fatalf(
				"Unexpected user login: %s\nExpected: %s",
				body.User.Login,
				testCase.response.User.Login,
			)
		}
	}
}