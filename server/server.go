package server

import (
	"net/http"

	"github.com/TheYeung1/yata-server/config"
	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/middleware"
	"github.com/TheYeung1/yata-server/middleware/auth"
	"github.com/google/uuid"
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
	r.Use(middleware.RequestLogger(func() string {
		u, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("failed to generate a uuid: %v", err) // If we get here bad things have happened (ie, we cannot read from crypto/rand's random reader).
		}
		return u.String()
	}))
	r.Use(mux.CORSMethodMiddleware(r),
		middleware.CORSAccessControlHeadersMiddleware,
		auth.CognitoJwtAuthMiddleware{Cfg: s.CognitoCfg}.Execute)
	r.HandleFunc("/items", s.GetAllItems).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/lists", s.GetLists).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/lists", s.InsertList).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/lists/{listID}/", s.GetList).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/lists/{listID}/items", s.GetListItems).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/lists/{listID}/items", s.InsertListItem).Methods(http.MethodPut, http.MethodOptions)
	log.Fatal(http.ListenAndServe(addr, r))
}
