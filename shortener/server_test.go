package shortener

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := NewMux(NewStore())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["ok"] != true {
		t.Errorf("expected ok=true, got %v", got["ok"])
	}
}

func TestShortenAndRedirectRoundTrip(t *testing.T) {
	mux := NewMux(NewStore())

	body := `{"url":"https://example.com/very/long/path"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rr.Code, rr.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	code, _ := got["code"].(string)
	if len(code) != codeLen {
		t.Fatalf("expected code length %d, got %d", codeLen, len(code))
	}

	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/r/"+code, nil)
	mux.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rr2.Code)
	}
	if loc := rr2.Header().Get("Location"); loc != "https://example.com/very/long/path" {
		t.Errorf("Location mismatch: %q", loc)
	}
}

func TestShortenRejectsInvalidJSON(t *testing.T) {
	mux := NewMux(NewStore())

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader([]byte("not-json")))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestShortenRequiresURL(t *testing.T) {
	mux := NewMux(NewStore())

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":""}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestShortenRejectsRelativeURL(t *testing.T) {
	mux := NewMux(NewStore())

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":"/just/a/path"}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestRedirectMissingCodeReturnsNotFound(t *testing.T) {
	mux := NewMux(NewStore())

	req := httptest.NewRequest(http.MethodGet, "/r/zzzzzzzz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestStoreGeneratesDistinctCodes(t *testing.T) {
	store := NewStore()
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := store.Save("https://example.com")
		if err != nil {
			t.Fatal(err)
		}
		if codes[code] {
			t.Fatalf("duplicate code generated: %s", code)
		}
		codes[code] = true
	}
}
