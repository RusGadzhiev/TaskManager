package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/RusGadzhiev/TaskManager/internal/config"
	"github.com/RusGadzhiev/TaskManager/internal/service"

	"github.com/redis/go-redis/v9"
)

var (
	ErrPingRedis = errors.New("error of ping redis db")
)

type SessionsRepoRedis struct {
	DB *redis.Client
}

func NewSessionsRepoRedis(ctx context.Context, cfg *config.RedisDb) *SessionsRepoRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(cfg.Host + ":" + cfg.Port),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrPingRedis)
	}

	return &SessionsRepoRedis{
		DB: rdb,
	}
}

func (repo *SessionsRepoRedis) GetUser(ctx context.Context, cookieVal string) (string, error) {
	val, err := repo.DB.Get(ctx, cookieVal).Result()
	if err == redis.Nil {
		return "", service.ErrNoUserBySession
	} else if err != nil {
		return "", fmt.Errorf("get redis error: %w", err)
	}
	return val, nil
}

func (repo *SessionsRepoRedis) Add(ctx context.Context, cookieVal string, username string, dur time.Duration) error {
	_, err := repo.DB.Set(ctx, cookieVal, username, dur).Result()
	if err != nil {
		return fmt.Errorf("insert redis error: %w", err)
	}
	return nil
}

func (repo *SessionsRepoRedis) Delete(ctx context.Context, cookieVal string) error {
	_, err := repo.DB.Del(ctx, cookieVal).Result()
	if err != nil {
		return fmt.Errorf("delete redis error: %w", err)
	}
	return nil
}