package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Token returns the current valid access token, fetching or refreshing it if
// necessary. It is safe for concurrent use.
func (c *Client) Token(ctx context.Context) (string, error) {
	return c.accessToken(ctx)
}

func (c *Client) accessToken(ctx context.Context) (string, error) {
	if t, ok := c.cachedToken(); ok {
		return t, nil
	}
	c.mu <- struct{}{}
	defer func() { <-c.mu }()
	if t, ok := c.cachedToken(); ok {
		return t, nil
	}
	t, exp, err := c.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	c.token = t
	c.expiresAt = exp
	c.forceTokenRefresh = false
	if c.tokenStore != nil {
		_ = c.tokenStore.Set(t, exp)
	}
	return t, nil
}

func (c *Client) cachedToken() (string, bool) {
	if c.forceTokenRefresh {
		return "", false
	}
	if c.token == "" {
		if c.tokenStore != nil {
			if t, exp, ok := c.tokenStore.Get(); ok && time.Now().Before(exp.Add(-c.refreshMargin)) {
				c.token = t
				c.expiresAt = exp
				return t, true
			}
		}
		return "", false
	}
	if time.Now().After(c.expiresAt.Add(-c.refreshMargin)) {
		return "", false
	}
	return c.token, true
}

func (c *Client) fetchToken(ctx context.Context) (string, time.Time, error) {
	payload := map[string]string{
		"client_id":     c.clientID,
		"client_secret": c.secret,
		"grant_type":    "client_credentials",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ozon: marshal token request: %w", err)
	}
	u := c.baseURL + TokenPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(data))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ozon: create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ozon: token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		apiErr.StatusCode = resp.StatusCode
		return "", time.Time{}, &apiErr
	}
	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", time.Time{}, fmt.Errorf("ozon: decode token response: %w", err)
	}
	if tr.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("ozon: token response missing access_token")
	}
	ttl := tr.ExpiresIn
	if ttl <= 0 {
		ttl = 1800
	}
	return tr.AccessToken, time.Now().Add(time.Duration(ttl) * time.Second), nil
}
