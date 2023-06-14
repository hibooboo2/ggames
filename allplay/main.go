package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"

	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen"
	"github.com/hibooboo2/ggames/allplay/pollen/db"
	"github.com/hibooboo2/ggames/allplay/rest"
	"github.com/hibooboo2/glog"
)

func main() {
	// logger.SetFlags(logger.Lshortfile)

	r := chi.NewRouter()
	r.NotFound(notFoundHandler)
	r.Handle("/static/js/{filename}", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./pollen/static/js"))))
	r.Handle("/static/images/{imagename}", http.StripPrefix("/static/images/", http.FileServer(http.Dir("./pollen/static/images/"))))
	r.Handle("/static/css/{filename}", http.StripPrefix("/static/css/", http.FileServer(http.Dir("./pollen/static/css/"))))
	r.HandleFunc("/login", rest.LoginHandler)
	r.HandleFunc("/logout", rest.LogoutHandler)
	r.HandleFunc("/register", rest.Register)

	a := r.With(rest.LogRequest, rest.Recover, rest.BasicAuth)
	a.HandleFunc("/", homePage)
	a.Post("/game", newgame)
	a.Post("/game/{game_id}/invite", inviteUser)
	a.Get("/game/{game_id}/join", joinGame)
	a.Post("/game/{game_id}/start/", startGame)
	a.Get("/game/{game_id}/play/", playGame)
	a.Get("/game/{game_id}/render/", renderGame)
	a.Post("/game/{game_id}/hints/toggle/", toggleHints)
	a.Post("/game/{game_id}/play/card/{card_id}", playCard)
	a.Post("/game/{game_id}/play/token/{token_id}", playToken)

	a.Post("/tempID/", rest.TempID)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		glog.Fatal(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	homePageTmpl := template.Must(template.New("home").Parse(`
	{{ $username :=.Username }}
	<!DOCTYPE html>
	<html>
		<head>
		    <script src="/static/js/functions.js"> </script>
		    <script src="/static/js/hotkeys.js"> </script>
			<meta charset="utf-8">
		</head>
		<body>
		    <div id="mainbox">
				<h1>Welcome to Pollen</h1>
				<button onclick="newGame()">New Game</button>
				{{ range $game := .Games }}
				<ul>
					<li>
						<a href="/game/{{ $game.GetID }}/play/">Play Game: {{ $game.Name }} </a>
						{{ if $game.IsOwner $username }}
						<button onclick='startGame({{ $game.GetID }})'>Start Game</button>
						{{ end}}
					</li>
				</ul>
				{{end}}
				{{ range $game := .InvitedGames }}
					<ul>
						<li>
							<a href="/game/{{ $game.GetID }}/join">Join Game: {{ $game.Name }} </a>
						</li>
					</ul>
				{{end}}
			</div>
		</body>
	</html>
	`))

	username := rest.GetUsername(r)
	activeGames := db.GetActiveGames(username)
	invitedGames := db.GetInvitedGames(username)
	logger.AtLevelln(glog.LevelInfo, username, " Active Games: ", len(activeGames))
	logger.Usersln(username, " Active Games: ", len(activeGames))
	logger.Usersln(username, " Invited Games: ", len(invitedGames))

	homePageTmpl.Execute(w, struct {
		Games        []*pollen.Game
		InvitedGames []*pollen.Game
		Username     string
	}{
		activeGames,
		invitedGames,
		username,
	})
}

func inviteUser(w http.ResponseWriter, r *http.Request) {
	username := rest.GetUsername(r)
	gameID := pollen.GetGameID(r)
	logger.Usersln(username, " Invite Game: ", gameID)

	uname := r.URL.Query().Get("username")

	if err := db.InviteUserToGame(gameID, uname); err != nil {
		logger.Usersln(err)
		rest.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	rest.RespondJson(w, http.StatusAccepted, map[string]string{
		"status": "accepted",
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
			<button onclick='createJoinGameLink({{ .GameID }})'>Copy Join Game Link</button>
			<button onclick='startGame({{ .GameID }})'>Start Game</button>
			<input type="text" id="username" placeholder="User to invite"/>
			<button onclick='inviteUser({{ .GameID }})'>Invite User</button>
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

	http.Redirect(w, r, "/game/"+gameID.String()+"/play/", http.StatusFound)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	gameID := pollen.GetGameID(r)

	err := db.StartGame(gameID)
	if err != nil {
		rest.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, "/game/"+gameID.String()+"/play/", http.StatusFound)
}

func renderGame(w http.ResponseWriter, r *http.Request) {
	var err error
	gameID := pollen.GetGameID(r)
	g := db.GetGame(gameID)
	if g == nil {
		rest.RespondError(w, http.StatusNotFound, "Game not found")
		return
	}

	fw, ok := w.(pollen.FlusherWriter)
	if !ok {
		http.Error(w, "Event Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	err = g.Render(r.Context().Done(), fw, rest.GetUsername(r))
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

func playToken(w http.ResponseWriter, r *http.Request) {
	gameID := pollen.GetGameID(r)
	username := rest.GetUsername(r)

	g := db.GetGame(gameID)

	err := g.PlayToken(username, pollen.GetTokenID(r), pollen.GetPosition(r))
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
			<script src="/static/js/hotkeys.js"> </script>
			<link rel="stylesheet" href="/static/css/main.css">
			<meta charset="utf-8">
		</head>
		<body>
		    <div id="mainbox">
				<p class="right">Playing as {{.Username}}</p>
				<div id="gamebox">
				</div>
		    </div>
		</body>
		<script>
			renderGame("{{.GameID}}")
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

func toggleHints(w http.ResponseWriter, r *http.Request) {
	gameID := pollen.GetGameID(r)
	username := rest.GetUsername(r)

	g := db.GetGame(gameID)
	if g == nil {
		rest.RespondError(w, http.StatusNotFound, "Game not found")
		return
	}

	g.ToggleHints(username)

	w.WriteHeader(http.StatusAccepted)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	t, err := pollen.LoadTemplate("404.html.tmpl")
	if err != nil {
		panic(err)
	}
	var info = struct {
		URL string
	}{
		URL: r.URL.String(),
	}

	err = t.ExecuteTemplate(w, "404", info)
	if err != nil {
		rest.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
