package database

import (
	"fmt"

	"github.com/TheYeung1/yata-server/model"
)

type ListExistsError struct {
	uid model.UserID
	lid model.ListID
}

func (e ListExistsError) Error() string {
	return fmt.Sprintf("List %s already exists for user %s", e.lid, e.uid)
}
