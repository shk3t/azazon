package server

import (
	"common/api/auth"
	orderapi "common/api/order"
	"common/api/stock"
	"common/pkg/consts"
	convpkg "common/pkg/conversion"
	"common/pkg/grpcutil"
	"common/pkg/helper"
	"common/pkg/log"
	serverpkg "common/pkg/server"
	servicepkg "common/pkg/service"
	"context"
	"net/http"
	"order/internal/config"
	conv "order/internal/conversion"
	db "order/internal/database"
	"order/internal/model"
	"order/internal/service"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var NewErr = grpcutil.NewGrpcError
var NewInternalErr = grpcutil.NewInternalGrpcError

type OrderServer struct {
	orderapi.UnimplementedOrderServiceServer
	GrpcServer     *grpc.Server
	service        *service.OrderService
	marshaler      convpkg.KafkaMarshaler
	grpcConnector  *serverpkg.GrpcConnector
	kafkaConnector *serverpkg.KafkaConnector
	outbox         *serverpkg.TransactionalOutboxManager
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
	orderapi.RegisterOrderServiceServer(s.GrpcServer, s)

	s.outbox = serverpkg.NewTransactionalOutboxManager(
		db.ConnPool, s.kafkaConnector, log.Loggers.Event,
	)

	allServers = append(allServers, s)
	return s
}

func (s *OrderServer) initGrpcClients() {
	s.grpcConnector.Connect(consts.Services.Auth, config.Env.GrpcUrls.Auth)
	s.grpcConnector.Connect(consts.Services.Stock, config.Env.GrpcUrls.Stock)
}

func (s *OrderServer) initKafka() {
	fetchHandlers := map[consts.TopicName]serverpkg.KafkaFetchHandlerFunc{
		consts.Topics.OrderCanceled:  s.CancelOrder,
		consts.Topics.OrderConfirmed: s.ConfirmOrder,
	}
	readerTopics := helper.MapKeys(fetchHandlers)
	readerConfig := kafka.ReaderConfig{
		Brokers:     config.Env.KafkaBrokerHosts,
		GroupID:     "order_group",
		StartOffset: kafka.LastOffset,
	}

	writerConfig := kafka.WriterConfig{
		Brokers:      config.Env.KafkaBrokerHosts,
		RequiredAcks: int(kafka.RequireAll),
	}
	writerTopics := []consts.TopicName{
		consts.Topics.OrderCreated,
	}

	s.kafkaConnector.ConnectAll(&readerTopics, &readerConfig, &writerTopics, &writerConfig)
	for topic, handler := range fetchHandlers {
		s.kafkaConnector.AttachFetchHandler(topic, handler)
	}
}

func (s *OrderServer) CreateOrder(
	ctx context.Context,
	in *orderapi.CreateOrderRequest,
) (*orderapi.CreateOrderResponse, error) {
	authClient, _ := s.grpcConnector.GetAuthClient()
	stockClient, _ := s.grpcConnector.GetStockClient()

	resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	} else if !resp.Valid {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	}

	claims, _ := servicepkg.ParseJwtToken(in.Token)

	order := model.Order{
		Id:      s.service.GetNextOrderId(ctx),
		UserId:  claims.UserId,
		Status:  model.CreatedStatus,
		Address: in.Address,
		Track:   uuid.New().String(),
	}
	for _, item := range in.Items {
		order.Items = append(order.Items, model.Item{
			ProductId: int(item.ProductId),
			Quantity:  int(item.Quantity),
		})
	}

	fullPrice100 := atomic.Int64{}

	eg, egCtx := errgroup.WithContext(ctx)
	for _, item := range in.Items {
		eg.Go(
			func() error {
				resp, err := stockClient.Reserve(egCtx, &stock.ReserveRequest{
					Token:     in.Token,
					OrderId:   int64(order.Id),
					ProductId: item.ProductId,
					Quantity:  item.Quantity,
				})
				if err != nil {
					return err
				}

				fullPrice100.Add(int64(resp.Stock.Product.Price) * item.Quantity * 100)

				return nil
			},
		)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	newOrder, err := s.service.CreateOrder(ctx, tx, order)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	orderEvent := conv.OrderEvent(newOrder)
	orderEvent.FullPrice = float64(fullPrice100.Load()) / 100
	msg := s.marshaler.MarshalOrderEvent(orderEvent)
	s.outbox.Enqueue(ctx, tx, consts.Topics.OrderCreated, msg)
	if err = tx.Commit(ctx); err != nil {
		return nil, NewInternalErr(err)
	}
	s.outbox.Notify()

	return &orderapi.CreateOrderResponse{OrderId: int64(order.Id)}, nil
}

func (s *OrderServer) GetOrderInfo(
	ctx context.Context,
	in *orderapi.GetOrderInfoRequest,
) (*orderapi.GetOrderInfoResponse, error) {
	order, err := s.service.GetOrderInfo(ctx, int(in.OrderId))
	if err != nil {
		return nil, err.Grpc()
	}
	return conv.GetOrderInfoResponse(order), nil
}

func (s *OrderServer) CancelOrder(
	ctx context.Context,
	msg kafka.Message,
	commit serverpkg.KafkaHandlerCommit,
) error {
	orderEvent, err := s.marshaler.UnmarshalOrderEvent(msg)
	if err != nil {
		return err
	}

	order, err := s.service.GetOrderInfo(ctx, orderEvent.OrderId)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return err
	}

	order.Status = model.CanceledStatus
	_, err = s.service.UpdateOrder(ctx, nil, *order)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return err
	}

	commit()
	return nil
}

func (s *OrderServer) ConfirmOrder(
	ctx context.Context,
	msg kafka.Message,
	commit serverpkg.KafkaHandlerCommit,
) error {
	orderEvent, err := s.marshaler.UnmarshalOrderEvent(msg)
	if err != nil {
		return err
	}

	order, err := s.service.GetOrderInfo(ctx, orderEvent.OrderId)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return err
	}

	order.Status = model.ConfirmedStatus
	_, err = s.service.UpdateOrder(ctx, nil, *order)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return err
	}

	commit()
	return nil
}

var allServers []*OrderServer

func Deinit() {
	for _, s := range allServers {
		s.outbox.Close()
		s.grpcConnector.DisconnectAll()
		s.kafkaConnector.DisconnectAll()
		s.GrpcServer.GracefulStop()
	}
}