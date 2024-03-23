package main

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/npmaile/htmx-pong/gamestate"
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

var id string

func main() {
	gs := gamestate.Init()
	go gs.StartProcessing()

	http.Handle("/static/", http.FileServerFS(staticFS))
	http.HandleFunc("/", index)
	http.HandleFunc("/pong", pongFunc(gs))
	http.HandleFunc("/update/{id}", updateFunc(gs))
	http.HandleFunc("/update/up/{id}", updateFuncUp(gs))
	http.HandleFunc("/update/down/{id}", updateFuncDown(gs))
	fmt.Println("running on 8080")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.templ.html", nil)
}
func pongFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Res := make(chan gamestate.Game)
		gs.NewGameRequests <- gamestate.NewGameRequest{
			Res: Res,
		}
		ng := <-Res
		err := templates.ExecuteTemplate(w, "pong.templ.html", ng)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

func updateFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: r.PathValue("id"),
			A:  gamestate.NoAction,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
	}
}

func updateFuncUp(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: r.PathValue("id"),
			A:  gamestate.Lup,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
	}
}

func updateFuncDown(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: r.PathValue("id"),
			A:  gamestate.Ldown,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
	}
}
