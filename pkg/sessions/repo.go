package sessions

import (
	"HW4/internal/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrPingRedis = errors.New("error of ping redis db")
	ErrNoSession = errors.New("session with this id does not exist")
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

func (repo *SessionsRepoRedis) GetSession(ctx context.Context, id uint64) (*Session, error) {
	val, err := repo.DB.HGet(ctx, "session", string(rune(id))).Result()
	if err == redis.Nil {
		return nil, ErrNoSession
	} else if err != nil {
		return nil, fmt.Errorf("get redis error: %w", err)
	}

	session := &Session{}
	err = json.Unmarshal([]byte(val), session)
	if err != nil {
		return nil, err
	}
	return session, err
}

func (repo *SessionsRepoRedis) Add(ctx context.Context, session *Session) error {
	id := uuid.New().String()
	jsonValue, err := json.Marshal(*session)
	if err != nil {
		return err
	}
	_, err = repo.DB.HSet(ctx, "session", id, jsonValue).Result()
	if err != nil {
		return fmt.Errorf("insert redis error: %w", err)
	}
	return nil
}
