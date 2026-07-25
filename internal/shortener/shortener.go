package shortener

import (
	"crypto/rand"
	"fmt"
	"net/url"
)

const (
	alphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength = 7
	maxRetries = 5
)

// Service implementa a logica de encurtamento sobre um Store.
type Service struct {
	store Store
}

// New cria um Service com o Store informado.
func New(store Store) *Service {
	return &Service{store: store}
}

// Shorten valida a URL, gera um codigo unico e o persiste.
func (s *Service) Shorten(rawURL string) (string, error) {
	if err := validate(rawURL); err != nil {
		return "", err
	}
	for i := 0; i < maxRetries; i++ {
		code, err := generateCode()
		if err != nil {
			return "", err
		}
		if _, err := s.store.Load(code); err == ErrNotFound {
			if err := s.store.Save(code, rawURL); err != nil {
				return "", err
			}
			return code, nil
		}
	}
	return "", fmt.Errorf("shortener: could not generate a unique code after %d retries", maxRetries)
}

// Resolve retorna a URL original de um codigo.
func (s *Service) Resolve(code string) (string, error) {
	return s.store.Load(code)
}

func validate(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("shortener: url is empty")
	}
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("shortener: invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("shortener: unsupported scheme %q", u.Scheme)
	}
	return nil
}

func generateCode() (string, error) {
	buf := make([]byte, codeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i, b := range buf {
		buf[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(buf), nil
}
