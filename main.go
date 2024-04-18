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
	http.HandleFunc("/matchmaking", startMatchMakingFunc(gs))
	http.HandleFunc("/friendroom", openFriendRoomFunc(gs))
	http.HandleFunc("/cancel/{id}", cancelMatchMakingFunc(gs))
	http.HandleFunc("/update/no-action/{id}/{player}", updateFunc(gs, gamestate.NoAction))
	http.HandleFunc("/friendConnect", friendConnectFunc(gs))
	http.HandleFunc("/update/up/{id}/{player}", updateFunc(gs, gamestate.Up))
	http.HandleFunc("/update/down/{id}/{player}", updateFunc(gs, gamestate.Down))
	http.HandleFunc("/singlePlayer", singlePlayerFunc(gs))
	fmt.Println("running on 8080")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.templ.html", nil)
}

func singlePlayerFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(chan gamestate.GameResponse)
		gs.SinglePlayerRequests <- gamestate.SinglePlayerRequest{
			Res: res,
		}
		result := <-res
		err := templates.ExecuteTemplate(w, "pong.templ.html", result)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}

	}
}

func cancelMatchMakingFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.PathValue("id")
		gs.CancelRequests <- gamestate.CancelRequest{
			ID: gameID,
		}
		index(w, r)
	}
}

func friendConnectFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			return
		}
		gameID := r.FormValue("friendCode")
		res := make(chan gamestate.GameResponse)
		gs.FriendJoinRequests <- gamestate.FriendJoinRequest{
			Res: res,
			ID:  gameID,
		}
		result := <-res
		if result.Error != nil {
			//do something with this error
			return
		}
		err = templates.ExecuteTemplate(w, "pong.templ.html", result)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

func openFriendRoomFunc(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(chan gamestate.GameResponse)
		gs.NewWaitingRoomRequests <- gamestate.WaitingRoomRequest{
			Res: res,
		}
		result := <-res
		if result.Error != nil {
			//do something with this error
			return
		}
		err := templates.ExecuteTemplate(w, "pong.templ.html", result)
		if err != nil {
			fmt.Println("error from template:", err.Error())
		}
	}
}

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

func matchMakingWaiting(gs gamestate.GameStateSingleton) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Res := make(chan gamestate.GameResponse)
		gs.GameUpdateRequests <- gamestate.GameUpdateRequest{
			Res:      Res,
			ID:       r.PathValue("id"),
			PlayerID: r.PathValue("player"),
		}
		ng := <-Res
		if ng.Error == gamestate.ErrNOTREADYYET {
			w.WriteHeader(http.StatusTooEarly)
		}
		err := templates.ExecuteTemplate(w, "pong.templ.html", ng)
		if err != nil {
			fmt.Println("err from template:", err.Error())
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
