package logger

import (
	"os"

	"github.com/hibooboo2/glog"
)

var (
	LAuth     = allplayLogger.RegisterNextLevel("AUTH")
	LBoard    = allplayLogger.RegisterNextLevel("BOARD")
	LCards    = allplayLogger.RegisterNextLevel("CARDS")
	LGames    = allplayLogger.RegisterNextLevel("GAMES")
	LInit     = allplayLogger.RegisterNextLevel("INIT")
	LPlayer   = allplayLogger.RegisterNextLevel("PLAYERS")
	LPosition = allplayLogger.RegisterNextLevel("POSITION")
	LScore    = allplayLogger.RegisterNextLevel("SCORE")
	LToken    = allplayLogger.RegisterNextLevel("TOKEN")
	LUsers    = allplayLogger.RegisterNextLevel("USERS")
	All       = LCards | LGames | LBoard | LAuth | LUsers | LPlayer | LInit | LToken | LScore
)

func init() {
	// LInit = allplayLogger.RegisterNextLevel("INIT")

	SetPrefix("ALLPLAY")
	SetLevel(glog.LevelDebug | LScore | LPlayer | LPosition | LToken | LGames | All)
	glog.DefaultLogger = glog.NewLogger(os.Stdout, glog.DefaultLevel)
}

var allplayLogger = glog.DefaultLogger

var (
	SetPrefix = allplayLogger.SetPrefix
	SetLevel  = allplayLogger.SetLevel
	AtLevel   = allplayLogger.AtLevel
	AtLevelf  = allplayLogger.AtLevelf
	AtLevelln = allplayLogger.AtLevelln
	Auth      = allplayLogger.CustomLogAtLevel(LAuth)
	Authf     = allplayLogger.CustomLogAtLevelf(LAuth)
	Authln    = allplayLogger.CustomLogAtLevelln(LAuth)
	Board     = allplayLogger.CustomLogAtLevel(LBoard)
	Boardf    = allplayLogger.CustomLogAtLevelf(LBoard)
	Boardln   = allplayLogger.CustomLogAtLevelln(LBoard)
	Cards     = allplayLogger.CustomLogAtLevel(LCards)
	Cardsf    = allplayLogger.CustomLogAtLevelf(LCards)
	Cardsln   = allplayLogger.CustomLogAtLevelln(LCards)
	Games     = allplayLogger.CustomLogAtLevel(LGames)
	Gamesf    = allplayLogger.CustomLogAtLevelf(LGames)
	Gamesln   = allplayLogger.CustomLogAtLevelln(LGames)
	Initf     = allplayLogger.CustomLogAtLevelf(LInit)
	Player    = allplayLogger.CustomLogAtLevel(LPlayer)
	Playerf   = allplayLogger.CustomLogAtLevelf(LPlayer)
	Playerln  = allplayLogger.CustomLogAtLevelln(LPlayer)
	Token     = allplayLogger.CustomLogAtLevel(LToken)
	Tokenf    = allplayLogger.CustomLogAtLevelf(LToken)
	Tokenln   = allplayLogger.CustomLogAtLevelln(LToken)
	Users     = allplayLogger.CustomLogAtLevel(LUsers)
	Usersf    = allplayLogger.CustomLogAtLevelf(LUsers)
	Usersln   = allplayLogger.CustomLogAtLevelln(LUsers)
)
