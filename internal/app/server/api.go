package server

import (
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type Server struct {
	config *Config
	router *mux.Router
	//repo *repository.
	logger *logrus.Logger
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
}

func (s *Server) Start() error {
	s.setRoutes()
	s.setLogger()
	srvString := s.config.Host + ":" + s.config.Port
	return http.ListenAndServe(srvString, s.router)
}

func (s *Server) setRoutes() {
	s.router.HandleFunc("/register", s.handleRegister()).Methods("GET")
}

func (s *Server) setLogger() {
	lvl, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	s.logger.SetLevel(lvl)
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("my first log")
		w.Write([]byte("HIIIIIIIIIIIIIII"))
	}
}
