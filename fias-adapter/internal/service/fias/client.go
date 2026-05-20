package fias

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	GetAddressItemsPath = "/api/spas/v2.0/GetAddressItems"
)

type Config struct {
	BaseURL     string
	MasterToken string
}

type Client struct {
	baseURL     *url.URL
	masterToken string
	httpClient  *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse fias base url %q: %w", cfg.BaseURL, err)
	}
	return &Client{
		baseURL:     u,
		masterToken: cfg.MasterToken,
		httpClient:  &http.Client{},
	}, nil
}

type GetAddressItemsParams struct {
	AddressLevels []int  `json:"address_levels"`
	AddressType   int    `json:"address_type"`
	Path          string `json:"path"`
}

func (c *Client) GetAddressItems(p GetAddressItemsParams) (GetAddressItemsResponse, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return GetAddressItemsResponse{}, err
	}

	path := c.baseURL.JoinPath(GetAddressItemsPath).String()
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(jsonData))
	if err != nil {
		return GetAddressItemsResponse{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("master-token", c.masterToken)
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetAddressItemsResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return GetAddressItemsResponse{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var out GetAddressItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return GetAddressItemsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}

func (c *Client) buildURL(path string, q url.Values) string {
	u := c.baseURL.JoinPath(path)
	u.RawQuery = q.Encode()
	return u.String()
}
