package notificationtest

import (
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"base/pkg/sugar"
	"context"
	"encoding/json"
	"fmt"
	"notification/internal/config"
	"notification/internal/service"
	"notification/internal/setup"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
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
		setup.DeinitAll()
		os.Exit(1)
	}

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")

	baseSetup.ServerDown(cmd, logger)
	setup.DeinitAll()
	os.Exit(exitCode)
}

func TestOrderCreated(t *testing.T) {
	require := require.New(t)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: config.Env.KafkaBrokerHosts,
		Topic:   "order_created",
	})
	defer writer.Close()

	for _, testCase := range orderCreatedTestCases {
		ctx := context.Background()

		payload, err := json.Marshal(testCase.order)
		require.NoError(err)

		err = writer.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(strconv.Itoa(testCase.order.Id)),
				Value: payload,
			},
		)
		require.NoError(err)

		time.Sleep(10 * time.Second)

		messages, err := service.ReadEmails(
			service.FmtUserById(testCase.order.UserId),
		)
		require.NoError(err)
		require.True(len(messages) > 0, "No new messages recieved")
		msg := messages[len(messages)-1]
		require.Contains(msg, service.FmtUserById(testCase.order.UserId))
		require.Contains(msg, fmt.Sprintf("Order %d", testCase.order.Id))
		require.Contains(msg, "created")
	}
}