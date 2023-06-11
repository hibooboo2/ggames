package rest

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func Recover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				ResponndJson(w, http.StatusInternalServerError, map[string]string{
					"panic": fmt.Sprintf("%v", r),
					"stack": string(stack),
				})
				log.Println("Stack:", string(stack))
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
