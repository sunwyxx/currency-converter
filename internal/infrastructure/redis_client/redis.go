package redis_client

import(
	"context" //todo
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sunwyx/currency-converter/internal/config"
	"strconv"
	"time"
)

type Cache struct {
	rdb *redis.Client
	Ttl time.Duration
}

func NewCache(ctx context.Context, cfg config.Redis) (*Cache, error) {

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
	}, nil
}

func rateKey(base, target string) string {
	return fmt.Sprintf("rate:%s:%s", base, target)
}

func (c *Cache) GetRate(ctx context.Context, base, target string) (float64, bool, error) {
	val, err:= c.rdb.Get(ctx, rateKey(base, target)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, false, nil
		}
		return 0, false, err
	}

	rate, err := strconv.ParseFloat(val, 64)
	if err != nil {
		_ = c.rdb.Del(ctx, rateKey(base, target)).Err()
		return 0, false, err
	}
	return rate, true, nil
}

func (c *Cache) SetRate(ctx context.Context, base, target string, rate float64, ttl time.Duration) error {
 	val := strconv.FormatFloat(rate, 'f', 2, 64)
    return c.rdb.Set(ctx, rateKey(base, target), val, ttl).Err()
}