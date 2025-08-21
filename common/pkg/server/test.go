package server

import (
	"common/api/auth"
	"common/pkg/setup"

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
	topicsToRead *[]string,
	readerConfig *kafka.ReaderConfig,
	topicsToWrite *[]string,
	writerConfig *kafka.WriterConfig,
) {
	c.grpcUrl = grpcUrl
	c.kafkaConnector = NewKafkaConnector()
	c.kafkaConnector.Connect(topicsToRead, readerConfig, topicsToWrite, writerConfig)
}

func (c *TestConnector) GetGrpcClient() (
	client auth.AuthServiceClient,
	closeConn func() error,
	err error,
) {
	return setup.GetGrpcClient(c.grpcUrl)
}

func (c *TestConnector) GetKafkaReader(topic string) *kafka.Reader {
	return c.kafkaConnector.Readers[topic]
}

func (c *TestConnector) GetKafkaWriter(topic string) *kafka.Writer {
	return c.kafkaConnector.Writers[topic]
}

func (c *TestConnector) Disconnect() {
	c.kafkaConnector.Disconnect()
}