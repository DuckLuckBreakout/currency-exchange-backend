package usecase

import (
	"mime/multipart"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/repository_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/usecase_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/user"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/hasher"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/s3_storage"

	"github.com/pkg/errors"
)

type UserUseCase struct {
	userRepo user.Repository
}

func NewUseCase(userRepo user.Repository) user.UseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) SignUp(signupUser *models.SignupUserRequest) (uint64, error) {
	hashOfPassword, err := hasher.GenerateHashFromPassword(signupUser.Password)
	if err != nil {
		return 0, errors.Wrap(usecase_errors.UcCanNotGenerateHash, err.Error())
	}

	userId, err := u.userRepo.InsertUser(&models.UserData{
		Email:    signupUser.Email,
		Password: hashOfPassword,
		Username: signupUser.Username,
	})
	if err != nil {
		return 0, errors.Wrap(usecase_errors.UcEmailAlreadyExists, err.Error())
	}

	return userId, nil
}

func (u *UserUseCase) LogIn(loginUser *models.LoginUserRequest) (uint64, error) {
	authData, err := u.userRepo.SelectUserByEmailOrUsername(loginUser.EmailOrUsername)
	if err != nil {
		return 0, errors.Wrap(usecase_errors.UcCanNotFindUser, err.Error())
	}

	if ok := hasher.CompareHashAndPassword(authData.Password, loginUser.Password); !ok {
		return 0, usecase_errors.UcPasswordsNotMatch
	}

	return authData.Id, nil
}

func (u *UserUseCase) GetSelfProfile(userId uint64) (*models.UserData, error) {
	userData, err := u.userRepo.SelectUserById(userId)
	if err != nil {
		return nil, errors.Wrap(usecase_errors.UcCanNotFindUserId, err.Error())
	}
	//userData.Avatar.Url = s3_storage.PathToFile(userData.Avatar.Url, s3_storage.Avatar)
	//userData.Background.Url = s3_storage.PathToFile(userData.Background.Url, s3_storage.Background)

	return userData, nil
}

func (u *UserUseCase) DeleteSelfProfile(userId uint64) error {
	err := u.userRepo.DeleteSelfProfile(userId)
	if err != nil {
		return errors.Wrap(usecase_errors.UcCanNotDeleteUser, err.Error())
	}

	return nil
}

func (u *UserUseCase) CheckUniqEmail(email string) (bool, error) {
	_, err := u.userRepo.SelectUserByEmailOrUsername(email)
	switch errors.Cause(err) {
	case nil:
		return false, nil
	case repository_errors.PgCanNotFind:
		return true, nil
	default:
		return false, errors.Wrap(usecase_errors.UcInternalError, err.Error())
	}
}

func (u *UserUseCase) UpdateSelfUsername(userId uint64, login string) error {
	err := u.userRepo.UpdateSelfUsername(userId, login)
	if err != nil {
		return errors.Wrap(usecase_errors.UcCanNotUpdateUserLogin, err.Error())
	}

	return nil
}

func (u *UserUseCase) UpdateSelfAvatar(userId uint64,
	fileHeader *multipart.FileHeader) (*models.UserAvatar, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Upload new user avatar
	fileType := fileHeader.Header.Get("Content-Type")
	fileName, err := s3_storage.UploadMultipartFile(&src, fileType, s3_storage.Avatar, userId, &configer.AppConfig.S3)
	if err != nil {
		return nil, usecase_errors.UcCanNotUploadAvatar
	}

	// Delete old user avatar
	userData, err := u.userRepo.SelectUserById(userId)
	if err == nil && userData.Avatar.Url != "" {
		if err = s3_storage.DeleteFileByKey(userData.Avatar.Url, &configer.AppConfig.S3); err != nil {
			return nil, usecase_errors.UcCanNotDeleteAvatar
		}
	}

	err = u.userRepo.UpdateSelfAvatar(userId, fileName)
	if err != nil {
		return nil, usecase_errors.UcCanNotUpdateAvatar
	}

	return &models.UserAvatar{
		//Url: s3_storage.PathToFile(fileName, s3_storage.Avatar),
		Url: fileName,
	}, nil
}
