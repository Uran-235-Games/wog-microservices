package db

import (
	"fmt"
	"context"
	"encoding/json"
	"wog-server/domain"

	"github.com/redis/go-redis/v9"
)

type GameRepo struct {
	client *redis.Client
}

func (r *GameRepo) SaveRequestRedis(rq domain.BattleRequest) (string, error) {
	data, err := json.Marshal(rq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal BattleRequest: %w", err)
	}

	key := rq.Id
	err = r.client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save BattleRequest to Redis: %w", err)
	}

	return rq.Id, nil
}

func (r *GameRepo) GetRequestRedis(reqId string) (*domain.BattleRequest, error) {
	var obj domain.BattleRequest
	key := reqId

	data, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return &obj, fmt.Errorf("failed to get BattleRequest from Redis: %w", err)
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return &obj, fmt.Errorf("failed to unmarshal BattleRequest data: %w", err)
	}

	return &obj, nil
}

func (r *GameRepo) SaveGameRedis(obj domain.BattleObj) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal BattleObj: %w", err)
	}

	key := obj.Id
	err = r.client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save BattleObj to Redis: %w", err)
	}

	return obj.Id, nil
}

func (r *GameRepo) GetGameRedis(gameID string) (*domain.BattleObj, error) {
	var obj domain.BattleObj
	key := gameID

	data, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return &obj, fmt.Errorf("failed to get game from Redis: %w", err)
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return &obj, fmt.Errorf("failed to unmarshal game data: %w", err)
	}

	return &obj, nil
}
