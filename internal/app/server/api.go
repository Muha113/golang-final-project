package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Muha113/golang-final-project/internal/app/model"
	"github.com/Muha113/golang-final-project/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

type Server struct {
	config *Config
	router *mux.Router
	repo   *repository.UsersRepositoryInMemory
	logger *logrus.Logger
}

func NewServer(config *Config) *Server {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DbConnectionString))
	return &Server{
		config: config,
		router: mux.NewRouter(),
		logger: logrus.New(),
		repo:   repository.NewUsersRepositoryInMemory(client.Database(config.DbName)),
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

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var user model.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			s.logger.Error(err)
			return
		}
		user.UserPasswordHash = hashPasswd(user.UserPasswordHash)
		err = s.repo.SaveUser(user)
		if err != nil {
			s.logger.Error(err)
		}
		json.NewEncoder(w).Encode(user)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser model.User
		err := json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := s.repo.GetUserByEmail(reqUser.UserEmail)
		if err != nil || user.UserPasswordHash != hashPasswd(reqUser.UserPasswordHash) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expTime := time.Now().Add(5 * time.Minute)
		claims := &model.JWTClaims{
			UserEmail: user.UserEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(s.config.Jwt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expTime,
		})
	}
}

func (s *Server) handleSubscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, status, err := s.isAuthorized(r)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(status)
			return
		}
		var reqUser model.User
		err = json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := s.repo.GetUserByEmail(claims.UserEmail)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tmp := user.UserFollowing
		user.UserFollowing = append(user.UserFollowing, reqUser.UserName)
		err = s.repo.UpdateUser(user)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			user.UserFollowing = tmp
			return
		}
	}
}

func (s *Server) handleTweetsPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) handleTweetsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) isAuthorized(r *http.Request) (*model.JWTClaims, int, error) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, http.StatusUnauthorized, err
		}
		return nil, http.StatusBadRequest, err
	}

	tknStr := c.Value
	claims := &model.JWTClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return s.config.Jwt, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, http.StatusUnauthorized, err
		}
		return nil, http.StatusBadRequest, err
	}
	if !tkn.Valid {
		return nil, http.StatusUnauthorized, fmt.Errorf("Error: %s", "token is not valid")
	}
	return claims, http.StatusAccepted, nil
}

func hashPasswd(passwd string) string {
	hasher := md5.New()
	hasher.Write([]byte(passwd))
	return hex.EncodeToString(hasher.Sum(nil))
}
