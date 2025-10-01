package server

import (
	"common/pkg/consts"
	"context"
	"errors"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaReadHandlerFunc func(ctx context.Context, msg kafka.Message) error
type KafkaFetchHandlerFunc func(ctx context.Context, msg kafka.Message, commit KafkaHandlerCommit) error
type KafkaHandlerCommit func() error

type KafkaConnector struct {
	Readers   map[consts.TopicName]*kafka.Reader
	Writers   map[consts.TopicName]*kafka.Writer
	ctx       context.Context
	cancelCtx context.CancelFunc
	logger    *log.Logger
}

func NewKafkaConnector(logger *log.Logger) *KafkaConnector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConnector{
		Readers:   map[consts.TopicName]*kafka.Reader{},
		Writers:   map[consts.TopicName]*kafka.Writer{},
		ctx:       ctx,
		cancelCtx: cancel,
		logger:    logger,
	}
}

func (c *KafkaConnector) ConnectAll(
	readerTopics *[]consts.TopicName,
	readerConfig *kafka.ReaderConfig,
	writeTopics *[]consts.TopicName,
	writerConfig *kafka.WriterConfig,
) {
	if readerTopics != nil && readerConfig != nil {
		for _, topic := range *readerTopics {
			readerConfig.Topic = string(topic)
			c.Readers[topic] = kafka.NewReader(*readerConfig)
		}
	}

	if writeTopics != nil && writerConfig != nil {
		for _, topic := range *writeTopics {
			writerConfig.Topic = string(topic)
			c.Writers[topic] = kafka.NewWriter(*writerConfig)
		}
	}
}

func (c *KafkaConnector) AttachReadHandler(
	topic consts.TopicName,
	handlerFunc KafkaReadHandlerFunc,
) {
	reader := c.Readers[topic]

	go func() {
		for {
			msg, err := reader.ReadMessage(c.ctx)
			if err != nil {
				c.logger.Println(err)
				if errors.Is(err, context.Canceled) {
					return
				}
				continue
			}

			if err = handlerFunc(c.ctx, msg); err != nil {
				c.logger.Println(err)
			}
			c.logger.Printf(
				"Message %s in topic %s handled",
				string(msg.Value), topic,
			)
		}
	}()
}

func (c *KafkaConnector) AttachFetchHandler(
	topic consts.TopicName,
	handlerFunc KafkaFetchHandlerFunc,
) {
	reader := c.Readers[topic]

	go func() {
		for {
			msg, err := reader.FetchMessage(c.ctx)
			if err != nil {
				c.logger.Println(err)
				if errors.Is(err, context.Canceled) {
					return
				}
				continue
			}

			err = handlerFunc(c.ctx, msg, func() error { return reader.CommitMessages(c.ctx, msg) })
			if err != nil {
				c.logger.Println(err)
			}
			c.logger.Printf(
				"Message %s in topic %s handled",
				string(msg.Value), topic,
			)
		}
	}()
}

func (c *KafkaConnector) DisconnectAll() {
	c.cancelCtx()

	var wg sync.WaitGroup
	for _, writer := range c.Writers {
		wg.Add(1)
		go func() {
			writer.Close()
			wg.Done()
		}()
	}
	for _, reader := range c.Readers {
		wg.Add(1)
		go func() {
			reader.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}