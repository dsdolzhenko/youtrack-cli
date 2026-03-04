package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	BaseURL string
	Token   string
}

type Client struct {
	cfg  Config
	http *http.Client
}

func New(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Get(path string, params url.Values) (*http.Response, error) {
	base := strings.TrimRight(c.cfg.BaseURL, "/")
	rawURL := base + "/" + strings.TrimLeft(path, "/")

	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("client: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: do request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("client: unexpected status %d %s for %s", resp.StatusCode, http.StatusText(resp.StatusCode), rawURL)
	}

	return resp, nil
}

func (c *Client) Post(path string, body io.Reader) (*http.Response, error) {
	base := strings.TrimRight(c.cfg.BaseURL, "/")
	rawURL := base + "/" + strings.TrimLeft(path, "/")

	req, err := http.NewRequest(http.MethodPost, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("client: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: do request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("client: unexpected status %d %s for %s", resp.StatusCode, http.StatusText(resp.StatusCode), rawURL)
	}

	return resp, nil
}
