package repository_errors

import "github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors"

var (
	PgInternalDbError = errors.Error{
		Message: "Internal database error",
	}
	PgBadConnection = errors.Error{
		Message: "Bad connection",
	}
	PgCanNotFind = errors.Error{
		Message: "Can't find data",
	}
	PgCanNotInsert = errors.Error{
		Message: "Can't insert data",
	}
	PgCanNotDelete = errors.Error{
		Message: "Can't delete data",
	}
	PgCanNotUpdate = errors.Error{
		Message: "Can't update data",
	}
)
