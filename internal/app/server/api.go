package server

import (
	"context"
	"log"
	"net/http"

	"github.com/Muha113/golang-final-project/pkg/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

//Server : represents server functionality
type Server struct {
	config *Config
	router *mux.Router
	repo   *repository.UsersRepositoryInMemory
	logger *logrus.Logger
}

//NewServer : builds new server using configuration 'config'
func NewServer(config *Config) *Server {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DbConnectionString))
	return &Server{
		config: config,
		router: mux.NewRouter(),
		logger: logrus.New(),
		repo:   repository.NewUsersRepositoryInMemory(client.Database(config.DbName)),
	}
}

//Start : sets routes and logger and starts listen and server
func (s *Server) Start() error {
	s.setRoutes()
	s.setLogger()
	srvString := s.config.Host + ":" + s.config.Port
	return http.ListenAndServe(srvString, s.router)
}

func (s *Server) setRoutes() {
	s.router.HandleFunc("/register", s.handleRegister()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/subscribe", s.handleSubscribe()).Methods("POST")
	s.router.HandleFunc("/tweets", s.handleTweetsPost()).Methods("POST")
	s.router.HandleFunc("/tweets", s.handleTweetsGet()).Methods("GET")
}

func (s *Server) setLogger() {
	lvl, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	s.logger.SetLevel(lvl)
}
