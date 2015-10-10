package main

import (
	"github.com/captncraig/ghauth"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()

	// first create the auth handler
	conf := &ghauth.Conf{
		ClientId:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Scopes:       []string{"user", "repo"},
		CookieName:   "ghuser",
		CookieSecret: os.Getenv("COOKIE_SECRET"),
	}
	ghauth.New(conf)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})
	r.Run(":8765")
}
