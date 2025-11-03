package scripts

import (
	"common/api/auth"
	"common/api/order"
	"common/api/stock"
	"common/internal/config"
	"common/pkg/consts"
	"common/pkg/helper"
	logpkg "common/pkg/log"
	serverpkg "common/pkg/server"
	"context"
	"fmt"
)

func newDefaultTestManager() *serverpkg.TestManager {
	manager := serverpkg.NewTestManager(logpkg.Loggers.Test)

	if config.Env.VirtualRuntime == helper.VirtualRuntimes.Kubernetes {
		if err := manager.UseAsIngress(config.Env.Domain, "cert/tls.crt"); err != nil {
			panic(err)
		}
	}
	manager.ConnectGrpc(map[consts.ServiceName]string{
		consts.Services.Auth:  config.Env.GrpcUrls.Auth,
		consts.Services.Order: config.Env.GrpcUrls.Order,
		consts.Services.Stock: config.Env.GrpcUrls.Stock,
	})

	manager.InitMarshaler(config.Env.KafkaSerialization)

	return manager
}

func CreateDefaultUsers() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()

	regAdminResp, err := authClient.Register(ctx, &auth.RegisterRequest{
		Login:    "admin",
		Password: "admin123",
	})
	if err != nil {
		panic(err)
	}
	_, err = authClient.Register(ctx, &auth.RegisterRequest{
		Login:    "customer",
		Password: "customer123",
	})
	if err != nil {
		panic(err)
	}
	_, err = authClient.UpdateUser(ctx, &auth.UpdateUserRequest{
		Token:   regAdminResp.Token,
		RoleKey: &config.Env.AdminKey,
	})
	if err != nil {
		panic(err)
	}
}

func FillStocks() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()
	stockClient, _ := manager.GetStockClient()

	loginAdminResp, err := authClient.Login(ctx, &auth.LoginRequest{
		Login:    "admin",
		Password: "admin123",
	})

	if err != nil {
		panic(err)
	}
	adminToken := loginAdminResp.Token

	saveBallResp, err := stockClient.SaveProduct(ctx, &stock.SaveProductRequest{
		Token:        adminToken,
		ProductName:  "ball",
		ProductPrice: 1000.00,
	})
	if err != nil {
		panic(err)
	}
	stockClient.ChangeStockQuantity(ctx, &stock.ChangeStockQuantityRequest{
		Token:         adminToken,
		ProductId:     saveBallResp.Stock.Product.Id,
		QuantityDelta: 50,
	})
	saveLampResp, err := stockClient.SaveProduct(ctx, &stock.SaveProductRequest{
		Token:        adminToken,
		ProductName:  "lamp",
		ProductPrice: 2000.00,
	})
	if err != nil {
		panic(err)
	}
	stockClient.ChangeStockQuantity(ctx, &stock.ChangeStockQuantityRequest{
		Token:         adminToken,
		ProductId:     saveLampResp.Stock.Product.Id,
		QuantityDelta: 10,
	})
}

func MakeGoodOrder() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()
	orderClient, _ := manager.GetOrderClient()

	loginCustomerResp, _ := authClient.Login(ctx, &auth.LoginRequest{
		Login:    "customer",
		Password: "customer123",
	})
	token := loginCustomerResp.Token

	_, err := orderClient.CreateOrder(ctx, &order.CreateOrderRequest{
		Token: token,
		Items: []*order.Item{
			&order.Item{ProductId: 1, Quantity: 3},
			&order.Item{ProductId: 2, Quantity: 2},
		},
		Address: "Nice street, 420 house",
	})
	if err != nil {
		panic(err)
	}
}

func MakeBadOrderProductNotEnoughStocks() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()
	orderClient, _ := manager.GetOrderClient()

	loginCustomerResp, _ := authClient.Login(ctx, &auth.LoginRequest{
		Login:    "customer",
		Password: "customer123",
	})
	token := loginCustomerResp.Token

	_, err := orderClient.CreateOrder(ctx, &order.CreateOrderRequest{
		Token: token,
		Items: []*order.Item{
			&order.Item{ProductId: 1, Quantity: 3},
			&order.Item{ProductId: 2, Quantity: 20},
		},
		Address: "Ultra street, 420 000 house",
	})
	if err == nil {
		panic("It just works")
	}
	fmt.Println(err)
}

func MakeBadOrderProductDoesNotExist() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()
	orderClient, _ := manager.GetOrderClient()

	loginCustomerResp, _ := authClient.Login(ctx, &auth.LoginRequest{
		Login:    "customer",
		Password: "customer123",
	})
	token := loginCustomerResp.Token

	_, err := orderClient.CreateOrder(ctx, &order.CreateOrderRequest{
		Token: token,
		Items: []*order.Item{
			&order.Item{ProductId: 1, Quantity: 3},
			&order.Item{ProductId: 2, Quantity: 2},
			&order.Item{ProductId: 3, Quantity: 1},
		},
		Address: "Wall street, 1337 house",
	})
	if err == nil {
		panic("It just works")
	}
	fmt.Println(err)
}

func MakeBadOrderProductPaymentError() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()
	orderClient, _ := manager.GetOrderClient()

	loginCustomerResp, _ := authClient.Login(ctx, &auth.LoginRequest{
		Login:    "admin",
		Password: "admin123",
	})
	token := loginCustomerResp.Token

	_, err := orderClient.CreateOrder(ctx, &order.CreateOrderRequest{
		Token: token,
		Items: []*order.Item{
			&order.Item{ProductId: 1, Quantity: 20},
			&order.Item{ProductId: 2, Quantity: 5},
		},
		Address: "Rich street, X house",
	})
	if err != nil {
		panic(err)
	}
}