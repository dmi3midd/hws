package hws

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Server struct {
	cfg     *Config
	service URLService
}

func NewServer(cfg *Config, service URLService) *http.Server {
	s := &Server{cfg: cfg, service: service}
	mux := s.RegisterRoutes()
	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}
}

func (s *Server) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("POST /create/random", s.CreateRandomHandler)
	mux.HandleFunc("POST /create/custom", s.CreateCustomHandler)
	mux.HandleFunc("GET /{alias}", s.RedirectHandler)
	return mux
}

type CreateRandomRequest struct {
	Url string `json:"url"`
}

func (s *Server) CreateRandomHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateRandomRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	alias, err := s.service.SaveRandom(r.Context(), reqBody.Url)
	if err != nil {
		slog.Error(
			"failed to create random alias",
			slog.String("error", err.Error()),
			slog.String("url", reqBody.Url),
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	slog.Info(
		"created random alias",
		slog.String("alias", alias),
		slog.String("url", reqBody.Url),
	)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(alias))
}

type CreateCustomRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

func (s *Server) CreateCustomHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateCustomRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	alias, err := s.service.SaveCustom(r.Context(), reqBody.Url, reqBody.Alias)
	if err != nil {
		slog.Error(
			"failed to create random alias",
			slog.String("error", err.Error()),
			slog.String("url", reqBody.Url),
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	slog.Info(
		"created random alias",
		slog.String("alias", alias),
		slog.String("url", reqBody.Url),
	)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(alias))
}

func (s *Server) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")
	if alias == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty alias"))
		return
	}

	url, err := s.service.Find(r.Context(), alias)
	if err != nil {
		slog.Info(
			"url not found",
			slog.String("error", err.Error()),
			slog.String("alias", alias),
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	slog.Info(
		"redirected successfully",
		slog.String("alias", alias),
		slog.String("to", url),
	)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
