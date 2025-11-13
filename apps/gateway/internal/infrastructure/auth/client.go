package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang-social-media/pkg/contracts/auth"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Register(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error) {
	var resp auth.RegisterResponse
	if err := c.do(ctx, http.MethodPost, "/auth/register", req, &resp); err != nil {
		return auth.RegisterResponse{}, err
	}
	return resp, nil
}

func (c *Client) Login(ctx context.Context, req auth.LoginRequest) (auth.LoginResponse, error) {
	var resp auth.LoginResponse
	if err := c.do(ctx, http.MethodPost, "/auth/login", req, &resp); err != nil {
		return auth.LoginResponse{}, err
	}
	return resp, nil
}

func (c *Client) GetProfile(ctx context.Context, userID string) (auth.ProfileResponse, error) {
	var resp auth.ProfileResponse
	path := fmt.Sprintf("/auth/profile/%s", userID)
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return auth.ProfileResponse{}, err
	}
	return resp, nil
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var apiErr auth.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("auth service error: %s", res.Status)
		}
		return fmt.Errorf("auth service error: %s", apiErr.Error)
	}

	if out == nil {
		return nil
	}

	return json.NewDecoder(res.Body).Decode(out)
}
