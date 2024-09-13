package storage

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

var _ SessionRepository = &RedisSessionRepository{}

type SessionRepository interface {
	AddSession(userID string) (string, error)
	GetUserBySession(jwtToken string) (uuid.UUID, error)
}

type RedisSessionRepository struct {
	redisClient *redis.Client
	jwtSecret   []byte
}

func (r *RedisSessionRepository) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(r.jwtSecret)
}

func (r *RedisSessionRepository) GetUserBySession(jwtToken string) (uuid.UUID, error) {
	ctx := context.Background()
	userIDStr, err := r.redisClient.Get(ctx, jwtToken).Result()
	if err == redis.Nil {
		return uuid.UUID{}, errors.New("invalid token")
	} else if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid user ID format")
	}
	return userID, nil
}

func (r *RedisSessionRepository) AddSession(userID string) (string, error) {
	jwtToken, err := r.generateToken(userID)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	err = r.redisClient.Set(ctx, jwtToken, userID, time.Minute*5).Err()
	return jwtToken, err
}
