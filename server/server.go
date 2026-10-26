package server

import (
	"encoding/json"
	"net/http"

	"github.com/partite-ai/example-ticket-pipeline-shortlink/store"
	"github.com/partite-ai/example-ticket-pipeline-shortlink/validate"
)

func New(s *store.Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /links", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json: " + err.Error()})
			return
		}
		if err := validate.URL(body.URL); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		code := s.Put(body.URL)
		writeJSON(w, http.StatusCreated, map[string]string{
			"code":      code,
			"short_url": "http://" + r.Host + "/" + code,
		})
	})
	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		target, ok := s.Get(r.PathValue("code"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, target, http.StatusFound)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
