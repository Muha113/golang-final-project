package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Muha113/golang-final-project/internal/app/model"
	"github.com/dgrijalva/jwt-go"
)

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
		return []byte(s.config.Jwt), nil
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
