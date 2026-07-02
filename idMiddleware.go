package main

import (
	"crypto/rand"
	"net/http"
	"strings"
)

func newIDRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		id := r.Header.Get("X-Request-ID")
		if strings.TrimSpace(id) == "" {
			randId := rand.Text()
			w.Header().Set("X-Request-ID", randId)
		} else {
			w.Header().Set("X-Request-ID", id)
		}

		next.ServeHTTP(w, r)
	})
}
