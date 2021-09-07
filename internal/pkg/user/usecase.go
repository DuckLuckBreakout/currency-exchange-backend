package user

import (
	"mime/multipart"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
)

type UseCase interface {
	SignUp(signupUser *models.SignupUserRequest) (uint64, error)
	LogIn(loginUser *models.LoginUserRequest) (uint64, error)
	GetSelfProfile(userId uint64) (*models.UserData, error)
	DeleteSelfProfile(userId uint64) error
	CheckUniqEmail(email string) (bool, error)
	UpdateSelfUsername(userId uint64, login string) error
	UpdateSelfAvatar(userId uint64, avatar *multipart.FileHeader) (*models.UserAvatar, error)
}
