package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/repository_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/user"
)

type UserRepository struct {
	dbConn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) user.Repository {
	return &UserRepository{
		dbConn: conn,
	}
}

func (r *UserRepository) InsertUser(user *models.UserData) (uint64, error) {
	var userId uint64
	err := r.dbConn.Get(
		&userId,
		"INSERT INTO users(email, password, username, avatar) "+
			"VALUES ($1, $2, $3, $4) RETURNING id",
		user.Email,
		user.Password,
		user.Username,
		user.Avatar.Url,
	)
	if err != nil {
		return 0, repository_errors.PgCanNotInsert
	}

	return userId, nil
}

func (r *UserRepository) SelectUserByEmailOrUsername(emailOrUsername string) (*models.UserData, error) {
	row := r.dbConn.QueryRow(
		"SELECT id, email, username, password, avatar "+
			"FROM users "+
			"WHERE email = $1 or username = $1",
		emailOrUsername,
	)

	userData := &models.UserData{}
	err := row.Scan(
		&userData.Id,
		&userData.Email,
		&userData.Username,
		&userData.Password,
		&userData.Avatar.Url,
	)

	switch err {
	case nil:
		return userData, nil
	case sql.ErrNoRows:
		return nil, repository_errors.PgCanNotFind
	default:
		return nil, repository_errors.PgInternalDbError
	}
}

func (r *UserRepository) SelectUserById(userId uint64) (*models.UserData, error) {
	row := r.dbConn.QueryRow(
		"SELECT id, email, username, password, avatar "+
			"FROM users "+
			"WHERE id = $1",
		userId,
	)

	userData := &models.UserData{}
	err := row.Scan(
		&userData.Id,
		&userData.Email,
		&userData.Username,
		&userData.Password,
		&userData.Avatar.Url,
	)

	switch err {
	case nil:
		return userData, nil
	case sql.ErrNoRows:
		return nil, repository_errors.PgCanNotFind
	default:
		return nil, repository_errors.PgInternalDbError
	}
}

func (r *UserRepository) DeleteSelfProfile(userId uint64) error {
	_, err := r.dbConn.Exec(
		"DELETE FROM users "+
			"WHERE id = $1",
		userId,
	)

	if err != nil {
		return repository_errors.PgCanNotDelete
	}

	return nil
}

func (r *UserRepository) UpdateSelfUsername(userId uint64, username string) error {
	_, err := r.dbConn.Exec(
		"UPDATE users "+
			"SET username = $2 "+
			"WHERE id = $1",
		userId,
		username,
	)

	if err != nil {
		return repository_errors.PgCanNotUpdate
	}

	return nil
}

func (r *UserRepository) UpdateSelfAvatar(userId uint64, avatarUrl string) error {
	_, err := r.dbConn.Exec(
		"UPDATE users "+
			"SET avatar = $2 "+
			"WHERE id = $1",
		userId,
		avatarUrl,
	)

	if err != nil {
		return repository_errors.PgCanNotUpdate
	}

	return nil
}
