package conversion

import (
	"common/api/common"
	"common/pkg/model"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func orderEventModel(o *common.OrderEvent) *model.OrderEvent {
	return &model.OrderEvent{
		OrderId:   int(o.OrderId),
		UserId:    int(o.UserId),
		FullPrice: o.FullPrice,
	}
}

func orderEventProto(o *model.OrderEvent) *common.OrderEvent {
	return &common.OrderEvent{
		OrderId:   int64(o.OrderId),
		UserId:    int64(o.UserId),
		FullPrice: o.FullPrice,
	}
}

type KafkaMarshaler interface {
	MarshalOrderEvent(event *model.OrderEvent) kafka.Message
	UnmarshalOrderEvent(msg kafka.Message) (*model.OrderEvent, error)
}

type protoKafkaMarshaler struct{}
type jsonKafkaMarshaler struct{}

func NewKafkaMarshaler(method string) KafkaMarshaler {
	method = strings.ToLower(method)

	switch method {
	case "proto":
		return &protoKafkaMarshaler{}
	case "json":
		return &jsonKafkaMarshaler{}
	default:
		panic("Serialization method is not implemented")
	}
}

func (m *protoKafkaMarshaler) MarshalOrderEvent(event *model.OrderEvent) kafka.Message {
	protoEvent := orderEventProto(event)
	encodedEvent, _ := proto.Marshal(protoEvent)
	return kafka.Message{
		Key:   []byte(strconv.Itoa(event.OrderId)),
		Value: encodedEvent,
	}
}

func (m *protoKafkaMarshaler) UnmarshalOrderEvent(msg kafka.Message) (*model.OrderEvent, error) {
	var protoEvent common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &protoEvent); err != nil {
		return nil, err
	}
	return orderEventModel(&protoEvent), nil
}

func (m *jsonKafkaMarshaler) MarshalOrderEvent(event *model.OrderEvent) kafka.Message {
	encodedEvent, _ := json.Marshal(event)
	return kafka.Message{
		Key:   []byte(strconv.Itoa(event.OrderId)),
		Value: encodedEvent,
	}
}

func (m *jsonKafkaMarshaler) UnmarshalOrderEvent(msg kafka.Message) (*model.OrderEvent, error) {
	var event model.OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return nil, err
	}
	return &event, nil
}