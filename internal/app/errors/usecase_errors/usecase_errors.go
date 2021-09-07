package usecase_errors

import "github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors"

var (
	UcCanNotDeleteToken = errors.Error{
		Message: "Can't delete token",
	}
	UcCanNotFindUser = errors.Error{
		Message: "Can't find user by id",
	}
	UcCanNotInsertToken = errors.Error{
		Message: "Can't insert token in database",
	}
	UcCanNotInsertCode = errors.Error{
		Message: "Can't insert code in database",
	}
	UcCanNotCheckCode = errors.Error{
		Message: "Can't check code in database",
	}
	UcCanNotCreateToken = errors.Error{
		Message: "Can't create new token",
	}
	UcCanNotFindUserId = errors.Error{
		Message: "Can't create new token",
	}
	UcCanNotGenerateHash = errors.Error{
		Message: "Can't generate hash",
	}
	UcCanNotInsertUser = errors.Error{
		Message: "Can't insert user in database",
	}
	UcCanNotExchangeToken = errors.Error{
		Message: "Can't exchange token from auth code",
	}
	UcPasswordsNotMatch = errors.Error{
		Message: "Passwords don't match",
	}
	UcCanNotGetDataFromRequest = errors.Error{
		Message: "Can't get data from request",
	}
	UcCanNotDecode = errors.Error{
		Message: "Can't decode data",
	}
	UcCanNotCreateUniqLogin = errors.Error{
		Message: "Can't create uniq login",
	}
	UcEmailAlreadyExists = errors.Error{
		Message: "Email already exists in database",
	}
	UcLoginAlreadyExists = errors.Error{
		Message: "Login already exists in database",
	}
	UcCanNotDeleteUser = errors.Error{
		Message: "Can't delete user",
	}
	UcCanNotDeleteCode = errors.Error{
		Message: "Can't delete code",
	}
	UcCanNotSendCode = errors.Error{
		Message: "Can't send code",
	}
	UcAccessTokenNotRotten = errors.Error{
		Message: "Access token not rotten",
	}
	UcCanNotGenerateLogin = errors.Error{
		Message: "Can't generate login",
	}
	UcInternalError = errors.Error{
		Message: "Internal error",
	}
	UcCanNotUpdateUserLogin = errors.Error{
		Message: "Can't update user login",
	}
	UcCanNotCreateOAuthClient = errors.Error{
		Message: "Can't create client for oauth",
	}
	UcCanNotGetUserFromVk = errors.Error{
		Message: "Can't get user from vk",
	}
	UcCanNotUpdateUserProfile = errors.Error{
		Message: "Can't update user profile",
	}
	UcCanNotUploadAvatar = errors.Error{
		Message: "Can't update user avatar",
	}
	UcCanNotDeleteAvatar = errors.Error{
		Message: "Can't delete user avatar",
	}
	UcCanNotUpdateAvatar = errors.Error{
		Message: "Can't update user avatar",
	}
	UcCanNotUploadBackground = errors.Error{
		Message: "Can't update user background",
	}
	UcCanNotDeleteBackground = errors.Error{
		Message: "Can't delete user background",
	}
	UcCanNotUpdateBackground = errors.Error{
		Message: "Can't update user background",
	}
)
