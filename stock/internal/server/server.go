package server

import (
	"common/api/auth"
	stockapi "common/api/stock"
	"common/pkg/consts"
	convpkg "common/pkg/conversion"
	"common/pkg/grpcutil"
	serverpkg "common/pkg/server"
	servicepkg "common/pkg/service"
	"common/pkg/sugar"
	"context"
	"net/http"
	"stock/internal/config"
	conv "stock/internal/conversion"
	"stock/internal/model"
	"stock/internal/service"
	"time"

	"google.golang.org/grpc"
)

var NewErr = grpcutil.NewGrpcError
var NewInternalErr = grpcutil.NewInternalGrpcError

type StockServer struct {
	stockapi.UnimplementedStockServiceServer
	GrpcServer    *grpc.Server
	service       *service.StockService
	marshaler     convpkg.KafkaMarshaler
	grpcConnector *serverpkg.GrpcConnector
	bgCancel      context.CancelFunc
}

func NewStockServer(opts grpc.ServerOption) *StockServer {
	s := &StockServer{
		service:       service.NewStockService(),
		marshaler:     convpkg.NewKafkaMarshaler(config.Env.KafkaSerialization),
		grpcConnector: serverpkg.NewGrpcConnector(),
	}

	s.initBackgroundJobs()
	s.initGrpcClients()

	s.GrpcServer = grpc.NewServer(opts)
	stockapi.RegisterStockServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *StockServer) initGrpcClients() {
	s.grpcConnector.Connect(consts.Services.Auth, config.Env.GrpcUrls.Auth)
}

func (s *StockServer) initBackgroundJobs() {
	ticker := time.NewTicker(time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	s.bgCancel = cancel

	go func() {
		for {
			select {
			case <-ticker.C:
				s.CancelReserve(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *StockServer) SaveProduct(
	ctx context.Context,
	in *stockapi.SaveProductRequest,
) (*stockapi.SaveProductResponse, error) {
	authClient, _ := s.grpcConnector.GetAuthClient()
	resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	} else if !resp.Valid {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	} else if claims, _ := servicepkg.ParseJwtToken(in.Token); !claims.IsAdmin() {
		return nil, NewErr(http.StatusForbidden, "Not enough user permissions")
	}

	product, err := s.service.SaveProduct(ctx, model.Product{
		Id:    sugar.If(in.ProductId != nil, int(*in.ProductId), 0),
		Name:  in.ProductName,
		Price: in.ProductPrice,
	})
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	stock, err := s.service.GetStockInfo(ctx, product.Id)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	return &stockapi.SaveProductResponse{
		Stock: conv.StockProto(product, stock),
	}, nil
}

func (s *StockServer) IncreaseStockQuantity(
	ctx context.Context,
	in *stockapi.ChangeStockQuantityRequest,
) (*stockapi.ChangeStockQuantityResponse, error) {
	authClient, _ := s.grpcConnector.GetAuthClient()
	resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	} else if !resp.Valid {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	} else if claims, _ := servicepkg.ParseJwtToken(in.Token); !claims.IsAdmin() {
		return nil, NewErr(http.StatusForbidden, "Not enough user permissions")
	}

	stock, err := s.service.ChangeStockQuantity(ctx, nil, int(in.ProductId), int(in.QuantityDelta))
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	product, err := s.service.GetProductInfo(ctx, int(in.ProductId))
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	return &stockapi.ChangeStockQuantityResponse{
		Stock: conv.StockProto(product, stock),
	}, nil
}

func (s *StockServer) Reserve(
	ctx context.Context,
	in *stockapi.ReserveRequest,
) (*stockapi.ReserveResponse, error) {
	authClient, _ := s.grpcConnector.GetAuthClient()
	resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	} else if !resp.Valid {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	}

	claims, err := servicepkg.ParseJwtToken(in.Token)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	reserve := model.Reserve{
		UserId:    claims.UserId,
		OrderId:   int(in.OrderId),
		ProductId: int(in.ProductId),
		Quantity:  int(in.Quantity),
	}
	undo := in.Undo != nil && *in.Undo
	stock, err := s.service.Reserve(ctx, reserve, undo)
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	product, err := s.service.GetProductInfo(ctx, int(in.ProductId))
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	return &stockapi.ReserveResponse{
		Stock: conv.StockProto(product, stock),
	}, nil
}

func (s *StockServer) GetStockInfo(
	ctx context.Context,
	in *stockapi.GetStockInfoRequest,
) (*stockapi.GetStockInfoResponse, error) {
	stock, err := s.service.GetStockInfo(ctx, int(in.ProductId))
	if err != nil {
		return nil, err.Grpc()
	}

	product, err := s.service.GetProductInfo(ctx, int(in.ProductId))
	if err != nil {
		return nil, err.Grpc()
	}

	return &stockapi.GetStockInfoResponse{
		Stock: conv.StockProto(product, stock),
	}, nil
}

func (s *StockServer) DeleteProduct(
	ctx context.Context,
	in *stockapi.DeleteProductRequest,
) (*stockapi.DeleteProductResponse, error) {
	authClient, _ := s.grpcConnector.GetAuthClient()
	resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	} else if !resp.Valid {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	} else if claims, _ := servicepkg.ParseJwtToken(in.Token); !claims.IsAdmin() {
		return nil, NewErr(http.StatusForbidden, "Not enough user permissions")
	}

	err = s.service.DeleteProduct(ctx, int(in.ProductId))
	if v, ok := err.(*grpcutil.ServiceError); ok && v != nil {
		return nil, v.Grpc()
	}

	return &stockapi.DeleteProductResponse{}, nil
}

func (s *StockServer) CancelReserve(ctx context.Context) {
	deprecationTime := time.Now().Add(-config.Env.ReserveTimeout)
	s.service.UndoOldReserves(ctx, deprecationTime)
}

var allServers []*StockServer

func Deinit() {
	for _, s := range allServers {
		s.GrpcServer.GracefulStop()
		s.bgCancel()
	}
}