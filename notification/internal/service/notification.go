package service

import "base/pkg/grpcutil"

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}