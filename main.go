package main

import (
	"log"
	"os"

	"github.com/captncraig/ghauth"
	"github.com/captncraig/temple"
	"github.com/gin-gonic/gin"
)

//go:generate govendor add +external
//go:generate govendor update +ven
//go:generate templeGen -pkg=main -var=myTemplates -o=templates.go -dir=templates

var templateManager temple.TemplateStore

func main() {
	r := gin.Default()

	var err error
	templateManager, err = temple.New(os.Getenv("TEMPLE_DEV") != "", myTemplates, "templates")
	if err != nil {
		log.Fatal(err)
	}

	conf := &ghauth.Conf{
		ClientId:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user", "repo"},
		CookieName:   "ghuser",
		CookieSecret: os.Getenv("COOKIE_SECRET"),
	}
	auth := ghauth.New(conf)
	auth.RegisterRoutes("/login", "/callback", "/logout", r)

	open := r.Group("/", auth.OpenHandler())
	open.GET("/", func(ctx *gin.Context) {
		u := ghauth.User(ctx)
		ctx.Header("Content-Type", "text/html")
		if err := templateManager.Execute(ctx.Writer, gin.H{"User": u}, "home.tpl"); err != nil {
			ctx.Error(err)
		}
	})
	r.Run(":8765")
}
