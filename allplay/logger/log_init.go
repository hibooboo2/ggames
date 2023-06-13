package logger

import (
	"os"

	"github.com/hibooboo2/glog"
)

const (
	LInit = 0b1000000000 >> iota
	LCards
	LGames
	LBoard
	LAuth
	LUsers
	LPlayer
	LToken
	LPosition

	All = LCards | LGames | LBoard | LAuth | LUsers | LPlayer | LInit | LToken
)

func init() {
	glog.RegisterLevel(LInit, "INIT")
	glog.RegisterLevel(LCards, "CARDS")
	glog.RegisterLevel(LGames, "GAMES")
	glog.RegisterLevel(LBoard, "BOARD")
	glog.RegisterLevel(LAuth, "AUTH")
	glog.RegisterLevel(LUsers, "USERS")
	glog.RegisterLevel(LPlayer, "PLAYERS")
	glog.RegisterLevel(LToken, "TOKEN")
	glog.RegisterLevel(LPosition, "POSITION")
	allplayLogger.SetPrefix("ALLPLAY")
}

var allplayLogger = glog.NewLogger(os.Stdout, glog.DefaultLevel|All)

func SetLevel(level int) {
	allplayLogger.SetLevel(level)
}

func AtLevel(level int, args ...interface{}) {
	allplayLogger.AtLevel(level, args...)
}

func AtLevelf(level int, format string, args ...interface{}) {
	allplayLogger.AtLevelf(level, format, args...)
}

func AtLevelln(level int, args ...interface{}) {
	allplayLogger.AtLevelln(level, args...)
}

func Cardsf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LCards, msg, args...)
}

func Gamesf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LGames, msg, args...)
}

func Boardf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LBoard, msg, args...)
}

func Authf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LAuth, msg, args...)
}

func Usersf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LUsers, msg, args...)
}

func Cardsln(args ...interface{}) {
	allplayLogger.AtLevelln(LCards, args...)
}

func Gamesln(args ...interface{}) {
	allplayLogger.AtLevelln(LGames, args...)
}

func Boardln(args ...interface{}) {
	allplayLogger.AtLevelln(LBoard, args...)
}

func Authln(args ...interface{}) {
	allplayLogger.AtLevelln(LAuth, args...)
}

func Usersln(args ...interface{}) {
	allplayLogger.AtLevelln(LUsers, args...)
}

func Cards(args ...interface{}) {
	allplayLogger.AtLevel(LCards, args...)
}

func Games(args ...interface{}) {
	allplayLogger.AtLevel(LGames, args...)
}

func Board(args ...interface{}) {
	allplayLogger.AtLevel(LBoard, args...)
}

func Auth(args ...interface{}) {
	allplayLogger.AtLevel(LAuth, args...)
}

func Users(args ...interface{}) {
	allplayLogger.AtLevel(LUsers, args...)
}

func Player(args ...interface{}) {
	allplayLogger.AtLevel(LPlayer, args...)
}

func Playerf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LPlayer, msg, args...)
}

func Playerln(args ...interface{}) {
	allplayLogger.AtLevelln(LPlayer, args...)
}

func Initf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LInit, msg, args...)
}

func Tokenf(msg string, args ...interface{}) {
	allplayLogger.AtLevelf(LToken, msg, args...)
}

func Tokenln(args ...interface{}) {
	allplayLogger.AtLevelln(LToken, args...)
}

func Token(args ...interface{}) {
	allplayLogger.AtLevel(LToken, args...)
}
