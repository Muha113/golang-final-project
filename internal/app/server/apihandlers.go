package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Muha113/golang-final-project/internal/app/model"
	"github.com/dgrijalva/jwt-go"
)

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.ID = s.repo.GetSize()
		type Response struct {
			ID        uint   `json:"id"`
			UserName  string `json:"username"`
			UserEmail string `json:"email"`
		}
		response := Response{
			ID:        user.ID,
			UserName:  user.UserName,
			UserEmail: user.UserEmail,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser model.User
		err := json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := s.repo.GetUserByEmail(reqUser.UserEmail)
		if err != nil || user.UserPasswordHash != hashPasswd(reqUser.UserPasswordHash) {
			s.logger.Error(err)
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
		tokenString, err := token.SignedString([]byte(s.config.Jwt))
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expTime,
		})
		type Response struct {
			Token string `json:"jwt"`
		}
		response := Response{
			Token: tokenString,
		}
		json.NewEncoder(w).Encode(response)
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
		for _, v := range user.UserFollowing {
			if v == reqUser.UserName {
				s.logger.Error("Error: Duplicate following username")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		tmp := user.UserFollowing
		user.UserFollowing = append(user.UserFollowing, reqUser.UserName)
		fmt.Println("Following: ", user.UserFollowing)
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
		claims, status, err := s.isAuthorized(r)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(status)
			return
		}
		var tweet model.Tweet
		err = json.NewDecoder(r.Body).Decode(&tweet)
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
		tweet.DateTime = time.Now()
		tweet.ID = uint(len(user.UserTweets)) + 1
		tweet.UserName = user.UserName
		tmp := user.UserTweets
		user.UserTweets = append(user.UserTweets, tweet)
		err = s.repo.UpdateUser(user)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			user.UserTweets = tmp
			return
		}
		type Respone struct {
			ID           uint   `json:"id"`
			TweetMessage string `json:"message"`
		}
		response := Respone{
			ID:           tweet.ID,
			TweetMessage: tweet.TweetMessage,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) handleTweetsGet() http.HandlerFunc { //sort by datetime
	return func(w http.ResponseWriter, r *http.Request) {
		claims, status, err := s.isAuthorized(r)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(status)
			return
		}
		user, err := s.repo.GetUserByEmail(claims.UserEmail)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//subscribtions := make([]model.Tweet, 0)
		type Response struct {
			Tweet    string    `json:"message"`
			DateTime time.Time `json:"datetime"`
		}
		response := []Response{}
		for _, v := range user.UserFollowing {
			sub, err := s.repo.GetUserByUserName(v)
			if err != nil {
				s.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, val := range sub.UserTweets {
				//subscribtions = append(subscribtions, val)
				response = append(response, Response{
					Tweet:    val.TweetMessage,
					DateTime: val.DateTime,
				})
			}
		}
		json.NewEncoder(w).Encode(response)
	}
}
