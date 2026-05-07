package shortener

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// NewMux wires the HTTP routes for the URL shortener. It uses Go 1.22+
// method-aware routing so the handlers stay readable.
func NewMux(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /shorten", shortenHandler(store))
	mux.HandleFunc("GET /r/{code}", redirectHandler(store))
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func shortenHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return
		}
		body.URL = strings.TrimSpace(body.URL)
		if body.URL == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "url is required"})
			return
		}
		if _, err := url.ParseRequestURI(body.URL); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "url must be absolute"})
			return
		}
		parsed, _ := url.Parse(body.URL)
		if parsed == nil || parsed.Scheme == "" || parsed.Host == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "url must be absolute"})
			return
		}

		code, err := store.Save(body.URL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"code":         code,
			"original_url": body.URL,
			"short_path":   "/r/" + code,
		})
	}
}

func redirectHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		target, ok := store.Get(code)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		http.Redirect(w, r, target, http.StatusFound)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
