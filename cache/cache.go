package cache

import "github.com/hyperf/go-cache/error_code"

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
}

func (c *Cache[T]) Get(key string, defaultValue T) error {
	res, err := c.Driver.Get(key)
	if err != nil {
		return err
	}

	if res == "" {
		return error_code.NotFound
	}

	return c.Packer.UnPack(res, defaultValue)
}

func (c *Cache[T]) Has(key string) (bool, error) {
	return c.Driver.Has(key)
}

func (c *Cache[T]) Set(key string, value T, seconds int) error {
	res, err := c.Packer.Pack(value)
	if err != nil {
		return err
	}

	return c.Driver.Set(key, res, seconds)
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
