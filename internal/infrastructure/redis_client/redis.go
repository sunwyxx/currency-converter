package redis_client

import(
	"context" 
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sunwyx/currency-converter/internal/config"
	"log/slog"
	"strconv"
	"time"
)

type Cache struct {
	rdb *redis.Client
	Ttl time.Duration
	log *slog.Logger
}

func NewCache(ctx context.Context, cfg config.Redis, log *slog.Logger) (*Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
		DB: cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Cache{
     rdb: client,
	 Ttl: time.Duration(cfg.TTL) * time.Hour,
	 log: log,
	}, nil
}

func RateKey(base, target string) string {
	return fmt.Sprintf("rate:%s:%s", base, target)
}

func (c *Cache) GetRate(ctx context.Context, base string, target []string) (map[string]float64, error) {
	m := make(map[string]float64)

	rates, err:= c.rdb.MGet(ctx, target...).Result()

	if err != nil {
		return m, fmt.Errorf("mget returned: %w", err)
	}

	for i, valRdb := range rates {
		if valRdb == nil {
			continue
		}
		strValRdb, ok := valRdb.(string)
		if !ok {
			_ = c.rdb.Del(ctx, RateKey(base, target[i]))
			c.log.Error("error in string type cast")
			continue
		}
		rate, err := strconv.ParseFloat(strValRdb, 64)
		if err != nil {
			_ = c.rdb.Del(ctx, RateKey(base, target[i])).Err()
			c.log.Error("error in float type cast")
			continue
		}
		m[target[i]] = rate
	}
	return m, nil
}

func (c *Cache) SetRate(ctx context.Context, rates map[string]float64) error {
	 pipe := c.rdb.Pipeline()
    for keysCurrency, rate := range rates {
        pipe.Set(ctx, keysCurrency, rate, c.Ttl)
    }
    _, err := pipe.Exec(ctx)
    if err != nil {
        return err
    }
    return nil
}