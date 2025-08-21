package helper

import (
	"common/pkg/log"
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageHandlerFunc func(ctx context.Context, msg kafka.Message) error

func KafkaBgRead(
	ctx context.Context,
	reader *kafka.Reader,
	handlerFunc KafkaMessageHandlerFunc,
	topic string,
) {
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Loggers.Event.Println(err)
			if errors.Is(err, context.Canceled) {
				return
			}
			continue
		}

		if err = handlerFunc(ctx, msg); err != nil {
			log.Loggers.Event.Println(err)
		}
		log.Loggers.Event.Printf(
			"Message %s in topic %s handled",
			string(msg.Value), topic,
		)
	}
}