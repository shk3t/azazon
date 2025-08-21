package server

import (
	"common/api/common"
	"common/api/notification"
	conv "common/pkg/conversion"
	"common/pkg/helper"
	"context"
	"notification/internal/config"
	"notification/internal/service"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	GrpcServer   *grpc.Server
	service      service.NotificationService
	kafkaCancel  context.CancelFunc
	kafkaReaders []*kafka.Reader
}

func NewNotificationServer(opts grpc.ServerOption) *NotificationServer {
	srv := &NotificationServer{
		service: *service.NewNotificationService(),
	}

	var kafkaCtx context.Context
	kafkaCtx, srv.kafkaCancel = context.WithCancel(context.Background())
	srv.initKafkaReaders(kafkaCtx)

	srv.GrpcServer = grpc.NewServer(opts)
	notification.RegisterNotificationServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

func (srv *NotificationServer) initKafkaReaders(ctx context.Context) {
	handlers := map[string]helper.KafkaMessageHandlerFunc{
		"order_created":   srv.HandleOrderCreated,
		"order_confirmed": srv.HandleOrderConfirmed,
		"order_canceled":  srv.HandleOrderCanceled,
	}

	for topic, handlerFunc := range handlers {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:          config.Env.KafkaBrokerHosts,
			Topic:            topic,
			GroupID:          "notification_group",
			StartOffset:      kafka.LastOffset,
			RebalanceTimeout: 2 * time.Second,
		})
		srv.kafkaReaders = append(srv.kafkaReaders, reader)

		go helper.KafkaBgRead(ctx, reader, handlerFunc, topic)
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
	for _, srv := range allServers {
		srv.kafkaCancel()

		var wg sync.WaitGroup
		for _, reader := range srv.kafkaReaders {
			wg.Add(1)
			go func() {
				reader.Close()
				wg.Done()
			}()
		}
		wg.Wait()

		srv.GrpcServer.GracefulStop()
	}
}