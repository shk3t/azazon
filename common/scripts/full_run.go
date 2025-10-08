package scripts

import (
	"common/api/auth"
	"common/api/order"
	"common/api/stock"
	"common/internal/config"
	"common/pkg/consts"
	logpkg "common/pkg/log"
	serverpkg "common/pkg/server"
	"context"
)

func newDefaultTestManager() *serverpkg.TestManager {
	manager := serverpkg.NewTestManager(logpkg.Loggers.Test)
	manager.ConnectGrpc(map[consts.ServiceName]string{
		consts.Services.Auth:  config.Env.GrpcUrls.Auth,
		consts.Services.Order: config.Env.GrpcUrls.Order,
		consts.Services.Stock: config.Env.GrpcUrls.Stock,
	})
	// allTopics := []consts.TopicName{
	// 	consts.Topics.OrderCreated,
	// 	consts.Topics.OrderConfirmed,
	// 	consts.Topics.OrderCanceled,
	// }
	// manager.ConnectKafka(
	// 	&allTopics,
	// 	&kafka.ReaderConfig{
	// 		Brokers:     config.Env.KafkaBrokerHosts,
	// 		GroupID:     "script_group",
	// 		StartOffset: kafka.LastOffset,
	// 	},
	// 	&allTopics,
	// 	&kafka.WriterConfig{Brokers: config.Env.KafkaBrokerHosts},
	// )
	manager.InitMarshaler(config.Env.KafkaSerialization)

	return manager
}

func CreateDefaultUsers() {
	ctx := context.Background()
	manager := newDefaultTestManager()
	defer manager.Close()

	authClient, _ := manager.GetAuthClient()

	regAdminResp, _ := authClient.Register(ctx, &auth.RegisterRequest{
		Login:    "admin",
		Password: "admin123",
	})
	_, err := authClient.Register(ctx, &auth.RegisterRequest{
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

	loginAdminResp, _ := authClient.Login(ctx, &auth.LoginRequest{
		Login:    "admin",
		Password: "admin123",
	})
	adminToken := loginAdminResp.Token

	saveBallResp, _ := stockClient.SaveProduct(ctx, &stock.SaveProductRequest{
		Token:        adminToken,
		ProductName:  "ball",
		ProductPrice: 1000.00,
	})
	stockClient.ChangeStockQuantity(ctx, &stock.ChangeStockQuantityRequest{
		Token:         adminToken,
		ProductId:     saveBallResp.Stock.Product.Id,
		QuantityDelta: 50,
	})
	saveLampResp, _ := stockClient.SaveProduct(ctx, &stock.SaveProductRequest{
		Token:        adminToken,
		ProductName:  "lamp",
		ProductPrice: 2000.00,
	})
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

// 1. When not enough stocks
// 2. When product does not exist
// 3. When payment error