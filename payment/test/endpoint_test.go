package notificationtest

import (
	"common/api/common"
	"common/pkg/consts"
	"common/pkg/log"
	commServer "common/pkg/server"
	commSetup "common/pkg/setup"
	"common/pkg/sugar"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"payment/internal/config"
	conv "payment/internal/conversion"
	"payment/internal/setup"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var connector = commServer.NewTestConnector()

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

	cmd, err := commSetup.ServerUp(workDir, grpcUrl, logger)
	if err != nil {
		commSetup.ServerDown(cmd, logger)
		logger.Println(err)
		setup.DeinitAll()
		os.Exit(1)
	}

	connector.Connect(
		grpcUrl,
		&[]string{consts.Topics.OrderConfirmed, consts.Topics.OrderCanceled},
		&kafka.ReaderConfig{
			Brokers:          config.Env.KafkaBrokerHosts,
			GroupID:          "payment_test_group",
			StartOffset:      kafka.LastOffset,
			RebalanceTimeout: 2 * time.Second,
		},
		&[]string{consts.Topics.OrderCreated},
		&kafka.WriterConfig{Brokers: config.Env.KafkaBrokerHosts},
	)

	logger.Println("Running tests...")
	exitCode := m.Run()
	logger.Println("Test run finished")

	commSetup.ServerDown(cmd, logger)
	connector.Disconnect()
	setup.DeinitAll()
	os.Exit(exitCode)
}

func TestStartPayment(t *testing.T) {
	require := require.New(t)
	createdWriter := connector.GetKafkaWriter(consts.Topics.OrderCreated)
	confirmedReader := connector.GetKafkaReader(consts.Topics.OrderConfirmed)
	canceledReader := connector.GetKafkaReader(consts.Topics.OrderCanceled)

	for _, testCase := range startPaymentTestCases {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		payload, err := proto.Marshal(conv.OrderEventProto(&testCase.order))
		require.NoError(err)

		err = createdWriter.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(strconv.Itoa(testCase.order.OrderId)),
				Value: payload,
			},
		)
		require.NoError(err)

		reader := sugar.If(testCase.order.FullPrice <= balance, confirmedReader, canceledReader)
		msg, err := reader.ReadMessage(ctx)
		require.NoError(err)

		var out common.OrderEvent
		err = proto.Unmarshal(msg.Value, &out)
		require.NoError(err)

		resultOrder := conv.OrderEventModel(&out)
		require.Equal(testCase.order.OrderId, resultOrder.OrderId)
		require.Equal(testCase.order.UserId, resultOrder.UserId)
		require.Equal(testCase.order.FullPrice, resultOrder.FullPrice)

		cancel()
	}
}