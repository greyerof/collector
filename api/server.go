package api

import (
	"net/http"
	"time"

	"github.com/redhat-best-practices-for-k8s/collector/storage"
	"github.com/redhat-best-practices-for-k8s/collector/util"
	"github.com/sirupsen/logrus"
)

type Server struct {
	database    *storage.MySQLStorage
	objectStore *storage.S3Storage
	server      *http.Server
}

func NewServer(listenAddr string, db *storage.MySQLStorage, objectStore *storage.S3Storage, rTimeout, wTimeout time.Duration) *Server {
	return &Server{
		database:    db,
		objectStore: objectStore,
		server: &http.Server{
			Addr:         listenAddr,
			ReadTimeout:  rTimeout,
			WriteTimeout: wTimeout,
		},
	}
}

func (s *Server) Start() {
	logrus.Info("Starting server")
	http.HandleFunc("/", s.handler)
	//nolint:errcheck
	s.server.ListenAndServe()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	printServerUpMessage(w)

	switch r.Method {
	case http.MethodGet:
		ResultsHandler(w, r, s.database)
	case http.MethodPost:
		ParserHandler(w, r, s.database)
	default:
		util.WriteMsg(w, util.InvalidRequestErr)
		logrus.Errorf(util.InvalidRequestErr)
	}
}

func printServerUpMessage(w http.ResponseWriter) {
	logrus.Info(util.ServerIsUpMsg)
	_, writeErr := w.Write([]byte(util.ServerIsUpMsg + "\n"))
	if writeErr != nil {
		logrus.Errorf(util.WritingResponseErr, writeErr)
	}
}
