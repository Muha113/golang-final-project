package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Muha113/golang-final-project/internal/app/model"
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
	s.router.HandleFunc("/register", s.handleRegister()).Methods("POST")
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
		w.Header().Set("Content-Type", "application/json")
		var user model.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			s.logger.Error(err)
			return
		}
		hasher := md5.New()
		hasher.Write([]byte(user.UserPasswordHash))
		user.UserPasswordHash = hex.EncodeToString(hasher.Sum(nil))
		//json.NewEncoder(w).Encode(user)
	}
}
