package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"strings"
)

type Token struct {
	UserID   int
	UserName string
	Role     string
	jwt.StandardClaims
}

func getTokenString(user User) (string, error) {
	tk := &Token{UserID: user.ID, Role: user.Role, UserName: user.UserName}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(ctx *fasthttp.RequestCtx) (*Token, int, error) {

	tokenHeader := string(ctx.Request.Header.Peek("Authorization"))

	if strings.TrimSpace(tokenHeader) == "" {
		return nil, http.StatusUnauthorized, errors.New("missing auth token")
	}

	//`Bearer {token-body}`
	splitHeader := strings.Split(tokenHeader, " ")
	if len(splitHeader) != 2 {
		return nil, http.StatusForbidden, errors.New("invalid/malformed authentication token")
	}

	tk := &Token{}
	token, err := jwt.ParseWithClaims(splitHeader[1], tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, http.StatusForbidden, errors.New("invalid/malformed authentication token")
	}

	if !token.Valid {
		return nil, http.StatusForbidden, errors.New("token is not valid")
	}

	return tk, 0, nil
}
