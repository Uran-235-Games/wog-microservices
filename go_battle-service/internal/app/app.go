package app

import (
    "log"
    "net/http"

    "battle_service/internal/config"
    "battle_service/internal/grpc"
)

type App struct {
    Config *config.Config
    Server *http.Server
    ServergRPC *grpc_app.App
}

func NewApp(cfg *config.Config) *App {
    mux := http.NewServeMux()

    mux.HandleFunc("/health", healthCheckHandler)

    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: mux,
    }

    grpc_app := grpc_app.New("1480")

    return &App{
        Config: cfg,
        Server: srv,
        ServergRPC: grpc_app,
    }
}

func (a *App) MustRunHTTPServer() {
	log.Printf("Starting server on port %s...\n", a.Config.Port)
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen: %v\n", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}