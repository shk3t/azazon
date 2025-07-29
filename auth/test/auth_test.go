package authtest

import (
	"auth/internal/config"
	"auth/internal/setup"
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"base/pkg/sugar"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcUrl string

func TestMain(m *testing.M) {
	workDir := filepath.Dir(sugar.Default(os.Getwd()))
	os.Setenv(config.ServiceName+"_TEST", "true")
	os.Setenv(config.ServiceName+"_PORT", "17071")

	err := setup.InitAll("../../.env", workDir)
	if err != nil {
		if log.TLog != nil {
			log.TLog(err)
		}
		setup.GracefullExit(1)
	}

	grpcUrl = fmt.Sprintf("localhost:%d", config.Env.Port)
	cmd, err := baseSetup.ServerUp(workDir, grpcUrl, log.TLog)
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
	client, closeConn, _ := baseSetup.GetClient(grpcUrl)
	defer closeConn()

	for _, testCase := range registerTestCases {
		ctx := context.Background()

		out, err := client.Register(ctx, testCase.payload)
		st, ok := status.FromError(err)
		if !ok {
			t.Fatal(err)
		}

		if st.Code() != testCase.statusCode {
			t.Fatalf(
				"Unexpected statusCode: %d\nExpected: %d",
				st.Code(),
				testCase.statusCode,
			)
		}

		if st.Code() != codes.OK {
			continue
		}

		if out.User.Login != testCase.response.User.Login {
			t.Fatalf(
				"Unexpected user login: %s\nExpected: %s",
				out.User.Login,
				testCase.response.User.Login,
			)
		}
	}
}
