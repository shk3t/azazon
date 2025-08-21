package server

import (
	"common/api/common"
	"common/api/notification"
	"common/pkg/consts"
	"common/pkg/helper"
	connServer "common/pkg/server"
	"context"
	"notification/internal/config"
	conv "notification/internal/conversion"
	"notification/internal/service"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	GrpcServer     *grpc.Server
	service        *service.NotificationService
	kafkaConnector *connServer.KafkaConnector
}

func NewNotificationServer(opts grpc.ServerOption) *NotificationServer {
	s := &NotificationServer{
		service:        service.NewNotificationService(),
		kafkaConnector: connServer.NewKafkaConnector(),
	}

	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	notification.RegisterNotificationServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *NotificationServer) initKafka() {
	readerHandlers := map[string]connServer.KafkaMessageHandlerFunc{
		consts.Topics.OrderCreated:   s.HandleOrderCreated,
		consts.Topics.OrderConfirmed: s.HandleOrderConfirmed,
		consts.Topics.OrderCanceled:  s.HandleOrderCanceled,
	}
	readerTopics := helper.MapKeys(readerHandlers)
	readerConfig := kafka.ReaderConfig{
		Brokers:          config.Env.KafkaBrokerHosts,
		GroupID:          "notification_group",
		StartOffset:      kafka.LastOffset,
		RebalanceTimeout: 2 * time.Second,
	}

	s.kafkaConnector.Connect(&readerTopics, &readerConfig, nil, nil)
	for topic, handler := range readerHandlers {
		s.kafkaConnector.AttachReaderHandler(topic, handler)
	}
}

func (s *NotificationServer) HandleOrderCreated(ctx context.Context, msg kafka.Message) error {
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

func (s *NotificationServer) HandleOrderConfirmed(ctx context.Context, msg kafka.Message) error {
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

func (s *NotificationServer) HandleOrderCanceled(ctx context.Context, msg kafka.Message) error {
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

var allServers []*NotificationServer

func Deinit() {
	for _, s := range allServers {
		s.kafkaConnector.Disconnect()
		s.GrpcServer.GracefulStop()
	}
}