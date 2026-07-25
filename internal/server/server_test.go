package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bruno1186/url-shortener/internal/shortener"
)

func newTestServer() http.Handler {
	svc := shortener.New(shortener.NewMemoryStore())
	return New(svc, "http://localhost:8080").Routes()
}

func TestHealthz(t *testing.T) {
	h := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestShortenEndpoint(t *testing.T) {
	h := newTestServer()
	body := strings.NewReader(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	var resp shortenResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Code == "" {
		t.Fatal("expected a non-empty code")
	}
	if !strings.HasSuffix(resp.ShortURL, resp.Code) {
		t.Fatalf("short_url %q should end with code %q", resp.ShortURL, resp.Code)
	}
}

func TestShortenEndpointInvalidBody(t *testing.T) {
	h := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("not-json"))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestRedirectFlow(t *testing.T) {
	h := newTestServer()

	body := strings.NewReader(`{"url":"https://example.com/docs"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var resp shortenResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	redirectReq := httptest.NewRequest(http.MethodGet, "/"+resp.Code, nil)
	redirectRec := httptest.NewRecorder()
	h.ServeHTTP(redirectRec, redirectReq)

	if redirectRec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", redirectRec.Code)
	}
	if loc := redirectRec.Header().Get("Location"); loc != "https://example.com/docs" {
		t.Fatalf("unexpected redirect location: %q", loc)
	}
}

func TestRedirectNotFound(t *testing.T) {
	h := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
