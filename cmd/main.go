package main

import (
	"context"
	"faceit-backend-test/internal/config"
	"faceit-backend-test/internal/health"
	"faceit-backend-test/internal/router"
	"faceit-backend-test/internal/user"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	cfg                      config.Config
	logger                   *logrus.Logger
	gracefulShutdownDuration = 5
)

func main() {
	s := &Service{}
	defer s.Shutdown()

	initConfig()
	initLogger()
	db := initDb()
	routes := initRoutes(db)

	go func() {
		if err := initHTTPHandler(s, routes...); err != nil {
			logger.WithFields(logrus.Fields{
				"transport": "http",
				"error":     err.Error(),
			}).Error("error occurred while serving")
			s.Shutdown()
		}
	}()

	signalHandler := make(chan os.Signal)
	signal.Notify(signalHandler, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

	receivedSignal := <-signalHandler
	logger.WithFields(logrus.Fields{
		"signal": receivedSignal,
	}).Info("received signal")
}

type Service struct {
	server *http.Server
}

func (s *Service) Shutdown() {
	logger.WithFields(logrus.Fields{
		"gracefulShutdownPeriod": gracefulShutdownDuration,
	}).Info("shutting down the app")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gracefulShutdownDuration)*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("server shutdown properly")
}

func resolveLogLevel(l string) (int, error) {
	switch strings.ToUpper(l) {
	case "PANIC":
		return 0, nil
	case "FATAL":
		return 1, nil
	case "ERROR":
		return 2, nil
	case "WARN":
		return 3, nil
	case "INFO":
		return 4, nil
	case "DEBUG":
		return 5, nil
	case "TRACE":
		return 6, nil
	}

	return -1, fmt.Errorf("unknown log level: %s", l)
}

func initConfig() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		stdLog.Fatalf(err.Error())
	}
}

func initLogger() {
	lvl, err := resolveLogLevel(cfg.Log.Level)
	if err != nil {
		fmt.Printf("error while resolving log level: %v\n", err.Error())
		os.Exit(1)
	}

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.Level(lvl))

	logger.WithFields(logrus.Fields{
		"logLevel":       lvl,
		"logLevelString": cfg.Log.Level,
	}).Info("logger initialized")
}

func initHTTPHandler(s *Service, routes ...router.Controller) error {
	httpRouter := router.NewHTTPRouter(routes...)

	s.server = &http.Server{
		Addr:    cfg.Server.HttpAddress,
		Handler: httpRouter,
	}

	logger.WithFields(logrus.Fields{
		"transport": "http",
		"addr":      cfg.Server.HttpAddress,
	}).Info("http server initialized")
	return s.server.ListenAndServe()
}

func initRoutes(db *sqlx.DB) []router.Controller {
	healthCheck := health.NewController()

	userRepository := user.NewRepository(user.WithDb(db))
	userService := user.NewService(user.WithRepository(userRepository))
	userService = user.NewServiceLoggingMiddleware(logger.WithField("service", "userService").Logger)(userService)
	users := user.NewController(
		user.WithService(userService),
	)

	return []router.Controller{healthCheck, users}
}

func initDb() *sqlx.DB {
	db, err := sqlx.Connect("postgres", cfg.Postgres.Uri)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Panicln("cannot connect to db")
	}

	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConnections)

	return db
}
