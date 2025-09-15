package server

import (
	"common/api/auth"
	"common/api/notification"
	"common/api/order"
	"common/api/payment"
	"common/api/stock"
	"common/pkg/consts"
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type TestConnector struct {
	logger         *log.Logger
	grpcConnector  *GrpcConnector
	kafkaConnector *KafkaConnector
}

func NewTestConnector(logger *log.Logger) *TestConnector {
	return &TestConnector{
		logger: logger,
	}
}

func (c *TestConnector) ConnectAll(
	grpcUrls map[consts.ServiceName]string,
	topicsToRead *[]consts.TopicName,
	readerConfig *kafka.ReaderConfig,
	topicsToWrite *[]consts.TopicName,
	writerConfig *kafka.WriterConfig,
) {
	c.grpcConnector = NewGrpcConnector()
	for serviceName, url := range grpcUrls {
		c.grpcConnector.Connect(serviceName, url)
	}

	c.kafkaConnector = NewKafkaConnector(c.logger)
	c.kafkaConnector.ConnectAll(topicsToRead, readerConfig, topicsToWrite, writerConfig)
}

func (c *TestConnector) GetAuthClient() (client auth.AuthServiceClient, err error) {
	return c.grpcConnector.GetAuthClient()
}

func (c *TestConnector) GetNotificationClient() (client notification.NotificationServiceClient, err error) {
	return c.grpcConnector.GetNotificationClient()
}

func (c *TestConnector) GetOrderClient() (client order.OrderServiceClient, err error) {
	return c.grpcConnector.GetOrderClient()
}

func (c *TestConnector) GetPaymentClient() (client payment.PaymentServiceClient, err error) {
	return c.grpcConnector.GetPaymentClient()
}

func (c *TestConnector) GetStockClient() (client stock.StockServiceClient, err error) {
	return c.grpcConnector.GetStockClient()
}

func (c *TestConnector) GetKafkaReader(topic consts.TopicName, sink bool) *kafka.Reader {
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

func (c *TestConnector) GetKafkaWriter(topic consts.TopicName) *kafka.Writer {
	return c.kafkaConnector.Writers[topic]
}

func (c *TestConnector) DisconnectAll() {
	c.grpcConnector.DisconnectAll()
	c.kafkaConnector.DisconnectAll()
}