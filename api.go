package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func respOk(ctx *fasthttp.RequestCtx, text string) {
	ctx.SetStatusCode(200)
	ctx.SetBodyString(text)
	ctx.SetContentType("application/json")
}

func respErr(ctx *fasthttp.RequestCtx, code int, text string) {
	ctx.SetStatusCode(code)
	ctx.SetBodyString("{\"error\":\"" + text + "\"}")
	ctx.SetContentType("application/json")
}

type Request struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func login() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		var req Request

		err := json.Unmarshal(ctx.Request.Body(), &req)
		if err != nil {
			respErr(ctx, http.StatusBadRequest, err.Error())
			return
		}

		user, ok := users[req.UserName]
		if !ok || user.Password != req.Password {
			respErr(ctx, http.StatusUnauthorized, "wrong username or password")
			return
		}

		tokenString, err := getTokenString(user)
		if err != nil {
			respErr(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		respOk(ctx, "{\"token\":\""+tokenString+"\"}")
	}
}

func getUsers() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		tk, errorCode, err := validateToken(ctx)
		if err != nil {
			respErr(ctx, errorCode, err.Error())
			return
		}

		if tk.Role == "Admin" {
			usersJson, err := json.Marshal(users)
			if err != nil {
				respErr(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			respOk(ctx, string(usersJson))
		} else {
			usersJson, err := json.Marshal(users[tk.UserName])
			if err != nil {
				respErr(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			respOk(ctx, string(usersJson))
		}
	}
}
