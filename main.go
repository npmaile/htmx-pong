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

	Res := make(chan gamestate.NewGameResponse)
	gs.NewGameRequests <- gamestate.NewGameRequest{
		Res: Res,
	}
	ng := <-Res
	id = ng.GameID
	http.Handle("/static/", http.FileServerFS(staticFS))
	http.HandleFunc("/", pong)
	http.HandleFunc("/update", updateFunc(gs))
	http.HandleFunc("/update/left/up", updateFuncLeftUp(gs))
	http.HandleFunc("/update/left/down", updateFuncLeftDown(gs))
	http.HandleFunc("/update/right/up", updateFuncRightUp(gs))
	http.HandleFunc("/update/right/down", updateFuncRightDown(gs))
	http.HandleFunc("/testkey", testkey)
	fmt.Println("running on 8080")
	http.ListenAndServe(":8080", nil)
}

func pong(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.templ.html", nil)
}

func testkey(w http.ResponseWriter, r *http.Request) {
	fmt.Println("up key pressed")
}

func updateFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: id,
			A:  gamestate.NoAction,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
		// render the gamestate template to the thing
	}
}

func updateFuncLeftUp(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: id,
			A:  gamestate.Lup,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
		// render the gamestate template to the thing
	}
}

func updateFuncLeftDown(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: id,
			A:  gamestate.Ldown,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
		// render the gamestate template to the thing
	}
}

func updateFuncRightUp(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: id,
			A:  gamestate.Rup,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
		// render the gamestate template to the thing
	}
}

func updateFuncRightDown(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch: ch,
			ID: id,
			A:  gamestate.Rdown,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
		// render the gamestate template to the thing
	}
}
