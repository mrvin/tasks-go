package mainform

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type formContent struct {
	Title string
	Name  string
	Hash  string
}

const defaultName = "Joe Bloggs"

var content = &formContent{
	Title: "Identidock",
	Name:  defaultName,
	Hash:  fmt.Sprintf("%x", sha256.Sum256([]byte(defaultName))),
}

var htmlForm = `<!DOCTYPE HTML>
<html>
	<head>
		<meta charset="utf-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<form action="/main" method="POST">
			<label>Name:</label>
			<input type="text" name="name" value="{{.Name}}">
			<input type="submit" value="submit">
		</form>
		<label>You look like a:</label>
		<img src="/monster?name={{.Hash}}">
	</body>
</html>`

var tempForm = template.Must(template.New("htmlForm").Parse(htmlForm))

func NewGet() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := tempForm.Execute(res, content); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func NewPost() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		content.Name = req.FormValue("name")
		content.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(content.Name)))
		if err := tempForm.Execute(res, content); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}
