package go_zero_driver

import "github.com/zeromicro/go-zero/core/stores/redis"

type GoZeroDriver struct {
	r *redis.Redis
}

func (g *GoZeroDriver) Get(key string) (string, error) {
	return g.r.Get(key)
}

func (g *GoZeroDriver) Has(key string) (bool, error) {
	return g.r.Exists(key)
}

func (g *GoZeroDriver) Set(key string, value string, seconds int) error {
	if seconds == 0 {
		return g.r.Set(key, value)
	}
	return g.r.Setex(key, value, seconds)
}
