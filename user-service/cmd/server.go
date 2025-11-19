package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	

	"wog-server/internal/app"
	"wog-server/internal/api/handler"
	"wog-server/internal/db"
	"wog-server/internal/lib/jwt"
	"wog-server/internal/logger"
	"wog-server/internal/service/auth"
	"wog-server/internal/grpc"
	"wog-server/internal/api/routes"

	"github.com/gin-gonic/gin"
)

var (
	ErrEnvLoad = errors.New("Incorrect environment variable")
)

func main() {
	a := app.Init()
	database := db.InitDB()

	// Repositories
	a.Repo.Game = database.GameRepo
	a.Repo.User = database.UserRepo

	// Libraries
	jwtLib := jwt.NewJWTLib()

	// Services
	authService := auth.NewAuthService(a.Repo.User, jwtLib)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup Gin router
	r := gin.Default()
	SetRoutes(r, authHandler)

	// HTTP server
	srv := &http.Server{
		Addr:    ":" + a.Config.SrvPort,
		Handler: r,
	}

	// gRPC server
	gRPCApp := grpc_app.New(logger.Log, a.Config.GRPCSrvPort, authService, jwtLib)

	// Channel for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run HTTP server in a goroutine
	go func() {
		logger.Log.Info("HTTP server is running on http://localhost:" + a.Config.SrvPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("HTTP server error", slog.String("err", fmt.Sprint(err)))
		}
	}()

	// Run gRPC server in a goroutine
	go func() {
		gRPCApp.MustRun();
	}()

	// Wait for OS signal
	<-quit
	logger.Log.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop gRPC and DB
	gRPCApp.Stop()
	database.Disconnect()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("HTTP server forced to shutdown", slog.String("err", fmt.Sprint(err)))
	}

	logger.Log.Info("Servers gracefully stopped")
}

// SetRoutes configures all API routes
func SetRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	v1 := r.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	routers.RegisterAuthRoutes(v1, authHandler)
}
