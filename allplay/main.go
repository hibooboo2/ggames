package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen"
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

	a := r.With(rest.LogRequest, rest.Recover, rest.BasicAuth)
	a.HandleFunc("/", homePage)
	a.Post("/game", newgame)
	a.Get("/game/join/{game_id}", joinGame)
	a.Post("/game/start/{game_id}", startGame)
	a.Get("/game/play/{game_id}", playGame)
	a.Get("/game/{game_id}/render", renderGame)
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
				{{ range $game := .Games }}
					<a href="/game/play/{{ $game.GetID }}">Play Game</a>
				{{end}}
			</div>
		</body>
	</html>
	`))

	username := rest.GetUsername(r)
	activeGames := db.GetActiveGames(username)
	log.Println("User: ", username, " Active Games: ", len(activeGames))

	homePageTmpl.Execute(w, struct {
		Games []*pollen.Game
	}{
		activeGames,
	})
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
	gameID := pollen.GetGameID(r)
	username := rest.GetUsername(r)

	err := db.AddGameUser(gameID, username)
	if err != nil {
		rest.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.(http.Flusher).Flush()
	http.Redirect(w, r, "/game/play/"+gameID.String(), http.StatusSeeOther)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	gameID := pollen.GetGameID(r)

	username := rest.GetUsername(r)

	err := db.AddGameUser(gameID, username)
	if err != nil {
		rest.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.StartGame(gameID)
	if err != nil {
		rest.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, "/game/play/"+gameID.String(), http.StatusFound)
}

func renderGame(w http.ResponseWriter, r *http.Request) {
	var err error
	gameID := pollen.GetGameID(r)
	g := db.GetGame(gameID)
	if g == nil {
		rest.RespondError(w, http.StatusNotFound, "Game not found")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	err = g.Render(w, rest.GetUsername(r))
	if err != nil {
		rest.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func playCard(w http.ResponseWriter, r *http.Request) {
	gameID := pollen.GetGameID(r)
	username := rest.GetUsername(r)

	g := db.GetGame(gameID)

	err := g.PlayCard(username, pollen.GetCardID(r), pollen.GetPosition(r))
	if err != nil {
		rest.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	g.NextPlayer()

	w.WriteHeader(http.StatusAccepted)
}

var playGameTmpl = template.Must(template.New("playgame").Parse(`
	<!DOCTYPE html>
	<html>
        <head>
			<script src="/static/js/functions.js"> </script>
			<link rel="stylesheet" href="/static/css/main.css">
			<meta charset="utf-8">
		</head>
		<body>
		    <div id="mainbox">
                <h1>Current Game {{.GameID}}</h1>
				<h2>Playing as {{.Username}}</h2>
				<div id="gamebox">
				</div>
		    </div>
		</body>
		<script>
			renderGame({{.GameID}})
		</script>
	</html>
`))

func playGame(w http.ResponseWriter, r *http.Request) {
	tmplContext := struct {
		GameID   string
		Username string
	}{
		pollen.GetGameID(r).String(),
		rest.GetUsername(r),
	}

	playGameTmpl.Execute(w, tmplContext)
}
