package custom_logic

import (
	"context"
	"fmt"

	"zaiboot/segmentIO.tests/internal/configs"

	"github.com/redis/go-redis/v9"
)

type Dedupe_storage interface {
	connect()
	seen(ctx context.Context, key string, config *configs.Config) (bool, error)
	close()
}

type Redis_dedupe struct {
	c *redis.Client
}

func (r *Redis_dedupe) close() {
	r.c.Close()
}
func (r *Redis_dedupe) connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r.c = client

}

func (c Redis_dedupe) seen(ctx context.Context, key string, config *configs.Config) (bool, error) {
	k := fmt.Sprintf("%s:%s:%s", config.APPLICATIONNAME, config.KAFKA_TOPIC, key)
	ok, err := c.c.SetNX(ctx, k, "true", config.DEDUP_EXPIRATION).Result()
	return ok, err
}
