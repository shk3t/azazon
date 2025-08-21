package server

import (
	"common/api/auth"
	"common/pkg/setup"
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type TestConnector struct {
	grpcUrl        string
	kafkaConnector *KafkaConnector
}

func NewTestConnector() *TestConnector {
	return &TestConnector{}
}

func (c *TestConnector) Connect(
	grpcUrl string,
	logger *log.Logger,
	topicsToRead *[]string,
	readerConfig *kafka.ReaderConfig,
	topicsToWrite *[]string,
	writerConfig *kafka.WriterConfig,
) {
	c.grpcUrl = grpcUrl
	c.kafkaConnector = NewKafkaConnector(logger)
	c.kafkaConnector.Connect(topicsToRead, readerConfig, topicsToWrite, writerConfig)
}

func (c *TestConnector) GetGrpcClient() (
	client auth.AuthServiceClient,
	closeConn func() error,
	err error,
) {
	return setup.GetGrpcClient(c.grpcUrl)
}

func (c *TestConnector) GetKafkaReader(topic string, sink bool) *kafka.Reader {
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

func (c *TestConnector) GetKafkaWriter(topic string) *kafka.Writer {
	return c.kafkaConnector.Writers[topic]
}

func (c *TestConnector) Disconnect() {
	c.kafkaConnector.Disconnect()
}