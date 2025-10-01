package server

import (
	"common/api/payment"
	"common/pkg/consts"
	convpkg "common/pkg/conversion"
	"common/pkg/helper"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	"context"
	"payment/internal/config"
	"payment/internal/service"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	GrpcServer     *grpc.Server
	service        *service.PaymentService
	marshaler      convpkg.KafkaMarshaler
	kafkaConnector *serverpkg.KafkaConnector
}

func NewPaymentServer(opts grpc.ServerOption) *PaymentServer {
	s := &PaymentServer{
		service:        service.NewPaymentService(),
		marshaler:      convpkg.NewKafkaMarshaler(config.Env.KafkaSerialization),
		kafkaConnector: serverpkg.NewKafkaConnector(log.Loggers.Event),
	}

	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	payment.RegisterPaymentServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *PaymentServer) initKafka() {
	fetchHandlers := map[consts.TopicName]serverpkg.KafkaFetchHandlerFunc{
		consts.Topics.OrderCreated: s.StartPayment,
	}
	readerTopics := helper.MapKeys(fetchHandlers)
	readerConfig := kafka.ReaderConfig{
		Brokers:          config.Env.KafkaBrokerHosts,
		GroupID:          "payment_group",
		StartOffset:      kafka.LastOffset,
	}

	writerConfig := kafka.WriterConfig{
		Brokers:      config.Env.KafkaBrokerHosts,
		RequiredAcks: int(kafka.RequireAll),
	}
	writerTopics := []consts.TopicName{
		consts.Topics.OrderConfirmed,
		consts.Topics.OrderCanceled,
	}

	s.kafkaConnector.ConnectAll(&readerTopics, &readerConfig, &writerTopics, &writerConfig)
	for topic, handler := range fetchHandlers {
		s.kafkaConnector.AttachFetchHandler(topic, handler)
	}
}

func (s *PaymentServer) StartPayment(
	ctx context.Context,
	msg kafka.Message,
	commit serverpkg.KafkaHandlerCommit,
) error {
	event, err := s.marshaler.UnmarshalOrderEvent(msg)
	if err != nil {
		return err
	}

	if time.Since(msg.Time) > config.Env.PayTimeout {
		newMsg := kafka.Message{Key: msg.Key, Value: msg.Value}
		err = s.kafkaConnector.Writers[consts.Topics.OrderCanceled].WriteMessages(ctx, newMsg)
		return err
	}

	err = s.service.StartPayment(ctx, *event)  // TODO: идемпотентно обрабатывать

	commit()

	// TODO: добавить Outbox
	newMsg := kafka.Message{Key: msg.Key, Value: msg.Value}
	if err == nil {
		err = s.kafkaConnector.Writers[consts.Topics.OrderConfirmed].WriteMessages(ctx, newMsg)
	} else {
		err = s.kafkaConnector.Writers[consts.Topics.OrderCanceled].WriteMessages(ctx, newMsg)
	}

	return err
}

var allServers []*PaymentServer

func Deinit() {
	for _, s := range allServers {
		s.kafkaConnector.DisconnectAll()
		s.GrpcServer.GracefulStop()
	}
}