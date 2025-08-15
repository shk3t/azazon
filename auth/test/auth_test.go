package authtest

import (
	"auth/internal/config"
	"auth/internal/setup"
	"base/api/auth"
	errpkg "base/pkg/errors"
	"base/pkg/log"
	"base/pkg/model"
	"base/pkg/service"
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

	err := setup.InitAll(workDir)
	if err != nil {
		setup.DeinitAll()
		panic(err)
	}
	logger := log.Loggers.Test

	grpcUrl = fmt.Sprintf("localhost:%d", config.Env.TestPort)
	cmd, err := baseSetup.ServerUp(workDir, grpcUrl, logger)
	if err != nil {
		baseSetup.ServerDown(cmd, logger)
		logger.Println(err)
		setup.GracefullExit(1)
	}

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")
	baseSetup.ServerDown(cmd, logger)
	setup.GracefullExit(exitCode)
}

func TestRegister(t *testing.T) {
	client, closeConn, _ := baseSetup.GetClient(grpcUrl)
	defer closeConn()

	for _, testCase := range registerTestCases {
		ctx := context.Background()

		resp, err := client.Register(ctx, &auth.RegisterRequest{
			Login:    testCase.request.Login,
			Password: testCase.request.Password,
		})
		st, ok := status.FromError(err)
		requireOk(t, ok, err)

		assertEqual(t, st.Code(), testCase.statusCode, "status code")
		if st.Code() != codes.OK {
			continue
		}

		claims, err := service.ParseJwtToken(resp.Token)
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

		respReg, _ := client.Register(ctx, &auth.RegisterRequest{
			Login:    testCase.registerRequest.Login,
			Password: testCase.registerRequest.Password,
		})
		requireNotNil(t, respReg, "response")

		resp, err := client.ValidateToken(ctx, &auth.ValidateTokenRequest{
			Token: respReg.Token,
		})
		st, ok := status.FromError(err)
		requireOk(t, ok, err)
		requireOk(t, resp.Valid, errpkg.InvalidToken)

		assertEqual(t, st.Code(), testCase.statusCode, "status code")
	}
}

func TestUpdateUser(t *testing.T) {
	client, closeConn, _ := baseSetup.GetClient(grpcUrl)
	defer closeConn()

	for _, testCase := range updateUserTestCases {
		ctx := context.Background()

		respReg, _ := client.Register(ctx, &auth.RegisterRequest{
			Login:    testCase.oldUser.Login,
			Password: testCase.oldUser.Password,
		})
		requireNotNil(t, respReg, "response")

		resp, err := client.UpdateUser(ctx, &auth.UpdateUserRequest{
			Token:       respReg.Token,
			NewLogin:    testCase.newUser.Login,
			NewPassword: testCase.newUser.Password,
			RoleKey: sugar.If(
				testCase.newUser.Role == model.UserRoles.Admin,
				&config.Env.AdminKey,
				nil,
			),
		})
		st, ok := status.FromError(err)
		requireOk(t, ok, err)
		assertEqual(t, st.Code(), testCase.statusCode, "UpdateUser status code")
		if st.Code() != codes.OK {
			continue
		}

		claims, err := service.ParseJwtToken(resp.Token)
		requireNoError(t, err)
		assertEqual(t, claims.Login, testCase.newUser.Login, "UpdateUser user login")
		assertEqual(t, claims.Role, testCase.newUser.Role, "UpdateUser user role")

		respLog, err := client.Login(ctx, &auth.LoginRequest{
			Login:    testCase.newUser.Login,
			Password: testCase.newUser.Password,
		})
		st, ok = status.FromError(err)
		requireOk(t, ok, err)

		claims, err = service.ParseJwtToken(respLog.Token)
		requireNoError(t, err)
		assertEqual(t, claims.Login, testCase.newUser.Login, "Login user login")
		assertEqual(t, claims.Role, testCase.newUser.Role, "Login user role")
	}
}