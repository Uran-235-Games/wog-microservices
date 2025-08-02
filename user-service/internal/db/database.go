package db

import (
	"log"
	"log/slog"
	"os"
	"wog-server/domain"
	"wog-server/internal/logger"
)

type Database interface {
	Connect(string) error
	Disconnect() error
}

type IUserRepo interface {
	Get(uid string) (*domain.DBUser, error)
	SaveRedis(uid string, obj domain.UserRedisData) error
	GetRedis(uid string) (*domain.UserRedisData, error)

	Insert(domain.User) (string, error)
	Find(string, any) (*domain.DBUser, error)
	IsFieldValueExists(field string, value any) (bool, error)
}

type IGameRepo interface {
	SaveRequestRedis(domain.BattleRequest) (id string, err error)
	GetRequestRedis(id string) (*domain.BattleRequest, error)

	SaveGameRedis(domain.BattleObj) (id string, err error)
	GetGameRedis(gameID string) (*domain.BattleObj, error)
	// Save(domain.BattleObj)
	// Get(gameID string) (domain.BattleObj, error)
}

type DB struct {
	mongoDB Database
	redisDB Database

	UserRepo IUserRepo
	GameRepo IGameRepo
}

func InitDB() *DB{
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not defined in .env")
	}

	var mongo MongoDB
	if err := mongo.Connect(mongoURI); err != nil {
		panic("Ошибка коннекта к mongoDB")
	}
	logger.Log.Info("MongoDB connected")

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		log.Fatal("REDIS_URI not defined in .env")
	}
	logger.Log.Info("RedisDB connected")
	
	var redis RedisDB
	redis.Connect(redisURI)

	db := DB{mongoDB: &mongo, redisDB: &redis}
	db.UserRepo = &UserRepo{coll: mongo.Database.Collection("users"), redis: redis.client}
	db.GameRepo = &GameRepo{client: redis.client}
	return &db
}

func (s *DB) Disconnect() {
	if err := s.mongoDB.Disconnect(); err != nil {
		logger.Log.Error("Ошибка дисконнекта mongoDB", slog.String("err", err.Error()))
	}
	if err := s.redisDB.Disconnect(); err != nil {
		logger.Log.Error("Ошибка дисконнекта redisDB", slog.String("err", err.Error()))
	}
}