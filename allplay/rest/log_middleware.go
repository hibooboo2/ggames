package rest

import (
	"net/http"
	"time"

	"github.com/hibooboo2/glog"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Debugf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		start := time.Now()
		next.ServeHTTP(w, r)
		glog.Debugf("request end took:%s", time.Since(start))
	})
}
