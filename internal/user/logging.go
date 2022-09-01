package user

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
		return serviceLoggingMiddleware{
			logger: logger,
			next:   s,
		}
	}
}

func (s serviceLoggingMiddleware) Create(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"endpoint": "Create",
		"request":  request,
	}).Debug("received request")
	var resp CreateUserResponse
	var err error
	defer func(start time.Time) {
		logger := s.logger.WithFields(logrus.Fields{
			"endpoint": "Create",
			"took":     time.Since(start).String(),
		})
		if err != nil {
			logger.WithField("error", err).Errorln("an error occurred")
			return
		}

		logger.WithField("response", resp).Debug()
	}(time.Now())
	resp, err = s.next.Create(ctx, request)
	return resp, err
}

func (s serviceLoggingMiddleware) Update(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"endpoint": "Update",
		"request":  request,
	}).Debug("received request")
	var resp UpdateUserResponse
	var err error
	defer func(start time.Time) {
		logger := s.logger.WithFields(logrus.Fields{
			"endpoint": "Update",
			"took":     time.Since(start).String(),
		})
		if err != nil {
			logger.WithField("error", err).Errorln("an error occurred")
			return
		}

		logger.WithField("response", resp).Debug()
	}(time.Now())
	resp, err = s.next.Update(ctx, request)
	return resp, err
}

func (s serviceLoggingMiddleware) DeleteById(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"endpoint": "DeleteById",
		"request":  request,
	}).Debug("received request")
	var resp DeleteUserResponse
	var err error
	defer func(start time.Time) {
		logger := s.logger.WithFields(logrus.Fields{
			"endpoint": "Delete",
			"took":     time.Since(start).String(),
		})
		if err != nil {
			logger.WithField("error", err).Errorln("an error occurred")
			return
		}

		logger.WithField("response", resp).Debug()
	}(time.Now())
	resp, err = s.next.DeleteById(ctx, request)
	return resp, err
}

func (s serviceLoggingMiddleware) GetMany(ctx context.Context) (GetUsersManyResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"endpoint": "GetMany",
	}).Debug("received request")
	var resp GetUsersManyResponse
	var err error
	defer func(start time.Time) {
		logger := s.logger.WithFields(logrus.Fields{
			"endpoint": "GetMany",
			"took":     time.Since(start).String(),
		})
		if err != nil {
			logger.WithField("error", err).Errorln("an error occurred")
			return
		}

		logger.WithField("response", resp).Debug()
	}(time.Now())
	resp, err = s.next.GetMany(ctx)
	return resp, err
}
