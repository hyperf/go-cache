package null_driver

import "github.com/hyperf/go-cache/cache"

type NullDriver struct{}

func (n NullDriver) Get(key string) (string, error) {
	return "", cache.NotFound
}

func (n NullDriver) Has(key string) (bool, error) {
	return false, nil
}

func (n NullDriver) Set(key string, value string, seconds int) error {
	return nil
}
