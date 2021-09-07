package usecase

import (
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/usecase_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/jwt_token"
)

type SessionUseCase struct {
	sessionRepo session.Repository
}

func NewUseCase(sessionRepo session.Repository) session.UseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo,
	}
}

func (u *SessionUseCase) GetUserIdByAccessToken(Uuid string) (uint64, error) {
	userId, err := u.sessionRepo.SelectUserIdByAccessToken(Uuid)
	if err != nil {
		return 0, usecase_errors.UcCanNotFindUserId
	}

	return userId, nil
}

func (u *SessionUseCase) CreateNewSession(userId uint64) (*models.Token, error) {
	token, err := jwt_token.CreateJwtToken(&configer.AppConfig.Secret)
	if err != nil {
		return nil, usecase_errors.UcCanNotCreateToken
	}

	if err = u.sessionRepo.InsertToken(userId, token); err != nil {
		return nil, usecase_errors.UcCanNotInsertToken
	}

	return &token.Token, nil
}

func (u *SessionUseCase) DestroySession(Uuid string) error {
	if err := u.sessionRepo.DeleteAccessToken(Uuid); err != nil {
		return usecase_errors.UcCanNotDeleteToken
	}

	if err := u.sessionRepo.DeleteRefreshToken(Uuid); err != nil {
		return usecase_errors.UcCanNotDeleteToken
	}

	return nil
}

func (u *SessionUseCase) RefreshSession(Uuid string) (*models.Token, error) {
	// Check access token
	if _, err := u.sessionRepo.SelectUserIdByAccessToken(Uuid); err == nil {
		return nil, usecase_errors.UcAccessTokenNotRotten
	}

	// Destroy refresh token
	userId, err := u.sessionRepo.SelectUserIdByRefreshToken(Uuid)
	if err != nil {
		return nil, usecase_errors.UcCanNotFindUserId
	}

	if err = u.sessionRepo.DeleteRefreshToken(Uuid); err != nil {
		return nil, usecase_errors.UcCanNotDeleteToken
	}

	// Create and save new user token
	token, err := jwt_token.CreateJwtToken(&configer.AppConfig.Secret)
	if err != nil {
		return nil, usecase_errors.UcCanNotCreateToken
	}

	if err = u.sessionRepo.InsertToken(userId, token); err != nil {
		return nil, usecase_errors.UcCanNotInsertToken
	}

	return &token.Token, nil
}
