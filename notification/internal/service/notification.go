package service

import (
	"base/pkg/grpcutil"
	"base/pkg/model"
	"context"
	"fmt"
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
	SendEmail(
		FmtUserById(body.UserId),
		fmt.Sprintf("Order %d created", body.OrderId),
	)
	return nil
}

func (s *NotificationService) HandleOrderConfirmed(
	ctx context.Context,
	body model.OrderEvent,
) error {
	SendEmail(
		FmtUserById(body.UserId),
		fmt.Sprintf("Order %d confirmed", body.OrderId),
	)
	return nil
}

func (s *NotificationService) HandleOrderCanceled(
	ctx context.Context,
	body model.OrderEvent,
) error {
	SendEmail(
		FmtUserById(body.UserId),
		fmt.Sprintf("Order %d canceled", body.OrderId),
	)
	return nil
}