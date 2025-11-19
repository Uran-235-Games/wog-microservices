package app

import (
	"log/slog"
	"wog-server/internal/service/auth"
	"wog-server/internal/service/battle"
	"wog-server/internal/db"
	"wog-server/internal/logger"

	"github.com/joho/godotenv"
	"os"
	"log"
)

type Config struct {
	SrvPort	string
	GRPCSrvPort	string
	SocketIOPort string
	LogMode	string
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func (con *Config) Init() *Config {
	return &Config{
		SrvPort:      getEnv("HTTP_PORT", "1488"),
		GRPCSrvPort:  getEnv("GRPCPORT", "1489"),
		SocketIOPort: getEnv("SOCKETIOPORT", "1487"),
		LogMode:      getEnv("MODE", "DEV"),
	}
}

type App struct {
	Config *Config
	Log *slog.Logger

	Repo struct {
		User db.IUserRepo
		Game db.IGameRepo
	}

	Service struct {
		Auth *auth.AuthService
		Battle *battle.BattleService
	}
}

// загружает окружение, настраивает логгер
func Init() *App {
	a := App{}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("load env error: %s", err.Error())
	}

	a.Config = a.Config.Init()

	logger.Init(a.Config.LogMode)
	a.Log = logger.Log
	a.Log.Info("Logger запущен в режиме: "+os.Getenv("MODE"))

	return &a
}