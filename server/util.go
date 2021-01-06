package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/TheYeung1/yata-server/model"
	"github.com/TheYeung1/yata-server/server/request"
)

// getUserIDFromContext returns the userID stored on the request context.
// A non-nil error will be returned if the userID cannot be found.
// Deprecated: use request.UserID instead.
func getUserIDFromContext(r *http.Request) (model.UserID, error) {
	val, ok := request.UserID(r.Context())
	if ok {
		return val, nil
	}
	return "", errors.New("userID context is not a string value")
}

// bindJSON decodes r into v.
// v must be a pointer.
// Q: Why such a small function?
// A: We might want to put limits on how much we read in the future.
// A: If we decide to ever support other input formats (XML?) this can become a generic "Bind" function.
// In an ideal world we might want to look at the incoming request's "Content-Type" header.
func bindJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

type responseError struct {
	Code    string
	Message string `json:",omitempty"`
}

// renderJSON writes the response code, sets the content type for JSON, and encodes v as JSON to w.
// In an ideal world we might want to look at the incoming request's "Accept" header.
func renderJSON(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json") // FYI: https://stackoverflow.com/questions/477816/what-is-the-correct-json-content-type
	if err := json.NewEncoder(w).Encode(v); err != nil {
		request.Logger(r.Context()).WithError(err).Warn("failed to render json")
	}
}

func renderInternalServerError(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, r, http.StatusInternalServerError, responseError{Code: "InternalServerError"})
}

func renderBadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	renderJSON(w, r, http.StatusBadRequest, responseError{Code: "BadRequest", Message: msg})
}

func validateListID(id model.ListID) error {
	if len(id) == 0 {
		return errors.New("ListID cannot be empty")
	}
	if len(id) > 100 {
		return errors.New("ListID length cannot exceed 100 characters")
	}
	if len(id) != len(strings.TrimSpace(string(id))) {
		return errors.New("ListID cannot be prefixed or suffixed with spaces")
	}
	return nil
}
