package main

import (
	"context"
	"errors"
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
	"wog-server/internal/socketio"
	"wog-server/internal/grpc"
	// "wog-server/internal/api/middleware"
	"wog-server/internal/api/routes"
	"wog-server/internal/service/auth"
	// "wog-server/internal/service/battle"

	"github.com/gin-gonic/gin"
	// socketio "github.com/doquangtan/socket.io/v4"
)

var (
	ErrEnvLoad = errors.New("Incorrect environment variable")
)

func main() {
	a := app.Init()
	database := db.InitDB()

	// repos
	a.Repo.Game = database.GameRepo
	a.Repo.User = database.UserRepo

	// libs
	jwtLib := jwt.NewJWTLib()

	// services
	authService := auth.NewAuthService(a.Repo.User, jwtLib)
	// battleService := battle.NewBattleService(a.Repo.Game)

	authHandler := handler.NewAuthHandler(authService)

	// Setup Server
	r := gin.Default()
	// r.Use(middleware.RequestLogger(log))

	SetRoutes(r, authHandler)

	// SocketIO
	ioHandler := socketio.Setup(a)

	r.GET("/socket.io/", gin.WrapH(ioHandler))
	r.Run(":"+a.Config.SocketIOPort)
	// socket.SetupSocketIO(r, jwtLib, a.Repo.User, a.Repo.Game, battleService)

	srv := &http.Server{
		Addr:    ":"+a.Config.SrvPort,
		Handler: r,
	}

	// gRPC App
	gRPCApp := grpc_app.New(logger.Log, authService, a.Config.GRPCSrvPort)
	gRPCApp.MustRun()

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Log.Info("Server is running on http://localhost:"+a.Config.SrvPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Starting server error", slog.String("err", err.Error()))
		}
	}()

	// 5. Канал для OS-сигналов (Ctrl+C, SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gRPCApp.Stop()
	database.Disconnect()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", slog.String("err", err.Error()))
	}
}

func SetRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	v1 := r.Group("/api/v1")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	routers.RegisterAuthRoutes(v1, authHandler)
}