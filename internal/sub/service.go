package sub

import (
	"context"
	"faceit-backend-test/internal/notify"
)

type service struct {
	notificationManager notify.NotificationManager
}

var _ Service = (*service)(nil)

type ServiceOpts func(*service)

func NewService(opts ...ServiceOpts) *service {
	s := &service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithNotificationManager(manager notify.NotificationManager) ServiceOpts {
	return func(s *service) {
		s.notificationManager = manager
	}
}

func (s service) Subscribe(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
	sub, err := s.notificationManager.Subscribe(request.Type, request.Callback, request.Secret)
	if err != nil {
		return SubscribeResponse{}, SubscribeError(err.Error())
	}

	return SubscribeResponse{
		Id:        sub.Id,
		Type:      request.Type,
		Status:    sub.Status,
		CreatedAt: sub.CreatedAt,
	}, nil
}
