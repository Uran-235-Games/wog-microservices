package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"wog-server/domain"
	"wog-server/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	coll *mongo.Collection
	redis *redis.Client
}

func (r *UserRepo) Insert(user domain.User) (string, error) {
	mongoUser, err := r.coll.InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}
	
	id := mongoUser.InsertedID.(bson.ObjectID)
	return id.Hex(), nil
}

// Возвращает (nil, nil) если юзер не найден
func (r *UserRepo) Find(key string, value any) (*domain.DBUser, error) {
	filter := bson.M{key: value}
	var user domain.DBUser
	err := r.coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Возвращает (nil, nil) если юзер не найден
func (r *UserRepo) Get(uid string) (*domain.DBUser, error) {
	user, err := r.Find("_id", uid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		cerr := errors.New("Ошибка получения юзера из mongoDB")
		logger.Log.Error(cerr.Error(), slog.String("err", err.Error()))
		return nil, cerr
	}
	return user, nil
}

func (r *UserRepo) IsFieldValueExists(field string, value interface{}) (bool, error) {
	filter := bson.M{field: value}
	count, err := r.coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) SaveRedis(uid string, obj domain.UserRedisData) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("Failed to marshal user data: %w", err)
	}

	if err := r.redis.Set(context.Background(), uid, data, 0).Err(); err != nil {
		return fmt.Errorf("Failed to save user data to Redis: %w", err)
	}

	return nil
}

func (r *UserRepo) GetRedis(uid string) (*domain.UserRedisData, error) {
	data, err := r.redis.Get(context.Background(), uid).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var res domain.UserRedisData
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal user data from Redis: %w", err)
	}
	return &res, nil
}