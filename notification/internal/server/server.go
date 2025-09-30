package server

import (
	"common/api/notification"
	"common/pkg/consts"
	convpkg "common/pkg/conversion"
	"common/pkg/helper"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	"context"
	"notification/internal/config"
	"notification/internal/service"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	GrpcServer     *grpc.Server
	service        *service.NotificationService
	marshaler      convpkg.KafkaMarshaler
	kafkaConnector *serverpkg.KafkaConnector
}

func NewNotificationServer(opts grpc.ServerOption) *NotificationServer {
	s := &NotificationServer{
		service:        service.NewNotificationService(),
		marshaler:      convpkg.NewKafkaMarshaler(config.Env.KafkaSerialization),
		kafkaConnector: serverpkg.NewKafkaConnector(log.Loggers.Event),
	}

	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	notification.RegisterNotificationServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *NotificationServer) initKafka() {
	readerHandlers := map[consts.TopicName]serverpkg.KafkaMessageHandlerFunc{
		consts.Topics.OrderCreated:    s.HandleOrderCreated,
		consts.Topics.OrderConfirmed:  s.HandleOrderConfirmed,
		consts.Topics.OrderCancelling: s.HandleOrderCanceled,
	}
	readerTopics := helper.MapKeys(readerHandlers)
	readerConfig := kafka.ReaderConfig{
		Brokers:     config.Env.KafkaBrokerHosts,
		GroupID:     "notification_group",
		StartOffset: kafka.LastOffset, // comment in, comment out if kafka bugs in tests
	}

	s.kafkaConnector.ConnectAll(&readerTopics, &readerConfig, nil, nil)
	for topic, handler := range readerHandlers {
		s.kafkaConnector.AttachReaderHandler(topic, handler)
	}
}

func (s *NotificationServer) HandleOrderCreated(ctx context.Context, msg kafka.Message) error {
	event, err := s.marshaler.UnmarshalOrderEvent(msg)
	if err != nil {
		return nil
	}
	return s.service.HandleOrderCreated(ctx, *event)
}

func (s *NotificationServer) HandleOrderConfirmed(ctx context.Context, msg kafka.Message) error {
	event, err := s.marshaler.UnmarshalOrderEvent(msg)
	if err != nil {
		return nil
	}
	return s.service.HandleOrderConfirmed(ctx, *event)
}

func (s *NotificationServer) HandleOrderCanceled(ctx context.Context, msg kafka.Message) error {
	event, err := s.marshaler.UnmarshalOrderEvent(msg)
	if err != nil {
		return nil
	}
	return s.service.HandleOrderCanceled(ctx, *event)
}

var allServers []*NotificationServer

func Deinit() {
	for _, s := range allServers {
		s.kafkaConnector.DisconnectAll()
		s.GrpcServer.GracefulStop()
	}
}