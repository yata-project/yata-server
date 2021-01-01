package database

import (
	"fmt"

	"github.com/TheYeung1/yata-server/model"
)

type ListNotFoundError struct {
	uid model.UserID
	lid model.ListID
}

func (e ListNotFoundError) Error() string {
	return fmt.Sprintf("list not found. UserID: %q, ListID: %q", e.uid, e.lid)
}

type ListExistsError struct {
	uid model.UserID
	lid model.ListID
}

func (e ListExistsError) Error() string {
	return fmt.Sprintf("list %q already exists for user %q", e.lid, e.uid)
}
