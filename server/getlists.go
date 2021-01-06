package server

import (
	"net/http"

	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/model"
	"github.com/TheYeung1/yata-server/server/request"
	"github.com/gorilla/mux"
)

type GetListOutput struct {
	List model.YataList
}

func (s *Server) GetList(w http.ResponseWriter, r *http.Request) {
	log := request.Logger(r.Context())
	uid, ok := request.UserID(r.Context())
	if !ok {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w, r)
		return
	}
	log.WithField("userID", uid).Debug("get list called")

	v := mux.Vars(r)
	listID := model.ListID(v["listID"])
	if err := validateListID(listID); err != nil {
		log.WithError(err).Info("failed to validate input")
		renderBadRequest(w, r, err.Error())
		return
	}

	yl, err := s.Ydb.GetList(uid, listID)
	if err != nil {
		if errnf, ok := err.(database.ListNotFoundError); ok {
			log.WithError(errnf).Info("list not found")
			renderJSON(w, r, http.StatusNotFound, responseError{Code: "ListDoesNotExist", Message: "List does not exist"})
			return
		}
		log.WithError(err).Error("failed to get list")
		renderInternalServerError(w, r)
		return
	}

	out := GetListOutput{List: yl}
	log.WithField("output", out).Debug("list retrieved")
	renderJSON(w, r, http.StatusOK, out)
}

type GetListsOutput struct {
	Lists []model.YataList
}

func (s *Server) GetLists(w http.ResponseWriter, r *http.Request) {
	log := request.Logger(r.Context())
	uid, ok := request.UserID(r.Context())
	if !ok {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w, r)
		return
	}
	log.WithField("userID", uid).Debug("get lists called")

	yl, err := s.Ydb.GetLists(uid)
	if err != nil {
		log.WithError(err).Error("failed to get lists")
		renderInternalServerError(w, r)
	}

	out := GetListsOutput{Lists: yl}
	log.WithField("output", out).Debug("lists retrieved")
	renderJSON(w, r, http.StatusOK, out)
}
