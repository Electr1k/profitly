package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"yandex-adapter/internal/entity"
)

type Config struct {
	APIKey  string
	BaseURL string
	Lang    string
	Timeout time.Duration
}

type Client struct {
	apiKey     string
	baseURL    string
	lang       string
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
		lang:       cfg.Lang,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

type GeocodeParams struct {
	Lng     float64
	Lat     float64
	Results int
	Skip    int
	Kind    string
}

type GeocodeResult struct {
	Items []entity.GeoObject
	Total int
}

func (c *Client) Geocode(ctx context.Context, p GeocodeParams) (GeocodeResult, error) {
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("format", "json")
	q.Set("lang", c.lang)
	q.Set("geocode", strconv.FormatFloat(p.Lng, 'f', -1, 64)+","+strconv.FormatFloat(p.Lat, 'f', -1, 64))
	if p.Results > 0 {
		q.Set("results", strconv.Itoa(p.Results))
	}
	if p.Skip > 0 {
		q.Set("skip", strconv.Itoa(p.Skip))
	}
	if p.Kind != "" {
		q.Set("kind", p.Kind)
	}

	endpoint := c.baseURL + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return GeocodeResult{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GeocodeResult{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return GeocodeResult{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var out GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return GeocodeResult{}, fmt.Errorf("decode response: %w", err)
	}
	return toGeocodeResult(&out), nil
}
