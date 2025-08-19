package server

import (
	"base/api/notification"
	"base/pkg/log"
	"base/pkg/model"
	"context"
	"encoding/json"
	"notification/internal/service"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	GrpcServer   *grpc.Server
	service      service.NotificationService
	kafkaReaders []*kafka.Reader
}

func NewNotificationServer(serverOpts grpc.ServerOption) *NotificationServer {
	srv := &NotificationServer{
		service: *service.NewNotificationService(),
	}

	srv.kafkaReaders = initKafkaReaders(srv)

	srv.GrpcServer = grpc.NewServer(serverOpts)
	notification.RegisterNotificationServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

func initKafkaReaders(srv *NotificationServer) []*kafka.Reader {
	handlers := map[string]kafkaMessageHandlerFunc{
		"order_created":   srv.HandleOrderCreated,
		"order_confirmed": srv.HandleOrderCreated,
		"order_canceled":  srv.HandleOrderCreated,
	}
	readers := []*kafka.Reader{}

	for topic, handlerFunc := range handlers {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			GroupID: "notification_group",
		})
		readers = append(readers, reader)

		go func() {
			ctx := context.Background()
			for {
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Loggers.Event.Println(err)
					continue
				}

				if err = handlerFunc(ctx, msg); err != nil {
					log.Loggers.Event.Println(err)
				}
			}
		}()
	}

	return readers
}

type kafkaMessageHandlerFunc func(ctx context.Context, msg kafka.Message) error

func (s *NotificationServer) HandleOrderCreated(ctx context.Context, msg kafka.Message) error {
	var in model.OrderEvent
	if err := json.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, in)
}

func (s *NotificationServer) HandleOrderConfirmed(ctx context.Context, msg kafka.Message) error {
	var in model.OrderEvent
	if err := json.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderConfirmed(ctx, in)
}

func (s *NotificationServer) HandleOrderCanceled(ctx context.Context, msg kafka.Message) error {
	var in model.OrderEvent
	if err := json.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCanceled(ctx, in)
}

var allServers []*NotificationServer

func Deinit() {
	for _, srv := range allServers {
		for _, reader := range srv.kafkaReaders {
			reader.Close()
		}

		srv.GrpcServer.GracefulStop()
	}
}