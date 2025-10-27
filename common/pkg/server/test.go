package server

import (
	"common/api/auth"
	"common/api/notification"
	"common/api/order"
	"common/api/payment"
	"common/api/stock"
	"common/pkg/consts"
	conv "common/pkg/conversion"
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type grpcConnector interface {
	UseAsIngress(authority string, tlsCertPath string) error
	GetAuthClient() (client auth.AuthServiceClient, err error)
	GetNotificationClient() (client notification.NotificationServiceClient, err error)
	GetOrderClient() (client order.OrderServiceClient, err error)
	GetPaymentClient() (client payment.PaymentServiceClient, err error)
	GetStockClient() (client stock.StockServiceClient, err error)
}

type TestManager struct {
	logger *log.Logger
	grpcConnector
	kafkaConnector *KafkaConnector
	conv.KafkaMarshaler
}

func NewTestManager(logger *log.Logger) *TestManager {
	return &TestManager{
		logger:        logger,
		grpcConnector: NewGrpcConnector(),
	}
}

func (c *TestManager) ConnectGrpc(
	grpcUrls map[consts.ServiceName]string,
) {
	for serviceName, url := range grpcUrls {
		c.grpcConnector.(*GrpcConnector).Connect(serviceName, url)
	}
}

func (c *TestManager) ConnectKafka(
	topicsToRead *[]consts.TopicName,
	readerConfig *kafka.ReaderConfig,
	topicsToWrite *[]consts.TopicName,
	writerConfig *kafka.WriterConfig,
) {
	c.kafkaConnector = NewKafkaConnector(c.logger)
	c.kafkaConnector.ConnectAll(topicsToRead, readerConfig, topicsToWrite, writerConfig)
}

func (c *TestManager) InitMarshaler(serializationMethod string) {
	c.KafkaMarshaler = conv.NewKafkaMarshaler(serializationMethod)
}

func (c *TestManager) GetKafkaReader(topic consts.TopicName, sink bool) *kafka.Reader {
	reader := c.kafkaConnector.Readers[topic]

	if sink {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		for {
			if _, err := reader.ReadMessage(ctx); err != nil {
				break
			}
		}
		cancel()
	}

	return reader
}

func (c *TestManager) GetKafkaWriter(topic consts.TopicName) *kafka.Writer {
	return c.kafkaConnector.Writers[topic]
}

func (c *TestManager) Close() {
	if grpcConn, ok := c.grpcConnector.(*GrpcConnector); ok {
		grpcConn.DisconnectAll()
	}
	if c.kafkaConnector != nil {
		c.kafkaConnector.DisconnectAll()
	}
}