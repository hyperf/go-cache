package cache

import (
	"errors"
	"github.com/hyperf/go-cache/error_code"
	"github.com/hyperf/go-cache/locker"
	"time"
)

type CacheInterface[T any] interface {
	Get(key string, defaultValue T) error
	Has(key string) (bool, error)
	// Set when the seconds is 0, the cache will be store forever.
	Set(key string, value T, seconds int) error
}

type DriverInterface interface {
	Get(key string) (string, error)
	Has(key string) (bool, error)
	Set(key string, value string, seconds int) error
}

type PackerInterface interface {
	Pack(data any) (string, error)
	UnPack(raw string, data any) error
}

type Cache[T any] struct {
	Driver DriverInterface
	Packer PackerInterface
	Prefix string
}

type CacheAHead[T any] struct {
	Data        T      `json:"data"`
	ExpiredAt   uint64 `json:"expired_at"`
	ExpiredTime int
}

func (a *CacheAHead[T]) fresh() *CacheAHead[T] {
	a.ExpiredAt = uint64(time.Now().Unix()) + uint64(a.ExpiredTime) - 60
	return a
}

func (c *Cache[T]) Key(key string) string {
	return c.Prefix + key
}

func (c *Cache[T]) Get(key string, defaultValue T) error {
	res, err := c.Driver.Get(c.Key(key))
	if err != nil {
		return err
	}

	if res == "" {
		return error_code.NotFound
	}

	return c.Packer.UnPack(res, defaultValue)
}

func (c *Cache[T]) Has(key string) (bool, error) {
	return c.Driver.Has(c.Key(key))
}

func (c *Cache[T]) Set(key string, value T, seconds int) error {
	res, err := c.Packer.Pack(value)
	if err != nil {
		return err
	}

	return c.Driver.Set(c.Key(key), res, seconds)
}

func (c *Cache[T]) Run(key string, defaultValue T, seconds int, fn func(T) error) error {
	key = c.Key(key)
	err := c.Get(key, defaultValue)
	if err != nil && !errors.Is(err, error_code.NotFound) {
		return err
	}

	if err == nil {
		return nil
	}

	err = fn(defaultValue)
	if err != nil {
		return err
	}

	err = c.Set(key, defaultValue, seconds)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache[T]) RunAHead(key string, data *CacheAHead[T], fn func(T) error) error {
	seconds := data.ExpiredTime
	key = c.Key(key)

	res, err := c.Driver.Get(key)
	if err != nil {
		return err
	}

	renew := func(d *CacheAHead[T], seconds int) error {
		err = fn(d.Data)
		if err != nil {
			return err
		}

		res, err = c.Packer.Pack(d.fresh())
		if err != nil {
			return err
		}

		return c.Driver.Set(key, res, seconds)
	}

	if res != "" {
		err = c.Packer.UnPack(res, data)
		if err != nil {
			return err
		}

		mu := locker.Manager.Get(key)
		locked := mu.TryLock()
		if locked {
			defer mu.Unlock()
			defer locker.Manager.Del(key)
		}

		if locked && data.ExpiredAt <= uint64(time.Now().Unix()) {
			return renew(data, seconds)
		}

		return nil
	}

	return renew(data, seconds)
}

func (c *Cache[T]) WithDriver(driver DriverInterface) *Cache[T] {
	c.Driver = driver
	return c
}

func (c *Cache[T]) WithPacker(p PackerInterface) *Cache[T] {
	c.Packer = p
	return c
}

func (c *Cache[T]) WithJsonPacker() *Cache[T] {
	c.Packer = &JsonPacker{}
	return c
}
