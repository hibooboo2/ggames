package logger

import (
	loggerinit "github.com/hibooboo2/ggames/allplay/logger/logger_init"
	"github.com/hibooboo2/glog"
)

var (
	LAuth     = loggerinit.Logger.RegisterNextLevel("AUTH")
	LBoard    = loggerinit.Logger.RegisterNextLevel("BOARD")
	LCards    = loggerinit.Logger.RegisterNextLevel("CARDS")
	LGames    = loggerinit.Logger.RegisterNextLevel("GAMES")
	LInit     = loggerinit.Logger.RegisterNextLevel("INIT")
	LPlayer   = loggerinit.Logger.RegisterNextLevel("PLAYERS")
	LPosition = loggerinit.Logger.RegisterNextLevel("POSITION")
	LScore    = loggerinit.Logger.RegisterNextLevel("SCORE")
	LToken    = loggerinit.Logger.RegisterNextLevel("TOKEN")
	LUsers    = loggerinit.Logger.RegisterNextLevel("USERS")
	All       = LCards | LGames | LBoard | LAuth | LUsers | LPlayer | LInit | LToken | LScore
)

func init() {
	SetPrefix("ALLPLAY")
	SetLevel(glog.LevelDebug | LScore | LPlayer | LPosition | LToken | LGames | All)
}

var (
	SetPrefix = loggerinit.Logger.SetPrefix
	SetLevel  = loggerinit.Logger.SetLevel
	AtLevel   = loggerinit.Logger.AtLevel
	AtLevelf  = loggerinit.Logger.AtLevelf
	AtLevelln = loggerinit.Logger.AtLevelln
	Auth      = loggerinit.Logger.CustomLogAtLevel(LAuth)
	Authf     = loggerinit.Logger.CustomLogAtLevelf(LAuth)
	Authln    = loggerinit.Logger.CustomLogAtLevelln(LAuth)
	Board     = loggerinit.Logger.CustomLogAtLevel(LBoard)
	Boardf    = loggerinit.Logger.CustomLogAtLevelf(LBoard)
	Boardln   = loggerinit.Logger.CustomLogAtLevelln(LBoard)
	Cards     = loggerinit.Logger.CustomLogAtLevel(LCards)
	Cardsf    = loggerinit.Logger.CustomLogAtLevelf(LCards)
	Cardsln   = loggerinit.Logger.CustomLogAtLevelln(LCards)
	Games     = loggerinit.Logger.CustomLogAtLevel(LGames)
	Gamesf    = loggerinit.Logger.CustomLogAtLevelf(LGames)
	Gamesln   = loggerinit.Logger.CustomLogAtLevelln(LGames)
	Initf     = loggerinit.Logger.CustomLogAtLevelf(LInit)
	Player    = loggerinit.Logger.CustomLogAtLevel(LPlayer)
	Playerf   = loggerinit.Logger.CustomLogAtLevelf(LPlayer)
	Playerln  = loggerinit.Logger.CustomLogAtLevelln(LPlayer)
	Token     = loggerinit.Logger.CustomLogAtLevel(LToken)
	Tokenf    = loggerinit.Logger.CustomLogAtLevelf(LToken)
	Tokenln   = loggerinit.Logger.CustomLogAtLevelln(LToken)
	Users     = loggerinit.Logger.CustomLogAtLevel(LUsers)
	Usersf    = loggerinit.Logger.CustomLogAtLevelf(LUsers)
	Usersln   = loggerinit.Logger.CustomLogAtLevelln(LUsers)

	Debug   = loggerinit.Logger.Debug
	Debugf  = loggerinit.Logger.Debugf
	Debugln = loggerinit.Logger.Debugln

	Info   = loggerinit.Logger.Info
	Infof  = loggerinit.Logger.Infof
	Infoln = loggerinit.Logger.Infoln

	Warn   = loggerinit.Logger.Warn
	Warnf  = loggerinit.Logger.Warnf
	Warnln = loggerinit.Logger.Warnln

	Error   = loggerinit.Logger.Error
	Errorf  = loggerinit.Logger.Errorf
	Errorln = loggerinit.Logger.Errorln

	Fatal   = loggerinit.Logger.Fatal
	Fatalf  = loggerinit.Logger.Fatalf
	Fatalln = loggerinit.Logger.Fatalln
)
