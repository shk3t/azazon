package server

import (
	"common/api/payment"
	"common/pkg/helper"
	"context"
	"payment/internal/config"
	"payment/internal/service"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	GrpcServer   *grpc.Server
	service      service.PaymentService
	kafkaCancel  context.CancelFunc
	kafkaReaders []*kafka.Reader
	kafkaWriters []*kafka.Writer
}

func NewPaymentServer(opts grpc.ServerOption) *PaymentServer {
	srv := &PaymentServer{
		service: *service.NewPaymentService(),
	}

	var kafkaCtx context.Context
	kafkaCtx, srv.kafkaCancel = context.WithCancel(context.Background())
	srv.initKafkaReaders(kafkaCtx)
	srv.initKafkaWriters(kafkaCtx)

	srv.GrpcServer = grpc.NewServer(opts)
	payment.RegisterPaymentServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

func (srv *PaymentServer) initKafkaReaders(ctx context.Context) {
	handlers := map[string]helper.KafkaMessageHandlerFunc{
		"order_created": srv.StartPayment,
	}

	for topic, handlerFunc := range handlers {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:          config.Env.KafkaBrokerHosts,
			Topic:            topic,
			GroupID:          "payment_group",
			StartOffset:      kafka.LastOffset,
			RebalanceTimeout: 2 * time.Second,
		})
		srv.kafkaReaders = append(srv.kafkaReaders, reader)

		go helper.KafkaBgRead(ctx, reader, handlerFunc, topic)
	}
}

func (srv *PaymentServer) initKafkaWriters(ctx context.Context) {
	topics := []string{"order_confirmed", "order_canceled"}

	for _, topic := range topics {
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers: config.Env.KafkaBrokerHosts,
			Topic:   topic,
		})
		srv.kafkaWriters = append(srv.kafkaWriters, writer)
	}
}

var allServers []*PaymentServer

func (s *PaymentServer) StartPayment(ctx context.Context, msg kafka.Message) error {
	return nil
}

func Deinit() {
	for _, srv := range allServers {
		srv.kafkaCancel()

		var wg sync.WaitGroup
		for _, writer := range srv.kafkaWriters {
			wg.Add(1)
			go func() {
				writer.Close()
				wg.Done()
			}()
		}
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