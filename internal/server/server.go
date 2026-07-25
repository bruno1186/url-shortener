package server

import (
	"encoding/json"
	"net/http"

	"github.com/bruno1186/url-shortener/internal/shortener"
)

// Server expõe a API HTTP do encurtador.
type Server struct {
	svc     *shortener.Service
	baseURL string
}

// New cria um Server. baseURL e usado para montar a URL curta na resposta.
func New(svc *shortener.Service, baseURL string) *Server {
	return &Server{svc: svc, baseURL: baseURL}
}

// Routes registra as rotas e retorna o handler raiz.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /api/shorten", s.handleShorten)
	mux.HandleFunc("GET /{code}", s.handleRedirect)
	return mux
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json body"})
		return
	}
	code, err := s.svc.Shorten(req.URL)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, shortenResponse{
		Code:     code,
		ShortURL: s.baseURL + "/" + code,
	})
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	target, err := s.svc.Resolve(code)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "code not found"})
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
