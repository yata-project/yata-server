package database

import (
	"github.com/TheYeung1/yata-server/model"
)

type YataDatabase interface {
	GetList(model.UserID, model.ListID) (model.YataList, error)
	GetLists(model.UserID) ([]model.YataList, error)
	InsertList(model.UserID, model.YataList) error
}
