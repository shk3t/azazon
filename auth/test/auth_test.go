package authtest

import (
	"auth/internal/config"
	"auth/internal/service"
	"auth/internal/setup"
	"base/api/auth"
	errpkg "base/pkg/error"
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
	os.Setenv(config.AppName+"_TEST", "true")
	os.Setenv(config.AppName+"_PORT", "17071")

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

		out, err := client.Register(ctx, &auth.RegisterRequest{
			Login:    testCase.request.Login,
			Password: testCase.request.Password,
		})
		st, ok := status.FromError(err)
		requireOk(t, ok, err)

		assertEqual(t, st.Code(), testCase.statusCode, "status code")
		if st.Code() != codes.OK {
			continue
		}

		claims, err := service.ParseJwtToken(out.Token)
		requireNoError(t, err)

		assertEqual(t, claims.Login, testCase.response.Login, "user login")
	}
}

func TestLogin(t *testing.T) {
	client, closeConn, _ := baseSetup.GetClient(grpcUrl)
	defer closeConn()

	for _, testCase := range loginTestCases {
		ctx := context.Background()

		// It is supposed to work
		_, _ = client.Register(ctx, &auth.RegisterRequest{
			Login:    testCase.request.Login,
			Password: testCase.request.Password,
		})

		_, err := client.Login(ctx, &auth.LoginRequest{
			Login:    testCase.request.Login,
			Password: testCase.request.Password,
		})
		st, ok := status.FromError(err)
		requireOk(t, ok, err)

		assertEqual(t, st.Code(), testCase.statusCode, "status code")
	}
}

func TestValidateToken(t *testing.T) {
	client, closeConn, _ := baseSetup.GetClient(grpcUrl)
	defer closeConn()

	for _, testCase := range validateTokenTestCases {
		ctx := context.Background()

		outR, _ := client.Register(ctx, &auth.RegisterRequest{
			Login:    testCase.request.Login,
			Password: testCase.request.Password,
		})
		requireNotNil(t, outR, "response")

		outV, err := client.ValidateToken(ctx, &auth.ValidateTokenRequest{
			Token: outR.Token,
		})
		st, ok := status.FromError(err)
		requireOk(t, ok, err)
		requireOk(t, outV.Valid, errpkg.InvalidToken)

		assertEqual(t, st.Code(), testCase.statusCode, "status code")
	}
}