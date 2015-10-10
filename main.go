package main

import (
	"fmt"
	"github.com/captncraig/ghauth"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()

	// first create the auth handler
	conf := &ghauth.Conf{
		ClientId:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user", "repo"},
		CookieName:   "ghuser",
		CookieSecret: os.Getenv("COOKIE_SECRET"),
	}
	auth := ghauth.New(conf)

	// register oauth routes
	auth.RegisterRoutes("/login", "/callback", "/logout", r)

	r.GET("/", func(ctx *gin.Context) {
		u := ghauth.User(ctx)
		fmt.Println(u)
		msg := "Not logged in"
		if u != nil {
			msg = "hello " + u.Login
		}
		ctx.String(200, msg)
	})
	r.Run(":8765")
}
