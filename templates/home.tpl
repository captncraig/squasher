{{if .User}}Hello, {{.User.Login}}. <a href="/logout">log out</a>{{else}}
Not logged in. <a href="/login">Log in.</a>

{{end}}