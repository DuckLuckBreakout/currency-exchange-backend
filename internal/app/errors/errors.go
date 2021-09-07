package errors

import (
	"fmt"
)

type Error struct {
	Message string `json:"message"`
}

func (err Error) Error() string {
	return fmt.Sprintf("error: happened %s", err.Message)
}
