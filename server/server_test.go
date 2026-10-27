package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/partite-ai/example-ticket-pipeline-shortlink/store"
)

func TestCreateAndRedirect(t *testing.T) {
	h := New(store.New())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "http://short.test/links", strings.NewReader(`{"url":"https://example.com/a?b=1#c"}`)))
	if rec.Code != 201 {
		t.Fatalf("status %d %s", rec.Code, rec.Body)
	}
	var body struct {
		Code     string `json:"code"`
		ShortURL string `json:"short_url"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.ShortURL != "http://short.test/"+body.Code {
		t.Fatalf("short_url %q", body.ShortURL)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/"+body.Code, nil))
	if rec.Code != 302 || rec.Header().Get("Location") != "https://example.com/a?b=1#c" {
		t.Fatalf("%d %q", rec.Code, rec.Header().Get("Location"))
	}
}

func TestCreateRejectsBadJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	New(store.New()).ServeHTTP(rec, httptest.NewRequest("POST", "/links", strings.NewReader(`{`)))
	if rec.Code != 400 || !strings.Contains(rec.Body.String(), `"error"`) {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestCreateRejectsBadURL(t *testing.T) {
	rec := httptest.NewRecorder()
	New(store.New()).ServeHTTP(rec, httptest.NewRequest("POST", "/links", strings.NewReader(`{"url":""}`)))
	if rec.Code != 400 {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestCreateRejectsSchemelessURL(t *testing.T) {
	rec := httptest.NewRecorder()
	New(store.New()).ServeHTTP(rec, httptest.NewRequest("POST", "/links", strings.NewReader(`{"url":"example.com"}`)))
	if rec.Code != 400 || !strings.Contains(rec.Body.String(), "http or https") {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestCreateRejectsTooLongURL(t *testing.T) {
	body := `{"url":"https://example.com/` + strings.Repeat("a", 2100) + `"}`
	rec := httptest.NewRecorder()
	New(store.New()).ServeHTTP(rec, httptest.NewRequest("POST", "/links", strings.NewReader(body)))
	if rec.Code != 400 || !strings.Contains(rec.Body.String(), `"error"`) {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestRedirectUnknown(t *testing.T) {
	rec := httptest.NewRecorder()
	New(store.New()).ServeHTTP(rec, httptest.NewRequest("GET", "/nothere", nil))
	if rec.Code != 404 {
		t.Fatalf("%d", rec.Code)
	}
}

func TestRedirectExpired(t *testing.T) {
	h := New(store.New(store.WithTTL(20 * time.Millisecond)))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/links", strings.NewReader(`{"url":"https://example.com"}`)))
	if rec.Code != 201 {
		t.Fatalf("status %d %s", rec.Code, rec.Body)
	}
	var body struct {
		Code string `json:"code"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	time.Sleep(50 * time.Millisecond)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/"+body.Code, nil))
	if rec.Code != 404 {
		t.Fatalf("%d", rec.Code)
	}
}

func TestNoServerErrors(t *testing.T) {
	h := New(store.New())
	for _, r := range []*http.Request{
		httptest.NewRequest("GET", "/links", nil),
		httptest.NewRequest("GET", "/a/b", nil),
		httptest.NewRequest("PUT", "/links", nil),
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		if rec.Code >= 500 {
			t.Fatalf("%s %s -> %d", r.Method, r.URL, rec.Code)
		}
	}
}
