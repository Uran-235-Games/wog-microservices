package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "battle_service/internal/app"
    "battle_service/internal/config"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    application := app.NewApp(cfg)

    // Канал для сигналов ОС
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    go application.MustRunHTTPServer()
    go application.ServergRPC.MustRun()

    <-stop
    log.Println("Shutting down gracefully...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    go func() {
        if err := application.Server.Shutdown(ctx); err != nil {
            log.Fatalf("Server forced to shutdown: %v", err)
        }
    }()

    go application.ServergRPC.Stop()
}
