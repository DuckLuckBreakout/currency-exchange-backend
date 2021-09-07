package handler_errors

import "github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors"

var (
	HttpSuccess = errors.Error{
		Message: "OK",
	}

	HttpEmailAlreadyExists = errors.Error{
		Message: "User email already exists",
	}
	HttpLoginAlreadyExists = errors.Error{
		Message: "User login already exists",
	}
	HttpIncorrectRequestBody = errors.Error{
		Message: "Body of request is incorrect",
	}
	HttpBadUserCredentials = errors.Error{
		Message: "Bad user credentials",
	}
	HttpRottenAccessToken = errors.Error{
		Message: "Access token is rotten",
	}
	HttpProfileNotConfigured = errors.Error{
		Message: "User profile not configured",
	}
	HttpWrongCode = errors.Error{
		Message: "Wrong code",
	}

	HttpInternalServerError = errors.Error{
		Message: "Internal server error",
	}
)
