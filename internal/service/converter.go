package service

import (
	"context"
	"fmt"
	"log/slog"
	"github.com/sunwyx/currency-converter/internal/infrastructure/redis_client"
	"github.com/sunwyx/currency-converter/internal/infrastructure/api_client"
	"strings"
)

type Converter struct {
	cache *redis_client.Cache
	api *api_client.Client
	log *slog.Logger
}

func NewConverter(cache *redis_client.Cache, api *api_client.Client, log *slog.Logger) *Converter {
	return &Converter {
		cache: cache,
		api:   api,
		log:   log,
	}
}

func (c *Converter) Convert(ctx context.Context, base string, targets []string, amount float64) (map[string]float64, error){
	if amount <=0 {
		return nil, fmt.Errorf("amount must be > 0")
	}
	targetsParams := make([]string, len(targets))
	for i, currency := range targets{
		targetsParams[i] = redis_client.RateKey(base, currency)
	}
	rates, err := c.cache.GetRate(ctx, base, targetsParams)
	if err != nil {
		c.log.Warn("error from GetRate", "err", err)
	}
	if len(targetsParams) != len(rates) {
		rates, err = c.api.GetRateAPI(ctx, base)
		if err != nil {
			return nil, fmt.Errorf("get rates from api: %w", err)
		}
	}
	err = c.cache.SetRate(ctx, rates)
	if err != nil {
		c.log.Warn("set rate", "err", err)
	}

	result := make(map[string]float64)
	for _, keyRateTarget := range targetsParams {
		parts := strings.Split(keyRateTarget, ":")
		part3 := parts[len(parts)-1]
		result[part3] = rates[keyRateTarget] * amount
	}
	c.log.Info("result", "result", rates)
	c.log.Info("result", "result", result)
	return result, nil
}