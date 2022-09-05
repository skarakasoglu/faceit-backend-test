package sub

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type ServiceLoggingMiddleware func(Service) Service

type serviceLoggingMiddleware struct {
	logger *logrus.Logger
	next   Service
}

var _ Service = (*serviceLoggingMiddleware)(nil)

func NewServiceLoggingMiddleware(logger *logrus.Logger) ServiceLoggingMiddleware {
	return func(s Service) Service {
		return &serviceLoggingMiddleware{
			logger: logger,
			next:   s,
		}
	}
}

func (s *serviceLoggingMiddleware) Subscribe(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"service":  "SubscribeService",
		"endpoint": "Subscribe",
		"request":  request,
	}).Debug("received request")
	var resp SubscribeResponse
	var err error
	defer func(start time.Time) {
		logger := s.logger.WithFields(logrus.Fields{
			"service":  "SubscribeService",
			"endpoint": "Subscribe",
			"took":     time.Since(start).String(),
		})
		if err != nil {
			logger.WithField("error", err).Errorln("an error occurred")
			return
		}

		logger.WithField("response", resp).Debug()
	}(time.Now())
	resp, err = s.next.Subscribe(ctx, request)
	return resp, err
}
