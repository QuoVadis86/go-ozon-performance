package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL is the Ozon Performance API base URL.
	DefaultBaseURL = "https://api-performance.ozon.ru:443"
	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 30 * time.Second

	// TokenPath is the OAuth2 token endpoint path.
	TokenPath = "/api/client/token"

	// HeaderAuthorization is the Authorization header name.
	HeaderAuthorization = "Authorization"
	// TokenType is the bearer token type prefix.
	TokenType = "Bearer"

	// Common API error messages (Performance API returns rpcStatus-style errors).
	ErrUnauthorized     = "unauthorized"
	ErrInvalidToken     = "invalid or expired token"
	ErrReportNotFound   = "report not found"
	ErrCampaignNotFound = "campaign not found"
)

// Options configures the transport Client.
type Options struct {
	// BaseURL overrides the default API base URL (mainly for testing).
	BaseURL string
	// Timeout overrides the default HTTP request timeout.
	Timeout time.Duration
	// HTTPClient overrides the underlying *http.Client.
	HTTPClient *http.Client
	// TokenStore allows persisting the access token across restarts.
	// When nil, the token is kept in memory only.
	TokenStore TokenStore
	// TokenRefreshMargin is how long before expiry the token is considered
	// stale and must be refreshed. Defaults to 60 seconds.
	TokenRefreshMargin time.Duration
}

// TokenStore persists the API access token.
type TokenStore interface {
	Get() (token string, expiresAt time.Time, ok bool)
	Set(token string, expiresAt time.Time) error
}

// Client is a token-authenticated HTTP client for the Ozon Performance API.
// Credentials are the `client_id` and `client_secret` obtained from the
// seller personal cabinet (Settings → API keys → Performance API).
//
// The client fetches and caches the access token automatically and refreshes
// it transparently once expired.
type Client struct {
	httpClient *http.Client
	baseURL    string
	clientID   string
	secret     string

	mu                chan struct{} // 1-buffered mutex for token refresh
	token             string
	expiresAt         time.Time
	refreshMargin     time.Duration
	tokenStore        TokenStore
	forceTokenRefresh bool
}

// New creates a Performance API client. `clientID` and `clientSecret` are
// the Performance API credentials. `opts` may be nil.
func New(clientID, clientSecret string, opts *Options) *Client {
	baseURL := DefaultBaseURL
	timeout := DefaultTimeout
	margin := 60 * time.Second
	var store TokenStore
	if opts != nil {
		if opts.BaseURL != "" {
			baseURL = opts.BaseURL
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		if opts.TokenRefreshMargin > 0 {
			margin = opts.TokenRefreshMargin
		}
		store = opts.TokenStore
	}
	hc := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     60 * time.Second,
		},
	}
	if opts != nil && opts.HTTPClient != nil {
		hc = opts.HTTPClient
	}
	return &Client{
		httpClient:    hc,
		baseURL:       baseURL,
		clientID:      clientID,
		secret:        clientSecret,
		mu:            make(chan struct{}, 1),
		refreshMargin: margin,
		tokenStore:    store,
	}
}

// APIError represents a non-2xx response from the Performance API.
type APIError struct {
	StatusCode int    `json:"-"`
	ErrMessage string `json:"error"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

func (e *APIError) Error() string {
	if e.ErrMessage != "" {
		return fmt.Sprintf("ozon performance API error (status=%d): %s", e.StatusCode, e.ErrMessage)
	}
	if e.Message != "" {
		return fmt.Sprintf("ozon performance API error (status=%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("ozon performance API error (status=%d)", e.StatusCode)
}

// Post sends a POST request with a JSON body.
func (c *Client) Post(ctx context.Context, path string, body, resp interface{}) error {
	_, err := c.do(ctx, http.MethodPost, path, nil, body, resp)
	return err
}

// Get sends a GET request. `query` is an optional struct tagged with `url:"..."`.
func (c *Client) Get(ctx context.Context, path string, query, resp interface{}) error {
	_, err := c.do(ctx, http.MethodGet, path, query, nil, resp)
	return err
}

// Put sends a PUT request with a JSON body.
func (c *Client) Put(ctx context.Context, path string, body, resp interface{}) error {
	_, err := c.do(ctx, http.MethodPut, path, nil, body, resp)
	return err
}

// Patch sends a PATCH request with a JSON body.
func (c *Client) Patch(ctx context.Context, path string, body, resp interface{}) error {
	_, err := c.do(ctx, http.MethodPatch, path, nil, body, resp)
	return err
}

// Delete sends a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, query, resp interface{}) error {
	_, err := c.do(ctx, http.MethodDelete, path, query, nil, resp)
	return err
}

// GetRaw sends a GET request and returns the raw response body (used for
// CSV/ZIP report downloads).
func (c *Client) GetRaw(ctx context.Context, path string, query interface{}) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, query, nil, nil)
}

// PostRaw sends a POST request and returns the raw response body.
func (c *Client) PostRaw(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// download returns raw bytes no matter the method; used internally by
// report-style generated methods.
func (c *Client) download(ctx context.Context, method, path string, query, body interface{}) ([]byte, error) {
	return c.do(ctx, method, path, query, body, nil)
}

func (c *Client) do(ctx context.Context, method, path string, query, body, resp interface{}) ([]byte, error) {
	// Retry once on authentication failure after forcing a token refresh.
	for attempt := 0; attempt < 2; attempt++ {
		out, err := c.doOnce(ctx, method, path, query, body, resp)
		if err == nil {
			return out, nil
		}
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusUnauthorized && attempt == 0 {
			c.forceTokenRefresh = true
			continue
		}
		return nil, err
	}
	return nil, nil
}

func (c *Client) doOnce(ctx context.Context, method, path string, query, body, resp interface{}) ([]byte, error) {
	token, err := c.accessToken(ctx)
	if err != nil {
		return nil, err
	}

	fullURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("ozon: build url: %w", err)
	}
	if query != nil {
		q, err := encodeQuery(query)
		if err != nil {
			return nil, err
		}
		if len(q) > 0 {
			fullURL += "?" + q.Encode()
		}
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("ozon: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("ozon: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderAuthorization, TokenType+" "+token)

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ozon: do request: %w", err)
	}
	defer httpResp.Body.Close()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("ozon: read response: %w", err)
	}

	if httpResp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: httpResp.StatusCode}
		_ = json.Unmarshal(raw, apiErr)
		return nil, apiErr
	}

	if resp != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, resp); err != nil {
			return nil, fmt.Errorf("ozon: unmarshal response: %w", err)
		}
		return raw, nil
	}
	return raw, nil
}

// encodeQuery builds url.Values from a struct using `url` field tags.
//
// Encoding rules:
//   - zero strings / numbers and empty slices are omitted;
//   - booleans are always included (`false` is meaningful);
//   - slices produce repeated query parameters.
func encodeQuery(v interface{}) (url.Values, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ozon: query must be a struct, got %T", v)
	}
	vals := url.Values{}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		name, ok := sf.Tag.Lookup("url")
		if !ok || name == "-" {
			continue
		}
		fv := rv.Field(i)
		if err := appendQueryValue(vals, name, fv); err != nil {
			return nil, err
		}
	}
	return vals, nil
}

func appendQueryValue(vals url.Values, name string, fv reflect.Value) error {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
		if fv.Len() == 0 {
			return nil
		}
		for i := 0; i < fv.Len(); i++ {
			if err := appendScalar(vals, name, fv.Index(i), true); err != nil {
				return err
			}
		}
		return nil
	}
	return appendScalar(vals, name, fv, false)
}

func appendScalar(vals url.Values, name string, fv reflect.Value, isElem bool) error {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	switch fv.Kind() {
	case reflect.String:
		s := fv.String()
		if !isElem && s == "" {
			return nil
		}
		vals.Add(name, s)
	case reflect.Bool:
		vals.Add(name, strconv.FormatBool(fv.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := fv.Int()
		if !isElem && n == 0 {
			return nil
		}
		vals.Add(name, strconv.FormatInt(n, 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := fv.Uint()
		if !isElem && n == 0 {
			return nil
		}
		vals.Add(name, strconv.FormatUint(n, 10))
	case reflect.Float32, reflect.Float64:
		f := fv.Float()
		if !isElem && f == 0 {
			return nil
		}
		vals.Add(name, strconv.FormatFloat(f, 'f', -1, 64))
	default:
		if str, ok := fv.Interface().(fmt.Stringer); ok {
			s := str.String()
			if !isElem && s == "" {
				return nil
			}
			vals.Add(name, s)
			return nil
		}
		return fmt.Errorf("ozon: unsupported query type %s for field %s", fv.Kind(), name)
	}
	return nil
}
