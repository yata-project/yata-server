package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/TheYeung1/yata-server/model"
	"github.com/TheYeung1/yata-server/server/request"
	"github.com/gorilla/mux"
)

type InsertListItemInput struct {
	ItemID  string
	Content string
}

// Validate returns an error if the input does not pass validation.
func (input *InsertListItemInput) Validate() error {
	if len(input.ItemID) == 0 {
		return errors.New("ItemID cannot be empty")
	}
	if len(input.ItemID) > 100 {
		return errors.New("ItemID length cannot exceed 100 characters")
	}
	if len(input.ItemID) != len(strings.TrimSpace(input.ItemID)) {
		return errors.New("ItemID cannot be prefixed or suffixed with spaces")
	}
	if len(input.Content) == 0 {
		return errors.New("Content cannot be empty")
	}
	if len(input.Content) > 100 {
		return errors.New("Content length cannot exceed 100 characters")
	}
	if len(input.Content) != len(strings.TrimSpace(input.Content)) {
		return errors.New("Content cannot be prefixed or suffixed with spaces")
	}
	return nil
}

type InsertListItemOutput struct {
	ItemID string
}

func (s *Server) InsertListItem(w http.ResponseWriter, r *http.Request) {
	log := request.Logger(r.Context())
	uid, ok := request.UserID(r.Context())
	if !ok {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w, r)
		return
	}

	var input InsertListItemInput
	if err := bindJSON(r.Body, &input); err != nil {
		log.WithError(err).Info("failed to bind input")
		renderBadRequest(w, r, "malformed input")
		return
	}
	log.WithField("input", input).Debug("input bound")

	v := mux.Vars(r)
	listID := model.ListID(v["listID"])
	if err := validateListID(listID); err != nil {
		log.WithError(err).Info("failed to validate input")
		renderBadRequest(w, r, err.Error())
		return
	}
	if err := input.Validate(); err != nil {
		log.WithError(err).Info("failed to normalize and validate input")
		renderBadRequest(w, r, err.Error())
		return
	}

	yi := model.YataItem{
		UserID:  uid,
		ListID:  model.ListID(v["listID"]),
		ItemID:  model.ItemID(input.ItemID),
		Content: input.Content,
	}
	log.WithField("item", yi).Debug("inserting item")
	if err := s.Ydb.InsertItem(yi); err != nil {
		log.WithError(err).Error("failed to insert item")
		renderInternalServerError(w, r)
		return
	}

	out := InsertListItemOutput{ItemID: input.ItemID}
	log.WithField("output", out).Debug("item inserted")
	renderJSON(w, r, http.StatusCreated, out)
}
