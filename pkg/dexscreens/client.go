package dexscreens

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
)

type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (c *Client) GetTopCoins(ctx context.Context, limit int) ([]models.Coin, error) {
	url := fmt.Sprintf("%s/tokens/top?limit=%d", c.baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Tokens []models.Coin `json:"tokens"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Tokens, nil
}

func (c *Client) GetCoinByID(ctx context.Context, id string) (*models.Coin, error) {
	url := fmt.Sprintf("%s/tokens/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var response struct {
		Token *models.Coin `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Token, nil
}

func (c *Client) handleResponse(resp *http.Response, v interface{}) error {
	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("http status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("api error: %s", errResp.Error)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}
