package main

import (
	"embed"
	"net/http"
	"text/template"
)

//go:embed static
var staticFS embed.FS

//go:embed templates
var templatesFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(templatesFS, "templates/*")
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/static/", http.FileServerFS(staticFS))
	http.HandleFunc("/", pong)
	http.HandleFunc("/gamestate", gamestate)
	http.ListenAndServe(":8080", nil)
}

func pong(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.templ.html", nil)
}

func gamestate(w http.ResponseWriter, r *http.Request) {
	// lookup the user
	// get their gamestate
	// render the gamestate template to the thing
}
