package server

import (
	"common/api/common"
	"common/api/notification"
	"common/pkg/consts"
	"common/pkg/helper"
	"common/pkg/log"
	commServer "common/pkg/server"
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
	kafkaConnector *commServer.KafkaConnector
}

func NewNotificationServer(opts grpc.ServerOption) *NotificationServer {
	s := &NotificationServer{
		service:        service.NewNotificationService(),
		kafkaConnector: commServer.NewKafkaConnector(log.Loggers.Event),
	}

	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	notification.RegisterNotificationServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *NotificationServer) initKafka() {
	readerHandlers := map[consts.TopicName]commServer.KafkaMessageHandlerFunc{
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

	s.kafkaConnector.ConnectAll(&readerTopics, &readerConfig, nil, nil)
	for topic, handler := range readerHandlers {
		s.kafkaConnector.AttachReaderHandler(topic, handler)
		log.Debug("ATTACHED")
	}
}

func (s *NotificationServer) HandleOrderCreated(ctx context.Context, msg kafka.Message) error {
	log.Debug("CALLED")
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

func (s *NotificationServer) HandleOrderConfirmed(ctx context.Context, msg kafka.Message) error {
	log.Debug("CALLED")
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

func (s *NotificationServer) HandleOrderCanceled(ctx context.Context, msg kafka.Message) error {
	log.Debug("CALLED")
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

var allServers []*NotificationServer

func Deinit() {
	for _, s := range allServers {
		s.kafkaConnector.DisconnectAll()
		s.GrpcServer.GracefulStop()
	}
}