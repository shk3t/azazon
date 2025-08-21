package notificationtest

import (
	"common/pkg/log"
	commSetup "common/pkg/setup"
	"common/pkg/sugar"
	"context"
	"fmt"
	"notification/internal/config"
	conv "notification/internal/conversion"
	"notification/internal/service"
	"notification/internal/setup"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
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

	cmd, err := commSetup.ServerUp(workDir, grpcUrl, logger)
	if err != nil {
		commSetup.ServerDown(cmd, logger)
		logger.Println(err)
		setup.DeinitAll()
		os.Exit(1)
	}

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")

	commSetup.ServerDown(cmd, logger)
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

	for i, testCase := range orderCreatedTestCases {
		ctx := context.Background()

		payload, err := proto.Marshal(conv.OrderEventProto(&testCase.order))
		require.NoError(err)

		err = writer.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(strconv.Itoa(testCase.order.OrderId)),
				Value: payload,
			},
		)
		require.NoError(err)

		if i == 0 {
			time.Sleep(3 * time.Second)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		messages, err := service.ReadEmails(
			service.FmtUserById(testCase.order.UserId),
		)
		require.NoError(err)
		require.True(len(messages) > 0, "No new messages recieved")
		msg := messages[len(messages)-1]
		require.Contains(msg, service.FmtUserById(testCase.order.UserId))
		require.Contains(msg, fmt.Sprintf("Order %d", testCase.order.OrderId))
		require.Contains(msg, "created")
	}
}