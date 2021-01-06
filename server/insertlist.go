package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/model"
	"github.com/TheYeung1/yata-server/server/request"
)

type InsertListInput struct {
	ListID string
	Title  string
}

// Validate returns an error if the input does not pass validation.
func (input *InsertListInput) Validate() error {
	if len(input.ListID) == 0 {
		return errors.New("ListID cannot be empty")
	}
	if len(input.ListID) > 100 {
		return errors.New("ListID length cannot exceed 100 characters")
	}
	if len(input.ListID) != len(strings.TrimSpace(input.ListID)) {
		return errors.New("ListID cannot be prefixed or suffixed with spaces")
	}
	if len(input.Title) == 0 {
		return errors.New("Title cannot be empty")
	}
	if len(input.Title) > 100 {
		return errors.New("Title length cannot exceed 100 characters")
	}
	if len(input.Title) != len(strings.TrimSpace(input.Title)) {
		return errors.New("Title cannot be prefixed or suffixed with spaces")
	}
	return nil
}

type InsertListOutput struct {
	ListID string
}

func (s *Server) InsertList(w http.ResponseWriter, r *http.Request) {
	log := request.Logger(r.Context())
	uid, ok := request.UserID(r.Context())
	if !ok {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w, r)
		return
	}
	log.WithField("userID", uid).Debug("insert list called")

	var input InsertListInput
	if err := bindJSON(r.Body, &input); err != nil {
		log.WithError(err).Info("failed to bind input")
		renderBadRequest(w, r, "malformed input")
		return
	}
	log.WithField("input", input).Debug("input bound")

	if err := input.Validate(); err != nil {
		log.WithError(err).Info("failed to normalize and validate input")
		renderBadRequest(w, r, err.Error())
		return
	}

	yl := model.YataList{
		UserID: uid,
		ListID: model.ListID(input.ListID),
		Title:  input.Title,
	}
	log.WithField("list", yl).Debug("inserting list")
	if err := s.Ydb.InsertList(yl.UserID, yl); err != nil {
		if errnf, ok := err.(database.ListExistsError); ok {
			log.WithError(errnf).Info("list not found")
			renderJSON(w, r, http.StatusConflict, responseError{Code: "ListExists", Message: "List already exists"})
			return
		}
		log.WithError(err).Error("failed to insert list")
		renderInternalServerError(w, r)
		return
	}

	out := InsertListOutput{ListID: input.ListID}
	log.WithField("output", out).Debug("list inserted")
	renderJSON(w, r, http.StatusCreated, out)
}
