package authtest

import (
	"auth/internal/config"
	conv "auth/internal/conversion"
	"auth/internal/model"
	"auth/internal/setup"
	"common/api/auth"
	"common/pkg/consts"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	"common/pkg/service"
	setuppkg "common/pkg/setup"
	"common/pkg/sugar"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var connector *serverpkg.TestConnector

func TestMain(m *testing.M) {
	workDir := filepath.Dir(sugar.Default(os.Getwd()))
	os.Setenv(config.AppName+"_TEST", "true")

	err := setup.InitAll(workDir)
	if err != nil {
		setup.DeinitAll()
		panic(err)
	}

	logger := log.Loggers.Test
	grpcUrl := fmt.Sprintf("localhost:%d", config.Env.TestPort)

	cmd, err := setuppkg.ServerUp(config.AppName, workDir, grpcUrl, logger)
	if err != nil {
		setuppkg.ServerDown(cmd, logger)
		logger.Println(err)
		setup.DeinitAll()
		os.Exit(1)
	}

	connector = serverpkg.NewTestConnector(logger)
	connector.ConnectAll(
		map[consts.ServiceName]string{consts.Services.Auth: grpcUrl},
		nil, nil, nil, nil,
	)

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")

	setuppkg.ServerDown(cmd, logger)
	connector.DisconnectAll()
	setup.DeinitAll()
	os.Exit(exitCode)
}

func TestRegister(t *testing.T) {
	require := require.New(t)
	client, err := connector.GetAuthClient()
	require.NoError(err)

	for _, testCase := range registerTestCases {
		ctx := context.Background()

		resp, err := client.Register(ctx, conv.RegisterRequest(&testCase.request))
		st, ok := status.FromError(err)
		require.True(ok, err)

		require.Equal(testCase.statusCode, st.Code())
		if st.Code() != codes.OK {
			continue
		}

		claims, err := service.ParseJwtToken(resp.Token)
		require.NoError(err)

		require.Equal(testCase.response.Login, claims.Login)
	}
}

func TestLogin(t *testing.T) {
	require := require.New(t)
	client, err := connector.GetAuthClient()
	require.NoError(err)

	for _, testCase := range loginTestCases {
		ctx := context.Background()

		// It is supposed to work
		_, _ = client.Register(ctx, conv.RegisterRequest(&testCase.request))

		_, err := client.Login(ctx, conv.LoginRequest(&testCase.request))
		st, ok := status.FromError(err)
		require.True(ok, err)

		require.Equal(testCase.statusCode, st.Code())
	}
}

func TestValidateToken(t *testing.T) {
	require := require.New(t)
	client, err := connector.GetAuthClient()
	require.NoError(err)

	for _, testCase := range validateTokenTestCases {
		ctx := context.Background()

		respReg, _ := client.Register(ctx, conv.RegisterRequest(&testCase.registerRequest))
		require.NotNil(respReg)

		resp, err := client.ValidateToken(ctx, &auth.ValidateTokenRequest{
			Token: respReg.Token,
		})
		st, ok := status.FromError(err)
		require.True(ok, err)
		require.True(resp.Valid, errors.New("Invalid token"))

		require.Equal(testCase.statusCode, st.Code())
	}
}

func TestUpdateUser(t *testing.T) {
	require := require.New(t)
	client, err := connector.GetAuthClient()
	require.NoError(err)

	for _, testCase := range updateUserTestCases {
		ctx := context.Background()

		respReg, _ := client.Register(ctx, conv.RegisterRequest(&testCase.oldUser))
		require.NotNil(respReg)

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
		require.True(ok, err)
		require.Equal(testCase.statusCode, st.Code())
		if st.Code() != codes.OK {
			continue
		}

		claims, err := service.ParseJwtToken(resp.Token)
		require.NoError(err)
		require.Equal(testCase.newUser.Login, claims.Login)
		require.Equal(string(testCase.newUser.Role), claims.Role)

		respLog, err := client.Login(ctx, conv.LoginRequest(&testCase.newUser))
		st, ok = status.FromError(err)
		require.True(ok, err)

		claims, err = service.ParseJwtToken(respLog.Token)
		require.NoError(err)
		require.Equal(testCase.newUser.Login, claims.Login)
		require.Equal(string(testCase.newUser.Role), claims.Role)
	}
}