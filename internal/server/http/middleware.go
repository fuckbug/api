package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fuckbug/api/internal/server/http/handlers"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	ResponseCode int
}

func (l *LoggingResponseWriter) WriteHeader(code int) {
	l.ResponseCode = code
	l.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger handlers.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lrw := LoggingResponseWriter{
			ResponseWriter: w,
			ResponseCode:   http.StatusOK,
		}
		next.ServeHTTP(&lrw, r)
		latency := time.Since(startTime)

		logger.Info(
			fmt.Sprintf(
				"%s [%s] %s %s HTTP/%s %d %s \"%s\"",
				r.RemoteAddr,
				time.Now().Format(time.RFC1123Z),
				r.Method, r.URL.Path, r.Proto,
				lrw.ResponseCode,
				latency,
				r.Header.Get("User-Agent"),
			),
		)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
