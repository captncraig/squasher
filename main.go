package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/captncraig/ghauth"
	"github.com/captncraig/temple"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
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

	locked := r.Group("/", auth.RequireAuth())
	locked.GET("/repo/:owner/:repo", repo)
	locked.GET("/repo/:owner/:repo/:pr", pull)

	r.Run(":8765")
}
func render(c *gin.Context, name string, data interface{}) {
	c.Header("Content-Type", "text/html")
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

func getIntQuery(ctx *gin.Context, name string, def int) int {
	q := ctx.Query(name)
	i, err := strconv.Atoi(q)
	if err != nil {
		return def
	}
	return i
}

func home(ctx *gin.Context) {
	u := ghauth.User(ctx)
	data := gin.H{"User": u}
	if u != nil {
		opts := &github.RepositoryListOptions{}
		opts.PerPage = 20
		opts.Page = getIntQuery(ctx, "page", 0)
		opts.Sort = "pushed"
		repos, res, err := u.Client().Repositories.List("", opts)
		if err != nil {
			ctx.Error(err)
			return
		}
		data["Result"] = res
		data["Repos"] = repos
	}
	render(ctx, "home", data)
}

func repo(ctx *gin.Context) {
	u := ghauth.User(ctx)
	owner, repo := ctx.Param("owner"), ctx.Param("repo")
	data := gin.H{"User": u, "Owner": owner, "Repo": repo}
	opts := &github.PullRequestListOptions{}
	opts.PerPage = 20
	opts.Page = getIntQuery(ctx, "page", 0)
	pulls, res, err := u.Client().PullRequests.List(owner, repo, opts)
	if err != nil {
		ctx.Error(err)
		return
	}
	data["Pulls"] = pulls
	data["Result"] = res
	render(ctx, "repo", data)
}

func pull(ctx *gin.Context) {
	owner, repo, num := ctx.Param("owner"), ctx.Param("repo"), ctx.Param("pr")
	number, _ := strconv.Atoi(num)
	u := ghauth.User(ctx)
	data := gin.H{"User": u, "Owner": owner, "Repo": repo, "Number": number}

	pr, res, err := u.Client().PullRequests.Get(owner, repo, number)
	if err != nil {
		ctx.Error(err)
		return
	}
	data["Pull"] = pr
	data["Result"] = res
	data["Merged"] = *pr.Merged
	data["Conflicts"] = !*pr.Mergeable
	render(ctx, "pr", data)
}
