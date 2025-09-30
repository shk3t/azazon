package notificationtest

import (
	"common/pkg/consts"
	conv "common/pkg/conversion"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	setuppkg "common/pkg/setup"
	"common/pkg/sugar"
	"context"
	"fmt"
	"notification/internal/config"
	"notification/internal/service"
	"notification/internal/setup"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

var connector *serverpkg.TestConnector
var marshaler conv.KafkaMarshaler

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
		nil,
		nil, nil,
		&[]consts.TopicName{consts.Topics.OrderCreated},
		&kafka.WriterConfig{Brokers: config.Env.KafkaBrokerHosts},
	)

	marshaler = conv.NewKafkaMarshaler(config.Env.KafkaSerialization)

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")

	setuppkg.ServerDown(cmd, logger)
	connector.DisconnectAll()
	setup.DeinitAll()
	os.Exit(exitCode)
}

func TestOrderCreated(t *testing.T) {
	require := require.New(t)
	createdWriter := connector.GetKafkaWriter(consts.Topics.OrderCreated)

	for i, testCase := range orderCreatedTestCases {
		ctx := context.Background()

		kMsg := marshaler.MarshalOrderEvent(&testCase.event)
		err := createdWriter.WriteMessages(ctx, kMsg)
		require.NoError(err)

		if i == 0 {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(100 * time.Millisecond)
		}

		messages, err := service.ReadEmails(
			service.FmtUserById(testCase.event.UserId),
		)
		require.NoError(err)
		require.True(len(messages) > 0, "No new messages recieved")
		eMsg := messages[len(messages)-1]
		require.Contains(eMsg, service.FmtUserById(testCase.event.UserId))
		require.Contains(eMsg, fmt.Sprintf("Order %d", testCase.event.OrderId))
		require.Contains(eMsg, "created")
	}
}