package service

import (
	"base/pkg/grpcutil"
	"base/pkg/log"
	"base/pkg/model"
	"context"
)

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) HandleOrderCreated(
	ctx context.Context,
	body model.OrderEvent,
) error {
	log.Loggers.Debug.Printf("Order created: %v\n", body)
	return nil
}

func (s *NotificationService) HandleOrderConfirmed(
	ctx context.Context,
	body model.OrderEvent,
) error {
	log.Loggers.Debug.Printf("Order confirmed: %v\n", body)
	return nil
}

func (s *NotificationService) HandleOrderCanceled(
	ctx context.Context,
	body model.OrderEvent,
) error {
	log.Loggers.Debug.Printf("Order canceled: %v\n", body)
	return nil
}