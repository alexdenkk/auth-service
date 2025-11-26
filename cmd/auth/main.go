package main

import (
	"alexdenkk/auth-service/internal/domain"
	"alexdenkk/auth-service/internal/handler"
	"alexdenkk/auth-service/internal/repository"
	"alexdenkk/auth-service/pkg/config"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	service := fx.New(
		fx.Provide(
			// Configuration
			config.NewConfigFromEnv,

			// Logger instance
			newLogger,

			// Echo instance
			newEcho,

			// Repository instance
			func(cfg *config.Config, logger *zap.Logger) domain.UserRepository {
				return repository.NewUserRepository(cfg.DBConfig, logger)
			},

			// Service instance
			func(cfg *config.Config, logger *zap.Logger, repo domain.UserRepository) domain.Service {
				return domain.NewService(cfg.JwtConfig, logger, repo)
			},

			// Handler instance
			handler.NewHandler,
		),

		fx.Invoke(registerHooks),
	)

	service.Run()
}

// Logger instance
func newLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	return logger, nil
}

// Echo instance
func newEcho(cfg *config.Config) *echo.Echo {
	// New echo server
	e := echo.New()

	// Timeouts
	e.Server.ReadTimeout = cfg.HttpConfig.ReadTimeout
	e.Server.WriteTimeout = cfg.HttpConfig.WriteTimeout
	e.Server.IdleTimeout = cfg.HttpConfig.IdleTimeout

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

// Registering hooks for fx lifecycle
func registerHooks(
	lc fx.Lifecycle,
	e *echo.Echo,
	handler *handler.Handler,
	cfg *config.Config,
	logger *zap.Logger,
) {
	handler.RegisterEndpoints(e)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Lauching web server in goroutine
			go func() {
				logger.Info("Starting HTTP server", zap.String("address", cfg.HttpConfig.Host))

				if err := e.Start(cfg.HttpConfig.Host); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")

			return e.Shutdown(ctx)
		},
	})
}
