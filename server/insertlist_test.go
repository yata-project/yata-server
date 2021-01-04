package server

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/model"
	"github.com/TheYeung1/yata-server/server/request"
	"github.com/stretchr/testify/assert"
)

func TestInsertListInput_Validate(t *testing.T) {
	tests := map[string]struct {
		input InsertListInput
		err   error
	}{
		"validate-input": {
			input: InsertListInput{ListID: "ID", Title: "Title"},
		},
		"list-id-empty": {
			input: InsertListInput{ListID: ""},
			err:   errors.New("ListID cannot be empty"),
		},
		"list-id-too-long": {
			input: InsertListInput{ListID: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque porta eros erat. Curabitur nam."},
			err:   errors.New("ListID length cannot exceed 100 characters"),
		},
		"list-id-with-trailing-space": {
			input: InsertListInput{ListID: "ID "},
			err:   errors.New("ListID cannot be prefixed or suffixed with spaces"),
		},
		"list-id-with-leading-space": {
			input: InsertListInput{ListID: " ID"},
			err:   errors.New("ListID cannot be prefixed or suffixed with spaces"),
		},
		"title-empty": {
			input: InsertListInput{ListID: "ID"},
			err:   errors.New("Title cannot be empty"),
		},
		"title-too-long": {
			input: InsertListInput{ListID: "ID", Title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque porta eros erat. Curabitur nam."},
			err:   errors.New("Title length cannot exceed 100 characters"),
		},
		"title-id-with-trailing-space": {
			input: InsertListInput{ListID: "ID", Title: "Title "},
			err:   errors.New("Title cannot be prefixed or suffixed with spaces"),
		},
		"title-id-with-leading-space": {
			input: InsertListInput{ListID: "ID", Title: " Title"},
			err:   errors.New("Title cannot be prefixed or suffixed with spaces"),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(fmt.Sprintf(name), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.err, test.input.Validate())
		})
	}
}

func TestServer_InsertList(t *testing.T) {
	tests := map[string]struct {
		input      string
		insertList func(*testing.T) func(id model.UserID, list model.YataList) error
		outCode    int
		outBody    string
	}{
		"happy-path": {
			input: "{\"ListID\":\"ID\",\"Title\":\"Title\"}",
			insertList: func(t *testing.T) func(id model.UserID, list model.YataList) error {
				return func(id model.UserID, list model.YataList) error {
					assert.Equal(t, "userID", string(id))
					assert.Equal(t, model.YataList{
						UserID: "userID",
						ListID: "ID",
						Title:  "Title",
					}, list)
					return nil
				}
			},
			outCode: http.StatusCreated,
			outBody: "{\"ListID\":\"ID\"}\n",
		},
		"insertion-error": {
			input: "{\"ListID\":\"ID\",\"Title\":\"Title\"}",
			insertList: func(t *testing.T) func(id model.UserID, list model.YataList) error {
				return func(id model.UserID, list model.YataList) error {
					return errors.New("boom")
				}
			},
			outCode: http.StatusInternalServerError,
			outBody: "{\"Code\":\"InternalServerError\"}\n",
		},
		"list-already-exists": {
			input: "{\"ListID\":\"ID\",\"Title\":\"Title\"}",
			insertList: func(t *testing.T) func(id model.UserID, list model.YataList) error {
				return func(id model.UserID, list model.YataList) error {
					return database.ListExistsError{}
				}
			},
			outCode: http.StatusConflict,
			outBody: "{\"Code\":\"ListExists\",\"Message\":\"List already exists\"}\n",
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(fmt.Sprintf(name), func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest( /* Method */ "", "https://does.not/matter", bytes.NewBufferString(test.input))

			srvr := Server{Ydb: mockYdb{MockInsertList: test.insertList(t)}}

			srvr.InsertList(rec, req.WithContext(request.WithUserID(req.Context(), "userID")))

			assert.Equal(t, test.outCode, rec.Code)
			assert.Equal(t, test.outBody, rec.Body.String())
		})
	}
}

var _ database.YataDatabase = mockYdb{}

type mockYdb struct {
	MockInsertList func(id model.UserID, list model.YataList) error
}

func (m mockYdb) GetList(id model.UserID, id2 model.ListID) (model.YataList, error) {
	panic("implement me")
}

func (m mockYdb) GetLists(id model.UserID) ([]model.YataList, error) {
	panic("implement me")
}

func (m mockYdb) InsertList(id model.UserID, list model.YataList) error {
	return m.MockInsertList(id, list)
}

func (m mockYdb) GetAllItems(id model.UserID) ([]model.YataItem, error) {
	panic("implement me")
}

func (m mockYdb) GetListItems(id model.UserID, id2 model.ListID) ([]model.YataItem, error) {
	panic("implement me")
}

func (m mockYdb) InsertItem(item model.YataItem) error {
	panic("implement me")
}
