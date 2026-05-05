package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		defer func() {
			log.Println(req.Method, req.URL.Path, time.Since(start))
		}()

		next.ServeHTTP(res, req)
	})
}
