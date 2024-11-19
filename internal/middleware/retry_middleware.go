package middleware

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func RetryMiddleware(maxRetries int, delay time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			for i := 0; i <= maxRetries; i++ {
				recorder := newResponseRecorder(w)
				next.ServeHTTP(recorder, r)

				if recorder.statusCode >= 500 {
					err = errors.New("requisição falhou, tentando novamente")
					time.Sleep(delay)
					continue
				}
				return
			}
			log.Printf("Falha após %d tentativas: %v", maxRetries, err)
			http.Error(w, "Falha no servidor, tente mais tarde", http.StatusInternalServerError)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
}
