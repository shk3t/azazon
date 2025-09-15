package server

import (
	"common/api/auth"
	"common/api/notification"
	"common/api/order"
	"common/api/payment"
	"common/api/stock"
	"common/pkg/consts"
	"common/pkg/helper"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type closeConnectionFunc func() error

type grpcClients struct {
	auth         auth.AuthServiceClient
	notification notification.NotificationServiceClient
	order        order.OrderServiceClient
	payment      payment.PaymentServiceClient
	stock        stock.StockServiceClient
}

type GrpcConnector struct {
	clients    grpcClients
	closeFuncs []closeConnectionFunc
}

func NewGrpcConnector() *GrpcConnector {
	return &GrpcConnector{}
}

func (c *GrpcConnector) Connect(serviceName consts.ServiceName, url string) error {
	switch serviceName {
	case consts.Services.Auth:
		conn, err := grpc.NewClient(
			url, grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}

		c.clients.auth = auth.NewAuthServiceClient(conn)
		c.closeFuncs = append(c.closeFuncs, conn.Close)

	case consts.Services.Notification:
		conn, err := grpc.NewClient(
			url, grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}

		c.clients.notification = notification.NewNotificationServiceClient(conn)
		c.closeFuncs = append(c.closeFuncs, conn.Close)

	case consts.Services.Order:
		conn, err := grpc.NewClient(
			url, grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}

		c.clients.order = order.NewOrderServiceClient(conn)
		c.closeFuncs = append(c.closeFuncs, conn.Close)

	case consts.Services.Payment:
		conn, err := grpc.NewClient(
			url, grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}

		c.clients.payment = payment.NewPaymentServiceClient(conn)
		c.closeFuncs = append(c.closeFuncs, conn.Close)

	case consts.Services.Stock:
		conn, err := grpc.NewClient(
			url, grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}

		c.clients.stock = stock.NewStockServiceClient(conn)
		c.closeFuncs = append(c.closeFuncs, conn.Close)
	}

	return nil
}

func (c *GrpcConnector) GetAuthClient() (client auth.AuthServiceClient, err error) {
	if c.clients.auth == nil {
		return nil, NewNotInitedServiceError(consts.Services.Auth)
	}

	return c.clients.auth, nil
}

func (c *GrpcConnector) GetNotificationClient() (client notification.NotificationServiceClient, err error) {
	if c.clients.notification == nil {
		return nil, NewNotInitedServiceError(consts.Services.Notification)
	}

	return c.clients.notification, nil
}

func (c *GrpcConnector) GetOrderClient() (client order.OrderServiceClient, err error) {
	if c.clients.order == nil {
		return nil, NewNotInitedServiceError(consts.Services.Order)
	}

	return c.clients.order, nil
}

func (c *GrpcConnector) GetPaymentClient() (client payment.PaymentServiceClient, err error) {
	if c.clients.payment == nil {
		return nil, NewNotInitedServiceError(consts.Services.Payment)
	}

	return c.clients.payment, nil
}

func (c *GrpcConnector) GetStockClient() (client stock.StockServiceClient, err error) {
	if c.clients.stock == nil {
		return nil, NewNotInitedServiceError(consts.Services.Stock)
	}

	return c.clients.stock, nil
}

func NewNotInitedServiceError(serviceName consts.ServiceName) error {
	return fmt.Errorf(
		"%s service client is not connected",
		helper.Capitalize(string(serviceName)),
	)
}

func (c *GrpcConnector) DisconnectAll() {
	wg := sync.WaitGroup{}

	for _, closeFunc := range c.closeFuncs {
		wg.Add(1)
		go func() {
			closeFunc()
			wg.Done()
		}()
	}

	wg.Wait()
}