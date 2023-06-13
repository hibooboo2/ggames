package rest

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/hibooboo2/glog"
)

func Recover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				RespondJson(w, http.StatusInternalServerError, map[string]string{
					"panic": fmt.Sprintf("%v", r),
					"stack": string(stack),
				})
				glog.Println("Stack:", string(stack))
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
