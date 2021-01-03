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
	tests := []struct {
		input InsertListInput
		err   error
	}{
		{
			input: InsertListInput{ListID: "ID", Title: "Title"},
		},
		{
			input: InsertListInput{ListID: ""},
			err:   errors.New("ListID cannot be empty"),
		},
		{
			input: InsertListInput{ListID: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque porta eros erat. Curabitur nam."},
			err:   errors.New("ListID length cannot exceed 100 characters"),
		},
		{
			input: InsertListInput{ListID: "ID "},
			err:   errors.New("ListID cannot be prefixed or suffixed with spaces"),
		},
		{
			input: InsertListInput{ListID: " ID"},
			err:   errors.New("ListID cannot be prefixed or suffixed with spaces"),
		},
		{
			input: InsertListInput{ListID: "ID"},
			err:   errors.New("Title cannot be empty"),
		},
		{
			input: InsertListInput{ListID: "ID", Title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque porta eros erat. Curabitur nam."},
			err:   errors.New("Title length cannot exceed 100 characters"),
		},
		{
			input: InsertListInput{ListID: "ID", Title: "Title "},
			err:   errors.New("Title cannot be prefixed or suffixed with spaces"),
		},
		{
			input: InsertListInput{ListID: "ID", Title: " Title"},
			err:   errors.New("Title cannot be prefixed or suffixed with spaces"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%+v", test.input), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.err, test.input.Validate())
		})
	}
}

func TestServer_InsertList(t *testing.T) {
	tests := []struct {
		input      string
		insertList func(*testing.T) func(id model.UserID, list model.YataList) error
		outCode    int
		outBody    string
	}{
		{
			input: "{\"ListID\":\"ID\",\"Title\":\"Title\"}",
			insertList: func(t *testing.T) func(id model.UserID, list model.YataList) error {
				return func(id model.UserID, list model.YataList) error {
					assert.Equal(t, model.UserID("u"), id)
					assert.Equal(t, model.YataList{
						UserID: "u",
						ListID: "ID",
						Title:  "Title",
					}, list)
					return nil
				}
			},
			outCode: http.StatusCreated,
			outBody: "{\"ListID\":\"ID\"}\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%+v", test.input), func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest( /* Method */ "", "https://does.not/matter", bytes.NewBufferString(test.input))

			srvr := Server{Ydb: MockTdb{MockInsertList: test.insertList(t)}}

			srvr.InsertList(rec, req.WithContext(request.WithUserID(req.Context(), "userID")))

			assert.Equal(t, test.outCode, rec.Code)
			assert.Equal(t, test.outBody, rec.Body.String())
		})
	}
}

var _ database.YataDatabase = MockTdb{}

type MockTdb struct {
	MockInsertList func(id model.UserID, list model.YataList) error
}

func (m MockTdb) GetList(id model.UserID, id2 model.ListID) (model.YataList, error) {
	panic("implement me")
}

func (m MockTdb) GetLists(id model.UserID) ([]model.YataList, error) {
	panic("implement me")
}

func (m MockTdb) InsertList(id model.UserID, list model.YataList) error {
	return m.MockInsertList(id, list)
}

func (m MockTdb) GetAllItems(id model.UserID) ([]model.YataItem, error) {
	panic("implement me")
}

func (m MockTdb) GetListItems(id model.UserID, id2 model.ListID) ([]model.YataItem, error) {
	panic("implement me")
}

func (m MockTdb) InsertItem(item model.YataItem) error {
	panic("implement me")
}

func TestNewInsertListHandler(t *testing.T) {
	tests := []struct {
		input      string
		insertList func(*testing.T) func(id model.UserID, list model.YataList) error
		outCode    int
		outBody    string
	}{
		{
			input: "{\"ListID\":\"ID\",\"Title\":\"Title\"}",
			insertList: func(t *testing.T) func(id model.UserID, list model.YataList) error {
				return func(id model.UserID, list model.YataList) error {
					assert.Equal(t, model.UserID("u"), id)
					assert.Equal(t, model.YataList{
						UserID: "u",
						ListID: "ID",
						Title:  "Title",
					}, list)
					return nil
				}
			},
			outCode: http.StatusCreated,
			outBody: "{\"ListID\":\"ID\"}\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%+v", test.input), func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest( /* Method */ "", "https://does.not/matter", bytes.NewBufferString(test.input))
			handler := NewInsertListHandler(test.insertList(t))

			handler(rec, req.WithContext(request.WithUserID(req.Context(), "userID")))

			assert.Equal(t, test.outCode, rec.Code)
			assert.Equal(t, test.outBody, rec.Body.String())
		})
	}
}
