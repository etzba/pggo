package server

import (
	"net/http"

	"github.com/etzba/pggo/pkg/logger"
	"github.com/etzba/pggo/wire"
	"github.com/gorilla/mux"
)

type Server struct {
	Logger     *logger.Log
	HTTPServer *http.Server
	Mux        *http.ServeMux
	Respoder   wire.Responder
}

func New(logger *logger.Log, address string) *Server {
	responder := wire.Respond{
		Logger: logger,
	}
	server := &Server{
		Logger:   logger,
		Respoder: responder,
	}
	router := server.getRouter()
	server.Mux = http.NewServeMux()
	server.Mux.Handle("/", router)
	server.HTTPServer = &http.Server{
		Addr:    address,
		Handler: router,
	}
	return server
}

func (s *Server) Run() error {
	s.Logger.Info("Start server in port 8080")
	if err := s.HTTPServer.ListenAndServe(); err != nil {
		s.Logger.Error("cannot run http server - listen and serve", err)
	}
	return nil
}

func (s *Server) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/locations", s.getLocations()).Methods("GET")
	router.HandleFunc("/locations/{id}", s.getLocationById()).Methods("GET")
	return router
}
