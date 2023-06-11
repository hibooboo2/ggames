package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen/db"
	"github.com/hibooboo2/ggames/allplay/rest"
)

func main() {
	log.SetFlags(log.Lshortfile)

	r := chi.NewRouter()
	r.Handle("/static/js/{filename}", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./pollen/static/js"))))
	r.Handle("/static/images/{imagename}", http.StripPrefix("/static/images/", http.FileServer(http.Dir("./pollen/static/images/"))))
	r.Handle("/static/css/{filename}", http.StripPrefix("/static/css/", http.FileServer(http.Dir("./pollen/static/css/"))))
	r.HandleFunc("/login", rest.LoginHandler)
	r.HandleFunc("/register", rest.Register)

	a := r.With(rest.Recover, rest.BasicAuth)
	a.HandleFunc("/", homePage)
	a.Post("/game", newgame)
	a.Get("/game/join/{game_id}", joinGame)
	a.Post("/game/start/{game_id}", startGame)
	a.Get("/game/{game_id}", renderGame)
	a.Post("/game/{game_id}/play/card/{card_id}", playCard)

	a.Post("/tempID/", rest.TempID)

	http.ListenAndServe(":8080", r)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	homePageTmpl := template.Must(template.New("home").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
		    <script src="/static/js/functions.js"> </script>
			<meta charset="utf-8">
		</head>
		<body>
		    <div id="mainbox">
				<h1>Welcome to Pollen</h1>
				<button onclick="newGame()">New Game</button>
			</div>
		</body>
	</html>
	`))

	homePageTmpl.Execute(w, nil)
}

func newgame(w http.ResponseWriter, r *http.Request) {
	id := uuid.Must(uuid.NewV4())

	err := db.NewGame(id, rest.GetUsername(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGameTmpl := template.Must(template.New("newgame").Parse(`
		<div>
		    <h1>Current Game {{.GameID}}</h1>
			<button onclick='createJoinGameLink({{.GameID}})'>Copy Join Game Link</button>
			<button onclick='startGame({{.GameID}})'>Start Game</button>
		</div>
	`))

	newGameTmpl.Execute(w, struct {
		GameID string
	}{id.String()})
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	gameID := rest.GetGameID(r)
	username := rest.GetUsername(r)

	err := db.AddGameUser(gameID, username)
	if err != nil {
		rest.ResponndError(w, http.StatusBadRequest, err.Error())
		return
	}

	rest.ResponndJson(w, http.StatusAccepted, map[string]string{"game": gameID.String(), "AddedUser": username})
}

func startGame(w http.ResponseWriter, r *http.Request) {
	gameID := rest.GetGameID(r)

	username := rest.GetUsername(r)

	err := db.AddGameUser(gameID, username)
	if err != nil {
		rest.ResponndError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.StartGame(gameID)
	if err != nil {
		rest.ResponndError(w, http.StatusBadRequest, err.Error())
		return
	}

	rest.ResponndJson(w, http.StatusCreated, map[string]string{"game": gameID.String(), "status": "started"})
}

func renderGame(w http.ResponseWriter, r *http.Request) {
	var err error
	gameID := rest.GetGameID(r)

	g := db.GetGame(gameID)

	err = g.Render(w)
	if err != nil {
		panic(err)
	}
}

func playCard(w http.ResponseWriter, r *http.Request) {
	gameID := rest.GetGameID(r)
	username := rest.GetUsername(r)

	g := db.GetGame(gameID)

	err := g.PlayCard(username, rest.GetCardID(r), rest.GetPosition(r))
	if err != nil {
		rest.ResponndError(w, http.StatusBadRequest, err.Error())
		return
	}

	g.NextPlayer()

	w.WriteHeader(http.StatusAccepted)
}
