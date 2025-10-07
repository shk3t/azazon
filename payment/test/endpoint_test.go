package notificationtest

import (
	"common/pkg/consts"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	setuppkg "common/pkg/setup"
	"common/pkg/sugar"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"payment/internal/config"
	"payment/internal/setup"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

var manager *serverpkg.TestManager

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

	manager = serverpkg.NewTestManager(logger)
	manager.ConnectKafka(
		&[]consts.TopicName{consts.Topics.OrderConfirmed, consts.Topics.OrderCanceled},
		&kafka.ReaderConfig{
			Brokers:          config.Env.KafkaBrokerHosts,
			GroupID:          "payment_test_group",
			StartOffset:      kafka.LastOffset,
		},
		&[]consts.TopicName{consts.Topics.OrderCreated},
		&kafka.WriterConfig{Brokers: config.Env.KafkaBrokerHosts},
	)
	manager.InitMarshaler(config.Env.KafkaSerialization)

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")

	setuppkg.ServerDown(cmd, logger)
	manager.Close()
	setup.DeinitAll()
	os.Exit(exitCode)
}

func TestStartPayment(t *testing.T) {
	require := require.New(t)
	createdWriter := manager.GetKafkaWriter(consts.Topics.OrderCreated)
	confirmedReader := manager.GetKafkaReader(consts.Topics.OrderConfirmed, true)
	canceledReader := manager.GetKafkaReader(consts.Topics.OrderCanceled, true)

	for _, testCase := range startPaymentTestCases {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		msg := manager.MarshalOrderEvent(&testCase.event)
		err := createdWriter.WriteMessages(ctx, msg)
		require.NoError(err)

		reader := sugar.If(testCase.event.FullPrice <= balance, confirmedReader, canceledReader)
		msg, err = reader.ReadMessage(ctx)
		require.NoError(err)

		resultOrder, err := manager.UnmarshalOrderEvent(msg)
		require.NoError(err)

		require.Equal(testCase.event.OrderId, resultOrder.OrderId)
		require.Equal(testCase.event.UserId, resultOrder.UserId)
		require.Equal(testCase.event.FullPrice, resultOrder.FullPrice)

		cancel()
	}
}