package shortener

import "testing"

func newService() *Service {
	return New(NewMemoryStore())
}

func TestShortenAndResolve(t *testing.T) {
	svc := newService()

	code, err := svc.Shorten("https://example.com/path?q=1")
	if err != nil {
		t.Fatalf("Shorten returned error: %v", err)
	}
	if len(code) != codeLength {
		t.Fatalf("expected code length %d, got %d", codeLength, len(code))
	}

	got, err := svc.Resolve(code)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if want := "https://example.com/path?q=1"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestShortenInvalidURL(t *testing.T) {
	svc := newService()

	cases := []string{"", "not-a-url", "ftp://example.com", "javascript:alert(1)"}
	for _, c := range cases {
		if _, err := svc.Shorten(c); err == nil {
			t.Errorf("expected error for input %q, got nil", c)
		}
	}
}

func TestResolveUnknownCode(t *testing.T) {
	svc := newService()

	if _, err := svc.Resolve("missing"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestShortenProducesDistinctCodes(t *testing.T) {
	svc := newService()

	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := svc.Shorten("https://example.com")
		if err != nil {
			t.Fatalf("Shorten returned error: %v", err)
		}
		if seen[code] {
			t.Fatalf("duplicate code generated: %s", code)
		}
		seen[code] = true
	}
}
