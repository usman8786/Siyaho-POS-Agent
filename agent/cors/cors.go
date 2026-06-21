package cors

import (
	"net/http"
	"strings"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	if !strings.Contains(w.Header().Get("Access-Control-Allow-Origin"), "://") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
