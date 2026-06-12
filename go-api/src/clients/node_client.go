package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jordanhuaman/go-api/src/models"
)

type NodeClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewNodeClient() *NodeClient {
	baseURL := os.Getenv("NODE_SERVICE_URL")

	timeout := 10 * time.Second
	if t := os.Getenv("NODE_SERVICE_TIMEOUT"); t != "" {
		if seconds, err := strconv.Atoi(t); err == nil && seconds > 0 {
			timeout = time.Duration(seconds) * time.Second
		}
	}

	return &NodeClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *NodeClient) CalculateStatistics(ctx context.Context, q, r models.Matrix2D) (*models.Statistics, error) {
	if c.baseURL == "" {
		return nil, nil
	}

	body := map[string]models.Matrix2D{"q": q, "r": r}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/matrix/statistics", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("node service call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node service status %d: %s", resp.StatusCode, string(respBody))
	}

	var stats models.Statistics
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return nil, fmt.Errorf("unmarshal statistics: %w", err)
	}

	return &stats, nil
}
