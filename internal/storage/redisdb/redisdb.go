package redisdb

import (
	"bot/internal/config"
	"bot/pkg/logging"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
)

type RedisDB struct {
	redisClient *redis.Client
	log         *logging.Logger
	cfg         *config.Config
}

func New(redisClient *redis.Client, log *logging.Logger, cfg *config.Config) *RedisDB {
	return &RedisDB{
		redisClient: redisClient,
		log:         log,
		cfg:         cfg,
	}
}

func (r *RedisDB) SetState(telegramID int64, stateName string, stateData *map[string]interface{}) error {
	stateJSON, err := json.Marshal(stateData)
	if err != nil {
		r.log.Error("failed to marshal state data", zap.Error(err))
		return err
	}
	err = r.redisClient.HSet(context.TODO(), strconv.FormatInt(telegramID, 10), stateName, stateJSON).Err()
	if err != nil {
		r.log.Error("failed to set state", zap.Error(err))
		return err
	}

	return nil
}

func (r *RedisDB) GetState(telegramID int64, stateName string) (map[string]interface{}, error) {
	stateJSON, err := r.redisClient.HGet(context.TODO(), strconv.FormatInt(telegramID, 10), stateName).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		r.log.Error("failed to get state", zap.Error(err))
		return nil, err
	}

	var stateData map[string]interface{}
	err = json.Unmarshal([]byte(stateJSON), &stateData)
	if err != nil {
		r.log.Error("failed to unmarshal state data", zap.Error(err))
		return nil, err
	}

	return stateData, nil
}

func (r *RedisDB) UpdateState(telegramID int64, stateName string, fieldName string, fieldValue interface{}) error {
	currentState, err := r.GetState(telegramID, stateName)
	if err != nil {
		r.log.Error("failed to get state", zap.Error(err))
		return err
	}

	currentState[fieldName] = fieldValue

	err = r.SetState(telegramID, stateName, &currentState)
	if err != nil {
		r.log.Error("failed to update state", zap.Error(err))
		return err
	}

	return nil
}

func (r *RedisDB) ClearState(telegramID int64, stateName string) error {
	err := r.redisClient.HDel(context.Background(), strconv.FormatInt(telegramID, 10), stateName).Err()
	if err != nil {
		r.log.Error("failed to clear state", zap.Error(err))
		return err
	}

	return nil
}
