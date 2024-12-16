package cache

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
	driver DriverInterface
	packer PackerInterface
}

func (c *Cache[T]) Get(key string, defaultValue T) error {
	res, err := c.driver.Get(key)
	if err != nil {
		return err
	}

	return c.packer.UnPack(res, defaultValue)
}

func (c *Cache[T]) Has(key string) (bool, error) {
	return c.driver.Has(key)
}

func (c Cache[T]) Set(key string, value T, seconds int) error {
	res, err := c.packer.Pack(value)
	if err != nil {
		return err
	}

	return c.driver.Set(key, res, seconds)
}

func NewCache[T any](driver DriverInterface, packer PackerInterface) *Cache[T] {
	return &Cache[T]{driver: driver, packer: packer}
}

func NewJsonCache[T any](driver DriverInterface) *Cache[T] {
	return &Cache[T]{driver: driver, packer: &JsonPacker{}}
}
