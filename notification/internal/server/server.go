package server

import (
	"base/api/notification"
	"base/api/order"
	conv "base/pkg/conversion"
	"base/pkg/log"
	"context"
	"errors"
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
			Brokers:          config.Env.KafkaBrokerHosts,
			Topic:            topic,
			GroupID:          "notification_group",
			StartOffset:      kafka.LastOffset,
			RebalanceTimeout: 1 * time.Second,
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
				log.Loggers.Event.Printf(
					"Message %s in topic %s handled",
					string(msg.Value), topic,
				)
			}
		}()
	}
}

type kafkaMessageHandlerFunc func(ctx context.Context, msg kafka.Message) error

func (s *NotificationServer) HandleOrderCreated(ctx context.Context, msg kafka.Message) error {
	var in order.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

func (s *NotificationServer) HandleOrderConfirmed(ctx context.Context, msg kafka.Message) error {
	var in order.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

func (s *NotificationServer) HandleOrderCanceled(ctx context.Context, msg kafka.Message) error {
	var in order.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}
	return s.service.HandleOrderCreated(ctx, *conv.OrderEventModel(&in))
}

var allServers []*NotificationServer

func Deinit() {
	for _, srv := range allServers {
		srv.readerCancelFunc()

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