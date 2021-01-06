package server

import (
	"net/http"

	"github.com/TheYeung1/yata-server/model"
	"github.com/TheYeung1/yata-server/server/request"
	"github.com/gorilla/mux"
)

type GetAllItemsOutput struct {
	Items []model.YataItem
}

func (s *Server) GetAllItems(w http.ResponseWriter, r *http.Request) {
	log := request.Logger(r.Context())
	uid, ok := request.UserID(r.Context())
	if !ok {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w, r)
		return
	}
	log.WithField("userID", uid).Debug("get all items called")

	items, err := s.Ydb.GetAllItems(uid)
	if err != nil {
		log.WithError(err).Error("failed to get all items")
		renderInternalServerError(w, r)
		return
	}

	out := GetAllItemsOutput{Items: items}
	log.WithField("output", out).Debug("items retrieved")
	renderJSON(w, r, http.StatusOK, out)
}

type GetListItemsOutput struct {
	Items []model.YataItem
}

func (s *Server) GetListItems(w http.ResponseWriter, r *http.Request) {
	log := request.Logger(r.Context())
	uid, ok := request.UserID(r.Context())
	if !ok {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w, r)
		return
	}
	log.WithField("userID", uid).Debug("get list items called")

	v := mux.Vars(r)
	listID := model.ListID(v["listID"])
	if err := validateListID(listID); err != nil {
		log.WithError(err).Info("failed to validate input")
		renderBadRequest(w, r, err.Error())
		return
	}

	items, err := s.Ydb.GetListItems(uid, listID)
	if err != nil {
		log.WithError(err).Error("failed to get list items")
		renderInternalServerError(w, r)
	}

	out := GetListItemsOutput{Items: items}
	log.WithField("output", out).Debug("list items retrieved")
	renderJSON(w, r, http.StatusOK, out)
}
