package server

import (
	"common/api/common"
	"common/api/payment"
	"common/pkg/consts"
	"common/pkg/helper"
	connServer "common/pkg/server"
	"context"
	"payment/internal/config"
	conv "payment/internal/conversion"
	"payment/internal/service"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	GrpcServer     *grpc.Server
	service        *service.PaymentService
	kafkaConnector *connServer.KafkaConnector
}

func NewPaymentServer(opts grpc.ServerOption) *PaymentServer {
	s := &PaymentServer{
		service:        service.NewPaymentService(),
		kafkaConnector: connServer.NewKafkaConnector(),
	}

	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	payment.RegisterPaymentServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *PaymentServer) initKafka() {
	readerHandlers := map[string]connServer.KafkaMessageHandlerFunc{
		consts.Topics.OrderCreated: s.StartPayment,
	}
	readerTopics := helper.MapKeys(readerHandlers)
	readerConfig := kafka.ReaderConfig{
		Brokers:          config.Env.KafkaBrokerHosts,
		GroupID:          "payment_group",
		StartOffset:      kafka.LastOffset,
		RebalanceTimeout: 2 * time.Second,
	}

	writerConfig := kafka.WriterConfig{
		Brokers: config.Env.KafkaBrokerHosts,
	}
	writerTopics := []string{
		consts.Topics.OrderConfirmed,
		consts.Topics.OrderCanceled,
	}

	s.kafkaConnector.Connect(&readerTopics, &readerConfig, &writerTopics, &writerConfig)
	for topic, handler := range readerHandlers {
		s.kafkaConnector.AttachReaderHandler(topic, handler)
	}
}

func (s *PaymentServer) StartPayment(ctx context.Context, msg kafka.Message) error {
	var in common.OrderEvent
	if err := proto.Unmarshal(msg.Value, &in); err != nil {
		return err
	}

	err := s.service.StartPayment(ctx, *conv.OrderEventModel(&in))

	if err == nil {
		s.kafkaConnector.Writers[consts.Topics.OrderConfirmed].WriteMessages(ctx, msg)
	} else {
		s.kafkaConnector.Writers[consts.Topics.OrderCanceled].WriteMessages(ctx, msg)
	}

	return err
}

var allServers []*PaymentServer

func Deinit() {
	for _, s := range allServers {
		s.kafkaConnector.Disconnect()
		s.GrpcServer.GracefulStop()
	}
}