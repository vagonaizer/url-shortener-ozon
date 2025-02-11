package restserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vagonaizer/url-shortener-ozon/internal/service"
)

type server struct {
	svc *service.URLShortenerService
}

// New возвращает настроенный HTTP-хендлер с зарегистрированными маршрутами.
func New(svc *service.URLShortenerService) http.Handler {
	s := &server{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", s.shortenHandler)
	mux.HandleFunc("/original/", s.originalHandler)
	return mux
}

func (s *server) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		OriginalURL string `json:"original_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OriginalURL == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	shortURL := s.svc.ShortenURL(req.OriginalURL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

func (s *server) originalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// URL: /original/shortURL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	shortURL := parts[2]
	originalURL, exists := s.svc.GetOriginalURL(shortURL)
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}
