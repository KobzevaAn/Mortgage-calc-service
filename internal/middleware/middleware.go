// Package middleware содержит HTTP middleware.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware логирует HTTP-запросы, фиксируя время обработки и статус-код ответа.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)
		duration := time.Since(start).Nanoseconds()
		log.Printf("%s status_code: %d, duration: %d ns", time.Now().Format("2006/01/02 15:04:05"), lw.statusCode, duration)
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *logResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}
