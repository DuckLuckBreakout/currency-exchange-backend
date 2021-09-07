package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/repository_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session"

	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	connAccessDB  *redis.Client
	connRefreshDB *redis.Client
}

func NewSessionRedisRepository(connAccessDB, connRefreshDB *redis.Client) session.Repository {
	return &SessionRepository{
		connAccessDB:  connAccessDB,
		connRefreshDB: connRefreshDB,
	}
}

func (r *SessionRepository) InsertToken(userId uint64, token *models.TokenDetails) error {
	data := fmt.Sprintf("%d", userId)
	err := r.connAccessDB.Set(context.TODO(), token.AccessDetails.Uuid, data, models.AccessTokenExpires).Err()
	if err != nil {
		return repository_errors.RsCanNotSet
	}

	err = r.connRefreshDB.Set(context.TODO(), token.RefreshDetails.Uuid, data, models.RefreshTokenExpires).Err()
	if err != nil {
		return repository_errors.RsCanNotSet
	}

	return nil
}

func (r *SessionRepository) SelectUserIdByAccessToken(Uuid string) (uint64, error) {
	data, err := r.connAccessDB.Get(context.TODO(), Uuid).Bytes()
	if err != nil {
		return 0, repository_errors.RsCanNotGetByKey
	}

	userId, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, repository_errors.RsCanNotParse
	}

	return userId, nil
}

func (r *SessionRepository) SelectUserIdByRefreshToken(Uuid string) (uint64, error) {
	data, err := r.connRefreshDB.Get(context.TODO(), Uuid).Bytes()
	if err != nil {
		return 0, repository_errors.RsCanNotGetByKey
	}

	userId, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, repository_errors.RsCanNotParse
	}

	return userId, nil
}

func (r *SessionRepository) DeleteAccessToken(Uuid string) error {
	err := r.connAccessDB.Del(context.TODO(), Uuid).Err()
	if err != nil {
		return repository_errors.RsCanNotGetByKey
	}

	return nil
}

func (r *SessionRepository) DeleteRefreshToken(Uuid string) error {
	err := r.connRefreshDB.Del(context.TODO(), Uuid).Err()
	if err != nil {
		return repository_errors.RsCanNotDelete
	}

	return nil
}
