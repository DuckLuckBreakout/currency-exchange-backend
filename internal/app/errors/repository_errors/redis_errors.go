package repository_errors

import "github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors"

var (
	RsCanNotDelete = errors.Error{
		Message: "Can't delete data",
	}
	RsCanNotGetByKey = errors.Error{
		Message: "Can't get data by key",
	}
	RsCanNotSet = errors.Error{
		Message: "Can't set data",
	}
	RsCanNotParse = errors.Error{
		Message: "Can't parse data",
	}
)
