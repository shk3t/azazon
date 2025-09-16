package server

import (
	"common/api/auth"
	"common/api/order"
	"common/api/stock"
	"common/pkg/consts"
	convpkg "common/pkg/conversion"
	"common/pkg/grpcutil"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	servicepkg "common/pkg/service"
	"context"
	"fmt"
	"net/http"
	"order/internal/config"
	conv "order/internal/conversion"
	"order/internal/model"
	"sync/atomic"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"order/internal/service"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

var NewErr = grpcutil.NewGrpcError
var NewInternalErr = grpcutil.NewInternalGrpcError

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	GrpcServer     *grpc.Server
	service        *service.OrderService
	marshaler      convpkg.KafkaMarshaler
	grpcConnector  *serverpkg.GrpcConnector
	kafkaConnector *serverpkg.KafkaConnector
}

func NewOrderServer(opts grpc.ServerOption) *OrderServer {
	s := &OrderServer{
		service:        service.NewOrderService(),
		marshaler:      convpkg.NewKafkaMarshaler(config.Env.KafkaSerialization),
		grpcConnector:  serverpkg.NewGrpcConnector(),
		kafkaConnector: serverpkg.NewKafkaConnector(log.Loggers.Event),
	}

	s.initGrpcClients()
	s.initKafka()

	s.GrpcServer = grpc.NewServer(opts)
	order.RegisterOrderServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *OrderServer) initGrpcClients() {
	s.grpcConnector.Connect(consts.Services.Auth, config.Env.GrpcUrls.Auth)
	s.grpcConnector.Connect(consts.Services.Stock, config.Env.GrpcUrls.Stock)
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
	authClient, _ := s.grpcConnector.GetAuthClient()
	stockClient, _ := s.grpcConnector.GetStockClient()

	resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	}
	if !resp.Valid {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	}

	claims, _ := servicepkg.ParseJwtToken(in.Token)

	order := model.Order{
		UserId:  claims.UserId,
		Address: in.Address,
		Track:   uuid.New().String(),
	}
	for _, item := range in.Items {
		order.Items = append(order.Items, model.Item{
			Id:       int(item.Id),
			Quantity: int(item.Quantity),
		})
	}

	fullPrice := atomic.Int64{}

	eg, egCtx := errgroup.WithContext(ctx)
	for _, item := range in.Items {
		eg.Go(
			func() error {
				resp, err := stockClient.GetStockInfo(egCtx, &stock.GetStockInfoRequest{
					ProductId: item.Id,
				})
				if err != nil {
					return err
				}

				if resp.Stock.Quantity < item.Quantity {
					return NewErr(
						http.StatusBadRequest,
						fmt.Sprintf("Not enough item_%d quantity in stock", item.Id),
					)
				}

				fullPrice.Add(resp.Stock.ProductPrice * item.Quantity)

				return nil
			},
		)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	order.Id, err = s.service.CreateOrder(ctx, order) // TODO
	if err != nil {
		return nil, err.(*grpcutil.ServiceError).Grpc()
	}

	orderEvent := conv.OrderEvent(&order)
	orderEvent.FullPrice = int(fullPrice.Load())

	msg := s.marshaler.MarshalOrderEvent(orderEvent)
	s.kafkaConnector.Writers[consts.Topics.OrderCreated].WriteMessages(ctx, msg)

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
		s.grpcConnector.DisconnectAll()
		s.kafkaConnector.DisconnectAll()
		s.GrpcServer.GracefulStop()
	}
}