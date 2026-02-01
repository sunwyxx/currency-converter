package service

import (
	"context"
	"fmt"
	"log/slog"
	"github.com/sunwyx/currency-converter/internal/infrastructure/redis_client"
	"github.com/sunwyx/currency-converter/internal/infrastructure/api_client"
)

type Converter struct {
	cache *redis_client.Cache
	api *api_client.Client
	Log *slog.Logger
}

func NewConverter(cache *redis_client.Cache, api *api_client.Client, log *slog.Logger) *Converter {
	return &Converter {
		cache: cache,
		api:   api,
		Log:   log,
	}
}

func (c *Converter) Convert(ctx context.Context, base, target string, amount float64) (float64, error){
	if amount <=0 {
		return 0, fmt.Errorf("amount must be > 0")
	}

	rate, found, err := c.cache.GetRate(ctx, base, target)
	if found != true && err != nil {
		c.Log.Error("redis did not return a value", slog.String("err", err.Error()))
	} else if found {
		return rate * amount, nil
	}

	rate, err = c.api.GetRateAPI(ctx, base, target)
	if err != nil {
		return 0, fmt.Errorf("api error: %w", err)
	}

	if err = c.cache.SetRate(ctx, base, target, rate, c.cache.Ttl); err != nil {
		c.Log.Error("failed to save to redis", slog.String("err", err.Error()))
	}

	return rate * amount, nil
}