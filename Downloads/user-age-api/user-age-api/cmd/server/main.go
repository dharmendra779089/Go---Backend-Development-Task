package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/example/user-age-api/config"
	"github.com/example/user-age-api/db/sqlc"
	"github.com/example/user-age-api/internal/handler"
	applogger "github.com/example/user-age-api/internal/logger"
	"github.com/example/user-age-api/internal/middleware"
	"github.com/example/user-age-api/internal/repository"
	"github.com/example/user-age-api/internal/routes"
	"github.com/example/user-age-api/internal/service"
)

func main() {
	// Initialize logger
	applogger.Init()
	defer applogger.Sync()
	logger := applogger.Log

	// Load configuration
	cfg := config.Load()

	// Connect to database
	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	// Wire up layers
	queries := sqlc.New(pool)
	repo := repository.NewUserRepository(queries)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc, logger)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logger))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Routes
	routes.RegisterUserRoutes(app, h)

	addr := ":" + cfg.AppPort
	logger.Info("starting server", zap.String("addr", addr))
	if err := app.Listen(addr); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
