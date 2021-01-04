package server

import (
	"net/http"

	"github.com/TheYeung1/yata-server/config"
	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/middleware/auth"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	CognitoCfg config.AwsCognitoUserPoolConfig
	Ydb        database.YataDatabase
}

func (s *Server) Start() {
	addr := ":8888"
	log.WithField("address", addr).Info("starting server")
	r := mux.NewRouter()
	r.Use(auth.CognitoJwtAuthMiddleware{Cfg: s.CognitoCfg}.Execute)
	r.HandleFunc("/items", s.GetAllItems).Methods(http.MethodGet)
	r.HandleFunc("/lists", s.GetLists).Methods(http.MethodGet)
	r.HandleFunc("/lists", s.InsertList).Methods(http.MethodPut)
	r.HandleFunc("/lists/{listID}/", s.GetList).Methods(http.MethodGet)
	r.HandleFunc("/lists/{listID}/items", s.GetListItems).Methods(http.MethodGet)
	r.HandleFunc("/lists/{listID}/items", s.InsertListItem).Methods(http.MethodPut)
	r.HandleFunc("/me", s.WhoAmI).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(addr, r))
}
