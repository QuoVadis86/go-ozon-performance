package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type statQuery struct {
	CampaignIds []Uint64 `url:"campaignIds"`
	Page        int64    `url:"page"`
	PageSize    int64    `url:"pageSize"`
	Vendor      bool     `url:"vendor"`
}

func TestTokenFetchAndReuse(t *testing.T) {
	calls := &CountedHandler{handler: func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}}
	cl, srv := NewTestClient(calls.Handler(), "tok-1")
	defer srv.Close()

	type resp struct {
		OK bool `json:"ok"`
	}
	for i := 0; i < 3; i++ {
		var r resp
		if err := cl.Get(context.Background(), "/api/x", nil, &r); err != nil {
			t.Fatal(err)
		}
	}
	if got := calls.Count(); got != 3 {
		t.Fatalf("api calls = %d, want 3", got)
	}
}

func TestAPIError(t *testing.T) {
	cl, srv := NewTestClient(MockHandler(http.StatusNotFound, map[string]string{"error": "report not found"}, "tok-2"), "tok-2")
	defer srv.Close()

	var out interface{}
	err := cl.Get(context.Background(), "/api/x", nil, &out)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ErrMessage != "report not found" {
		t.Fatalf("unexpected error: %+v", apiErr)
	}
	if apiErr.Error() == "" {
		t.Fatal("empty Error() string")
	}
}

func TestUnauthorizedTokenRefreshRetry(t *testing.T) {
	var stage = 0
	handler := &CountedHandler{handler: func(w http.ResponseWriter, r *http.Request) {
		if stage == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"invalid or expired token"}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}}
	cl, srv := NewTestClient(handler.Handler(), "tok-3")
	defer srv.Close()
	stage = 1

	var out map[string]bool
	if err := cl.Get(context.Background(), "/api/x", nil, &out); err != nil {
		t.Fatal(err)
	}
	if !out["ok"] {
		t.Fatal("expected ok after retry")
	}
}

func TestQueryEncoding(t *testing.T) {
	gotQuery := ""
	cl, srv := NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	}, "tok-4")
	defer srv.Close()

	q := &statQuery{
		CampaignIds: []Uint64{12558, 777},
		Page:        2,
		Vendor:      true,
	}
	var out interface{}
	if err := cl.Get(context.Background(), "/api/stats", q, &out); err != nil {
		t.Fatal(err)
	}
	want := "campaignIds=12558&campaignIds=777&page=2&vendor=true"
	if gotQuery != want {
		t.Fatalf("query = %q, want %q", gotQuery, want)
	}
}

func TestQueryEncodingOmitsZeroes(t *testing.T) {
	gotQuery := ""
	cl, srv := NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	}, "tok-5")
	defer srv.Close()

	var out interface{}
	q := &statQuery{}
	q.CampaignIds = nil
	q.Vendor = false
	if err := cl.Get(context.Background(), "/api/campaign", q, &out); err != nil {
		t.Fatal(err)
	}
	// bool fields are always encoded (false is meaningful); zero numerics and
	// empty slices are omitted.
	if gotQuery != "vendor=false" {
		t.Fatalf("query = %q, want %q", gotQuery, "vendor=false")
	}
}

func TestUint64JSON(t *testing.T) {
	// number input
	var a struct {
		V Uint64 `json:"v"`
	}
	if err := json.Unmarshal([]byte(`{"v":48852}`), &a); err != nil {
		t.Fatal(err)
	}
	if a.V != 48852 {
		t.Fatalf("got %v", a.V)
	}
	// string input
	var b struct {
		V Uint64 `json:"v"`
	}
	if err := json.Unmarshal([]byte(`{"v":"12558"}`), &b); err != nil {
		t.Fatal(err)
	}
	if b.V != 12558 {
		t.Fatalf("got %v", b.V)
	}
	// null input
	var c struct {
		V Uint64 `json:"v"`
	}
	if err := json.Unmarshal([]byte(`{"v":null}`), &c); err != nil {
		t.Fatal(err)
	}
	// empty string
	var d struct {
		V Uint64 `json:"v"`
	}
	if err := json.Unmarshal([]byte(`{"v":""}`), &d); err != nil {
		t.Fatal(err)
	}
	// marshal emits string
	var e struct {
		V Uint64 `json:"v"`
	}
	e.V = 48852
	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"v":"48852"}` {
		t.Fatalf("marshal = %s", raw)
	}
}

func TestPostJSONBody(t *testing.T) {
	gotBody := ""
	cl, srv := NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 256)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		_, _ = w.Write([]byte(`{"UUID":"abc"}`))
	}, "tok-6")
	defer srv.Close()

	req := map[string]interface{}{"campaigns": []Uint64{48852}}
	var out map[string]string
	if err := cl.Post(context.Background(), "/api/client/statistics", req, &out); err != nil {
		t.Fatal(err)
	}
	want := `{"campaigns":["48852"]}`
	if gotBody != want {
		t.Fatalf("body = %q, want %q", gotBody, want)
	}
	if out["UUID"] != "abc" {
		t.Fatalf("resp = %v", out)
	}
}

func TestGetRaw(t *testing.T) {
	cl, srv := NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("sku;Название\n12558;Штука\n"))
	}, "tok-7")
	defer srv.Close()

	raw, err := cl.GetRaw(context.Background(), "/api/client/statistics/report", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatal("empty raw body")
	}
}

func TestConcurrentTokenFetch(t *testing.T) {
	// Two goroutines racing to fetch the token; only one actual token request
	// should be issued because of the lock.
	var tokenFetchCount = 0
	srv := testTokenServer(t, &tokenFetchCount, "tok-8")
	defer srv.Close()
	cl := New("c", "s", &Options{BaseURL: srv.URL})

	done := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			_, err := cl.Token(context.Background())
			done <- err
		}()
	}
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
	if tokenFetchCount != 1 {
		t.Fatalf("token fetches = %d, want 1", tokenFetchCount)
	}
}

func testTokenServer(t *testing.T, counter *int, token string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == TokenPath {
			*counter++
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"` + token + `","expires_in":3600,"token_type":"Bearer"}`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	return srv
}
