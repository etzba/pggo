package server

import (
	"fmt"
	"net/http"

	"github.com/etzba/pggo/dat"
	"github.com/etzba/pggo/pkg/env"
	"github.com/etzba/pggo/pkg/logger"
	"github.com/etzba/pggo/wire"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Logger     *logger.Log
	HTTPServer *http.Server
	Mux        *http.ServeMux
	Respoder   wire.Responder
	Database   *dat.Context
	shipper    Shipper
}

func New(logger *logger.Log, address string) *Server {
	responder := wire.Respond{
		Logger: logger,
	}
	server := &Server{
		Logger:   logger,
		Respoder: responder,
	}
	logger.Info("configuring prometheus shipper and register metrics")
	server.shipper = NewShipper(logger)
	server.shipper.Register()
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
	datCtx, err := dat.GetDatabaseContext()
	if err != nil {
		s.Logger.Error("failed to connect to database", err)
		return err
	}

	s.Database = datCtx
	s.Logger.Info(fmt.Sprintf("Connected to database %s:%d", env.PostgresHost, env.PostgresPort))
	if err := s.Database.InitDB(); err != nil {
		s.Logger.Error("failed to run db migrations", err)
		return err
	}

	s.Logger.Info("Database migration completed")
	s.Logger.Info("Start server in port 8080")
	if err := s.HTTPServer.ListenAndServe(); err != nil {
		s.Logger.Error("cannot run http server - listen and serve", err)
	}
	return nil
}

func (s *Server) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = notFound
	router.MethodNotAllowedHandler = methodNotAllowed
	router.HandleFunc("/locations", s.getLocations()).Methods("GET")
	router.HandleFunc("/locations/{id}", s.getLocationById()).Methods("GET")
	router.HandleFunc("/locations", s.addLocation()).Methods("POST")
	router.HandleFunc("/locations/{id}", s.updateLocation()).Methods("PUT")
	router.HandleFunc("/locations/{id}", s.deleteLocationById()).Methods("DELETE")
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		})).Methods("GET")
	return router
}

var notFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found")) //nolint:errcheck
})

var methodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed")) //nolint:errcheck
})
