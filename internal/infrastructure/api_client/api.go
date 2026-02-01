package api_client

import(
	"encoding/json"
	"github.com/sunwyx/currency-converter/internal/config"
	"net/http"
	"time"
	"context"
	"fmt"
)


type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg config.API) *Client {
	return &Client{
		apiKey: cfg.ApiKey,
		baseURL: cfg.CurrencyApiEndpoint,
		httpClient: &http.Client{
			Timeout: time.Second*10,
		},
	}
}

type latestResp struct {
	Data map[string]struct{
		Value float64 `json:"value"`
	} `json:"data"`
}

func(c *Client) GetRateAPI(ctx context.Context, base, target string) (float64, error) {
	url := fmt.Sprintf("%s/latest?base_currency=%s&currencies=%s", c.baseURL, base, target)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err!= nil {
		return 0, err
	}

	req.Header.Set("apikey", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("currency api returned %d", resp.StatusCode)
	}

	var r latestResp
	if err:= json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	v, ok:= r.Data[target]
	if !ok {
		return 0, fmt.Errorf("rate %s not found", target)
	}
	return v.Value, nil
}