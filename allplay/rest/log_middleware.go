package rest

import (
	"net/http"

	"github.com/hibooboo2/glog"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Debugf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
		glog.Debug("REQUEST END\n\n\n\n")
	})
}
