package server

import (
	"common/api/order"
	"common/pkg/consts"
	commServer "common/pkg/server"
	"context"
	"order/internal/config"
	"order/internal/service"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	GrpcServer     *grpc.Server
	service        *service.OrderService
	kafkaConnector *commServer.KafkaConnector
}

func NewOrderServer(opts grpc.ServerOption) *OrderServer {
	s := &OrderServer{
		service: service.NewOrderService(),
	}

	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	order.RegisterOrderServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *OrderServer) initKafka() {
	writerConfig := kafka.WriterConfig{Brokers: config.Env.KafkaBrokerHosts}
	writerTopics := []consts.TopicName{consts.Topics.OrderCreated}
	s.kafkaConnector.ConnectAll(nil, nil, &writerTopics, &writerConfig)
}

func (s *OrderServer) CreateOrder(
	ctx context.Context,
	in *order.CreateOrderRequest,
) (*order.CreateOrderResponse, error) {

	// TODO: check stocks

	var err error
	var msg kafka.Message

	// TODO: create order

	if err != nil {
		s.kafkaConnector.Writers[consts.Topics.OrderCreated].WriteMessages(ctx, msg)
	}

	return nil, nil
}

func (s *OrderServer) GetOrderInfo(
	ctx context.Context,
	in *order.GetOrderInfoRequest,
) (*order.GetOrderInfoResponse, error) {
	// TODO: conversions
	return nil, nil
}

var allServers []*OrderServer

func Deinit() {
	for _, s := range allServers {
		s.kafkaConnector.DisconnectAll()
		s.GrpcServer.GracefulStop()
	}
}