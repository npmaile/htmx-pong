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
	http.HandleFunc("/friendroom", startFriendRoomFunc(gs))
	http.HandleFunc("/matchmaking", startMatchMakingFunc(gs))
	//http.HandleFunc("/friendConnect", friendConnectFunc(gs))

	http.HandleFunc("/update/{id}/{player}", updateFunc(gs, gamestate.NoAction))
	http.HandleFunc("/update/up/{id}/{player}", updateFunc(gs, gamestate.Up))
	http.HandleFunc("/update/down/{id}/{player}", updateFunc(gs, gamestate.Down))
	fmt.Println("running on 8080")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.templ.html", nil)
}

func startFriendRoomFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gs.NewWaitingRoomRequests <- gamestate.WaitingRoomRequest{}
	}
}

/*
func friendConnectFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
*/

func startMatchMakingFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(chan gamestate.GameResponse)
		gs.NewMatchMakingRequests <- gamestate.MatchMakingRequest{
			Res: res,
		}
		result := <-res
		err := templates.ExecuteTemplate(w, "pong.templ.html", result)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

func pongFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Res := make(chan gamestate.GameResponse)
		ng := <-Res
		err := templates.ExecuteTemplate(w, "pong.templ.html", ng)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

func updateFunc(gs gamestate.GameStateSingleton, updateAction gamestate.Action) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.Game)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Ch:       ch,
			ID:       r.PathValue("id"),
			A:        updateAction,
			PlayerID: r.PathValue("player"),
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
	}
}
