package logger

import (
	"os"

	"github.com/hibooboo2/glog"
)

var (
	LAuth     = glog.RegisterNextLevel("AUTH")
	LBoard    = glog.RegisterNextLevel("BOARD")
	LCards    = glog.RegisterNextLevel("CARDS")
	LGames    = glog.RegisterNextLevel("GAMES")
	LInit     = glog.RegisterNextLevel("INIT")
	LPlayer   = glog.RegisterNextLevel("PLAYERS")
	LPosition = glog.RegisterNextLevel("POSITION")
	LScore    = glog.RegisterNextLevel("SCORE")
	LToken    = glog.RegisterNextLevel("TOKEN")
	LUsers    = glog.RegisterNextLevel("USERS")
	All       = LCards | LGames | LBoard | LAuth | LUsers | LPlayer | LInit | LToken | LScore
)

func init() {
	SetPrefix("ALLPLAY")
	SetLevel(glog.LevelDebug | LScore | LPlayer | LPosition | LToken | LGames | All)
	glog.DefaultLogger = glog.NewLogger(os.Stdout, glog.DefaultLevel)

}

var allpplayLogger = glog.DefaultLogger

var (
	SetPrefix = glog.SetPrefix
	SetLevel  = glog.SetLevel
	AtLevel   = glog.AtLevel
	AtLevelf  = glog.AtLevelf
	AtLevelln = glog.AtLevelln
	Auth      = glog.CustomLogAtLevel(LAuth)
	Authf     = glog.CustomLogAtLevelf(LAuth)
	Authln    = glog.CustomLogAtLevelln(LAuth)
	Board     = glog.CustomLogAtLevel(LBoard)
	Boardf    = glog.CustomLogAtLevelf(LBoard)
	Boardln   = glog.CustomLogAtLevelln(LBoard)
	Cards     = glog.CustomLogAtLevel(LCards)
	Cardsf    = glog.CustomLogAtLevelf(LCards)
	Cardsln   = glog.CustomLogAtLevelln(LCards)
	Games     = glog.CustomLogAtLevel(LGames)
	Gamesf    = glog.CustomLogAtLevelf(LGames)
	Gamesln   = glog.CustomLogAtLevelln(LGames)
	Initf     = glog.CustomLogAtLevelf(LInit)
	Player    = glog.CustomLogAtLevel(LPlayer)
	Playerf   = glog.CustomLogAtLevelf(LPlayer)
	Playerln  = glog.CustomLogAtLevelln(LPlayer)
	Token     = glog.CustomLogAtLevel(LToken)
	Tokenf    = glog.CustomLogAtLevelf(LToken)
	Tokenln   = glog.CustomLogAtLevelln(LToken)
	Users     = glog.CustomLogAtLevel(LUsers)
	Usersf    = glog.CustomLogAtLevelf(LUsers)
	Usersln   = glog.CustomLogAtLevelln(LUsers)

	Debug   = glog.Debug
	Debugf  = glog.Debugf
	Debugln = glog.Debugln

	Info   = glog.Info
	Infof  = glog.Infof
	Infoln = glog.Infoln

	Warn   = glog.Warn
	Warnf  = glog.Warnf
	Warnln = glog.Warnln

	Error   = glog.Error
	Errorf  = glog.Errorf
	Errorln = glog.Errorln

	Fatal   = glog.Fatal
	Fatalf  = glog.Fatalf
	Fatalln = glog.Fatalln
)
