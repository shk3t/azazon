package authtest

import (
	"auth/internal/config"
	m "auth/internal/model"
	"auth/pkg/log"
	"auth/pkg/setup"
	"auth/pkg/sugar"
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
	os.Setenv(config.ServiceName+"_TEST", "true")
	os.Setenv(config.ServiceName+"_PORT", "17071")

	log.Init(workDir)
	defer log.Deinit()

	if err := config.LoadEnvs("../../.env"); err != nil {
		log.TLog(err)
		os.Exit(1)
	}
	baseUrl = fmt.Sprintf("http://localhost:%d", config.Env.Port)

	cmd, err := setup.ServerUp(workDir, baseUrl, log.TLog)
	if err != nil {
		log.TLog(err)
		os.Exit(1)
	}
	defer setup.ServerDown(cmd, log.TLog)

	log.TLog("Running tests...")
	exitCode := m.Run()
	log.TLog("Tests running finished with exit code:", exitCode)
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

		body := m.AuthResponse{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			t.Error(err)
		}

		if body.User.Login != testCase.response.User.Login {
			t.Errorf(
				"Unexpected user login: %s\nExpected: %s",
				body.User.Login,
				testCase.response.User.Login,
			)
		}
	}
}