package cors

import (
	"net/http"
	"strings"

	"github.com/usman8786/Siyaho-POS-Agent/agent/license"
)

type Middleware struct {
	licenses *license.Manager
}

func New(licenses *license.Manager) *Middleware {
	return &Middleware{licenses: licenses}
}

func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && !m.licenses.IsOriginAllowed(origin) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"ok":false,"error":"Origin not allowed"}`))
			return
		}

		if origin != "" && m.licenses.IsOriginAllowed(origin) {
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
