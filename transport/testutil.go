package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
)

// MockHandler returns an http.HandlerFunc that replies with a fixed response.
// The optional `token` is checked against the Authorization header.
func MockHandler(statusCode int, body interface{}, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			if r.Header.Get(HeaderAuthorization) != TokenType+" "+token {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"invalid token"}`))
				return
			}
		}
		w.WriteHeader(statusCode)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}
}

// MockTokenHandler returns a handler for /api/client/token that issues the
// given token. One request is consumed per call unless `once` is set.
func MockTokenHandler(token string, ttlSeconds int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`
{"access_token":"` + token + `","expires_in":` + jsonNumber(ttlSeconds) + `,"token_type":"Bearer"}`))
	}
}

func jsonNumber(v int) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// NewTestClient builds a Client pointed at an httptest server. The server
// mounts the token endpoint and a fallback handler for any other path.
func NewTestClient(apiHandler http.HandlerFunc, token string) (*Client, *httptest.Server) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == TokenPath {
			MockTokenHandler(token, 3600)(w, r)
			return
		}
		apiHandler(w, r)
	}))
	return New("test-client", "test-secret", &Options{BaseURL: srv.URL}), srv
}

// CountedHandler counts how many times it was called, for asserting retries.
type CountedHandler struct {
	mu      sync.Mutex
	handler http.HandlerFunc
	count   int
}

func (c *CountedHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.Lock()
		c.count++
		c.mu.Unlock()
		c.handler(w, r)
	}
}

func (c *CountedHandler) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
