package main

import (
	"fmt"
	"log"
	"os"

	"github.com/captncraig/ghauth"
	"github.com/captncraig/temple"
	"github.com/gin-gonic/gin"
)

//go:generate govendor add +external
//go:generate govendor update +ven
//go:generate templeGen -pkg=main -var=myTemplates -o=templates.go -dir=templates
//go:generate esc -o static.go -prefix static static

var templateManager temple.TemplateStore

func main() {
	r := gin.Default()

	var err error
	useDev := os.Getenv("TEMPLE_DEV") != ""
	templateManager, err = temple.New(useDev, myTemplates, "templates")
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
	r.Use(renderError, auth.AuthCheck())
	r.StaticFS("/static", FS(useDev))
	r.GET("/", home)

	r.Run(":8765")
}
func render(c *gin.Context, name string, data interface{}) {
	if err := templateManager.Execute(c.Writer, data, name); err != nil {
		c.AbortWithError(500, err)
	}
}
func renderError(c *gin.Context) {
	c.Next()
	errs := c.Errors.Errors()
	fmt.Println(errs)
	if len(errs) > 0 {
		u := ghauth.User(c)
		render(c, "error", gin.H{"User": u, "Errors": errs})
	}
}
func home(ctx *gin.Context) {
	u := ghauth.User(ctx)
	ctx.Header("Content-Type", "text/html")
	data := gin.H{"User": u}
	render(ctx, "home", data)
}
