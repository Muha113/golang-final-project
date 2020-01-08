package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config *Config
	router *mux.Router
	//repo *repository.
	//logger *logrus.Logger
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.setRoutes()
	srvString := s.config.Host + ":" + s.config.Port
	return http.ListenAndServe(srvString, s.router)
}

func (s *Server) setRoutes() {
	s.router.HandleFunc("/register", s.handleRegister()).Methods("GET")
}

// func (s *Server) setLogger() error {

// }

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HIIIIIIIIIIIIIII"))
	}
}
