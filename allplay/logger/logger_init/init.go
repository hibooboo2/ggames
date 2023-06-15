package loggerinit

import (
	"os"

	"github.com/hibooboo2/glog"
)

var Logger = glog.NewLogger(os.Stdout, glog.DefaultLevel)
