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
	http.HandleFunc("/pong/{id}/{player}", matchMakingWaiting(gs))
	http.HandleFunc("/friendroom", startFriendRoomFunc(gs))
	http.HandleFunc("/matchmaking", startMatchMakingFunc(gs))
	http.HandleFunc("/update/no-action/{id}/{player}", updateFunc(gs, gamestate.NoAction))
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

func startMatchMakingFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(chan gamestate.GameResponse)
		fmt.Println("trying to send on the channel")
		gs.NewMatchMakingRequests <- gamestate.MatchMakingRequest{
			Res: res,
		}
		fmt.Println("waiting in startMatchmaking")
		result := <-res
		fmt.Println("done waiting in startMatchmaking")
		err := templates.ExecuteTemplate(w, "pong.templ.html", result)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

func matchMakingWaiting(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		Res := make(chan gamestate.GameResponse)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Res:      Res,
			ID:       r.PathValue("id"),
			PlayerID: r.PathValue("player"),
		}
		fmt.Println("waiting in matchMakingWaiting")
		ng := <-Res
		fmt.Println("done waiting in matchMakingWaiting")
		err := templates.ExecuteTemplate(w, "pong.templ.html", ng)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

func updateFunc(gs gamestate.GameStateSingleton, updateAction gamestate.Action) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup the user
		ch := make(chan gamestate.GameResponse)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Res:      ch,
			ID:       r.PathValue("id"),
			PlayerID: r.PathValue("player"),
			A:        updateAction,
		}
		a := <-ch
		err := templates.ExecuteTemplate(w, "gamestate.templ.css", a)
		if err != nil {
			panic(err)
		}
	}
}
