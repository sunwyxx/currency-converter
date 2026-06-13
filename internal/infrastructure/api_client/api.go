package api_client

import(
	"encoding/json"
	"github.com/sunwyx/currency-converter/internal/config"
	"log/slog"
	"net/http"
	"time"
	"context"
	"fmt"
)


type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	log *slog.Logger
}

func NewClient(cfg config.API, log *slog.Logger) *Client {
	return &Client{
		apiKey: cfg.ApiKey,
		baseURL: cfg.CurrencyApiEndpoint,
		httpClient: &http.Client{
			Timeout: time.Second*10,
		},
		log: log,
	}
}

type latestResp struct {
	Data map[string]struct{
		Code string `json:"code"`
		Value float64 `json:"value"`
	} `json:"data"`
}

func(c *Client) GetRateAPI(ctx context.Context, base string) (map[string]float64, error) {
	var resMap latestResp
	url := fmt.Sprintf("%s/latest?base_currency=%s", c.baseURL, base)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err!= nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("apikey", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("currency api returned %d", resp.StatusCode)
	}

	if err:= json.NewDecoder(resp.Body).Decode(&resMap); err != nil {
		return nil, fmt.Errorf("decodejson: %w", err)
	}
	return resMap.toMap(base), nil
}

func (r latestResp) toMap(base string) map[string]float64{
	res := make(map[string]float64)
	for _, r := range r.Data{
		res[fmt.Sprintf("rate:%s:%s", base, r.Code)] = r.Value
	}
	return res
}