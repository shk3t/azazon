package server

import (
	"common/api/common"
	"common/api/payment"
	"common/pkg/consts"
	"common/pkg/helper"
	"context"
	"payment/internal/config"
	conv "payment/internal/conversion"
	"payment/internal/service"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	GrpcServer   *grpc.Server
	service      service.PaymentService
	kafkaCancel  context.CancelFunc
	kafkaReaders map[string]*kafka.Reader
	kafkaWriters map[string]*kafka.Writer
}

func NewPaymentServer(opts grpc.ServerOption) *PaymentServer {
	srv := &PaymentServer{
		service:      *service.NewPaymentService(),
		kafkaReaders: map[string]*kafka.Reader{},
		kafkaWriters: map[string]*kafka.Writer{},
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
		consts.Topics.OrderCreated: srv.StartPayment,
	}

	for topic, handlerFunc := range handlers {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:          config.Env.KafkaBrokerHosts,
			Topic:            topic,
			GroupID:          "payment_group",
			StartOffset:      kafka.LastOffset,
			RebalanceTimeout: 2 * time.Second,
		})
		srv.kafkaReaders[topic] = reader

		go helper.KafkaBgRead(ctx, reader, handlerFunc, topic)
	}
}

func (srv *PaymentServer) initKafkaWriters(ctx context.Context) {
	topics := []string{consts.Topics.OrderConfirmed, consts.Topics.OrderCanceled}

	for _, topic := range topics {
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers: config.Env.KafkaBrokerHosts,
			Topic:   topic,
		})
		srv.kafkaWriters[topic] = writer
	}
}

var allServers []*PaymentServer

func (s *PaymentServer) StartPayment(ctx context.Context, msg kafka.Message) error {
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}

	err := s.service.StartPayment(ctx, *conv.OrderEventModel(&in))

	if err == nil {
		s.kafkaWriters[consts.Topics.OrderConfirmed].WriteMessages(ctx, msg)
	} else {
		s.kafkaWriters[consts.Topics.OrderCanceled].WriteMessages(ctx, msg)
	}

	return err
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