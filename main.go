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
			A:  0,
		}
		a := <-ch
		fmt.Printf("got: %+v", a)
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
		// render the gamestate template to the thing
	}
}
