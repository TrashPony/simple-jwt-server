package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"os"
)

func init() {
	_ = godotenv.Load()
}

var users = map[string]User{
	"User_1": {
		ID:       1,
		UserName: "User_1",
		Password: "qwe",
		Role:     "Respondent",
	},
	"User_2": {
		ID:       2,
		UserName: "User_2",
		Password: "xcv",
		Role:     "Respondent",
	},
	"User_3": {
		ID:       3,
		UserName: "User_3",
		Password: "arwr",
		Role:     "Admin",
	},
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func main() {

	newRouter := router.New()
	newRouter.POST("/login", login())
	newRouter.GET("/get_users", getUsers())
	newRouter.ServeFiles("/{filepath:*}", "./static")

	fmt.Println("http server start: " + os.Getenv("SERVER_URL"))
	err := fasthttp.ListenAndServe(os.Getenv("SERVER_URL"), func(ctx *fasthttp.RequestCtx) {

		originHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "HEAD,GET,POST,PUT,DELETE,OPTIONS",
		}

		for h, v := range originHeaders {
			ctx.Response.Header.Set(h, v)
		}
		newRouter.Handler(ctx)
	})

	if err != nil {
		panic(err)
	}
}
