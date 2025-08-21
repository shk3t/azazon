package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageHandlerFunc func(ctx context.Context, msg kafka.Message) error

type KafkaConnector struct {
	Readers   map[string]*kafka.Reader
	Writers   map[string]*kafka.Writer
	ctx       context.Context
	cancelCtx context.CancelFunc
	logger    *log.Logger
}

func NewKafkaConnector(logger *log.Logger) *KafkaConnector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConnector{
		Readers:   map[string]*kafka.Reader{},
		Writers:   map[string]*kafka.Writer{},
		ctx:       ctx,
		cancelCtx: cancel,
		logger:    logger,
	}
}

func (c *KafkaConnector) Connect(
	readerTopics *[]string,
	readerConfig *kafka.ReaderConfig,
	writeTopics *[]string,
	writerConfig *kafka.WriterConfig,
) {
	if readerTopics != nil && readerConfig != nil {
		for _, topic := range *readerTopics {
			readerConfig.Topic = topic
			c.Readers[topic] = kafka.NewReader(*readerConfig)
		}
	}

	if writeTopics != nil && writerConfig != nil {
		for _, topic := range *writeTopics {
			writerConfig.Topic = topic
			c.Writers[topic] = kafka.NewWriter(*writerConfig)
		}
	}
}

func (c *KafkaConnector) AttachReaderHandler(topic string, handlerFunc KafkaMessageHandlerFunc) {
	reader := c.Readers[topic]
	c.logger.Println(topic, "INITEDDD")  // TODO

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
			c.logger.Println("Я ПРОШЕЛ ДАЛЬШЕ!") // TODO

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

func (c *KafkaConnector) Disconnect() {
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