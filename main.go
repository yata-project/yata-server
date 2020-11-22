package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/model"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/gorilla/mux"
)

type server struct {
	ydb database.YataDatabase
}

type InsertListInput struct {
	ListID string
	Title  string
}

type InsertListOutput struct {
	ListID string
}

type InsertListItemInput struct {
	ItemID  string
	Content string
}

func writeInternalErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Sorry! Something went wrong"))
}

func (s *server) GetList(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	userID := model.UserID(r.Header.Get("User"))
	listID := model.ListID(v["listID"])

	yl, err := s.ydb.GetList(userID, listID)
	if err != nil {
		if lnf, ok := err.(database.ListNotFoundError); ok {
			log.Infoln(lnf)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("List not found"))
			return
		}
		log.Errorln(err)
		writeInternalErrorResponse(w)
		return
	}

	res, err := json.Marshal(yl)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
	}
	_, err = w.Write(res)
	if err != nil {
		log.Errorln(err)
	}
}

func (s *server) GetLists(w http.ResponseWriter, r *http.Request) {
	userID := model.UserID(r.Header.Get("User"))

	yl, err := s.ydb.GetLists(userID)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
	}

	res, err := json.Marshal(yl)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println(err)
	}
}

func (s *server) InsertList(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
		return
	}

	var in InsertListInput
	err = json.Unmarshal(b, &in)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
		return
	}

	uid, ok := r.Header["User"]
	if !ok {
		log.Errorln("UserId not provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("UserId is missing"))
		return
	}

	// TODO: assert input lengths
	yl := model.YataList{
		UserID: model.UserID(uid[0]),
		ListID: model.ListID(in.ListID),
		Title:  in.Title,
	}

	// insert list to db here
	err = s.ydb.InsertList(yl.UserID, yl)
	if err != nil {
		if errnf, ok := err.(database.ListExistsError); ok {
			log.Warnln(errnf)
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("List already exists"))
			return
		}
		log.Errorln(err)
		writeInternalErrorResponse(w)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte{})
	if err != nil {
		log.Errorln(err)
	}
}

func (s *server) InsertListItem(w http.ResponseWriter, r *http.Request) {
	//TODO: add validation to inputs

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
		return
	}

	var in InsertListItemInput
	err = json.Unmarshal(b, &in)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
		return
	}

	uid, ok := r.Header["User"]
	if !ok {
		log.Errorln("UserId not provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("UserId is missing"))
		return
	}
	v := mux.Vars(r)

	// TODO: assert input lengths
	item := model.YataItem{
		UserID:  model.UserID(uid[0]),
		ListID:  model.ListID(v["listID"]),
		ItemID:  model.ItemID(in.ItemID),
		Content: in.Content,
	}

	// insert list to db here
	err = s.ydb.InsertItem(item)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte{})
	if err != nil {
		log.Errorln(err)
	}
}

func (s *server) GetListItems(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	userID := model.UserID(r.Header.Get("User"))
	listID := model.ListID(v["listID"])

	items, err := s.ydb.GetListItems(userID, listID)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
	}

	res, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println(err)
	}
}

func (s *server) GetAllItems(w http.ResponseWriter, r *http.Request) {
	userID := model.UserID(r.Header.Get("User"))

	items, err := s.ydb.GetAllItems(userID)
	if err != nil {
		log.Errorln(err)
		writeInternalErrorResponse(w)
	}

	res, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println(err)
	}
}

func (s *server) start() {
	r := mux.NewRouter()
	r.HandleFunc("/items", s.GetAllItems).Methods(http.MethodGet)
	r.HandleFunc("/lists", s.GetLists).Methods(http.MethodGet)
	r.HandleFunc("/lists", s.InsertList).Methods(http.MethodPut)
	r.HandleFunc("/lists/{listID}/", s.GetList).Methods(http.MethodGet)
	r.HandleFunc("/lists/{listID}/items", s.GetListItems).Methods(http.MethodGet)
	r.HandleFunc("/lists/{listID}/items", s.InsertListItem).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8888", r))
}

func main() {
	log.SetReportCaller(true)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "yata"),
	})

	if err != nil {
		log.Fatal(err)
	}

	yataDynamo := &database.DynamoDbYataDatabase{
		Dynamo: dynamodb.New(sess),
	}

	s := server{
		ydb: yataDynamo,
	}
	s.start()
}
