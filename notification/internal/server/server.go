package server

import (
	"base/api/notification"
	"base/pkg/log"
	"base/pkg/model"
	"context"
	"encoding/json"
	"errors"
	"notification/internal/config"
	"notification/internal/service"
	"sync"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	GrpcServer       *grpc.Server
	service          service.NotificationService
	kafkaReaders     []*kafka.Reader
	readerCancelFunc context.CancelFunc
}

func NewNotificationServer(serverOpts grpc.ServerOption) *NotificationServer {
	srv := &NotificationServer{
		service: *service.NewNotificationService(),
	}

	srv.initKafkaReaders()

	srv.GrpcServer = grpc.NewServer(serverOpts)
	notification.RegisterNotificationServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

func (srv *NotificationServer) initKafkaReaders() {
	handlers := map[string]kafkaMessageHandlerFunc{
		"order_created":   srv.HandleOrderCreated,
		"order_confirmed": srv.HandleOrderConfirmed,
		"order_canceled":  srv.HandleOrderCanceled,
	}

	var readerCtx context.Context
	readerCtx, srv.readerCancelFunc = context.WithCancel(context.Background())
	srv.kafkaReaders = []*kafka.Reader{}

	for topic, handlerFunc := range handlers {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: config.Env.KafkaBrokerHosts,
			Topic:   topic,
			GroupID: "notification_group",
		})
		srv.kafkaReaders = append(srv.kafkaReaders, reader)

		go func() {
			for {
				msg, err := reader.ReadMessage(readerCtx)
				if err != nil {
					log.Loggers.Event.Println(err)
					if errors.Is(err, context.Canceled) {
						return
					}
					continue
				}

				if err = handlerFunc(readerCtx, msg); err != nil {
					log.Loggers.Event.Println(err)
				}
			}
		}()
	}
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
		srv.readerCancelFunc()

		var wg sync.WaitGroup
		for i, reader := range srv.kafkaReaders {
			wg.Add(1)
			go func() {
				reader.Close()
				log.Debug(i, "kafka reader closed")
				wg.Done()
			}()
		}
		wg.Wait()

		srv.GrpcServer.GracefulStop()
	}
}