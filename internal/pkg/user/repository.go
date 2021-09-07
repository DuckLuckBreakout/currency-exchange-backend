package user

import "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"

type Repository interface {
	InsertUser(user *models.UserData) (uint64, error)
	SelectUserByEmailOrUsername(emailOrUsername string) (*models.UserData, error)
	SelectUserById(userId uint64) (*models.UserData, error)
	DeleteSelfProfile(userId uint64) error
	UpdateSelfUsername(userId uint64, username string) error
	UpdateSelfAvatar(userId uint64, avatarUrl string) error
}
