package cache

import (
	"errors"
	"github.com/hyperf/go-cache/driver/go_zero_driver"
	"github.com/hyperf/go-cache/driver/null_driver"
	"github.com/hyperf/go-cache/error_code"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"testing"
	"time"
)

func TestCacheNullDriver(t *testing.T) {
	res := &Cache[*Foo]{Driver: null_driver.NullDriver{}, Packer: &JsonPacker{}}

	var result Foo
	err := res.Get("xx", &result)
	assert.Equal(t, 0, result.Id)
	assert.True(t, errors.Is(err, error_code.NotFound))

	result2 := &Foo{Id: 10}
	_ = res.Get("xx", result2)
	assert.Equal(t, 10, result2.Id)
}

func TestCacheGoZeroDriver(t *testing.T) {
	driver := go_zero_driver.NewGoZeroDriver(redis.MustNewRedis(redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
	}))

	fooCache := &Cache[*Foo]{Driver: driver, Packer: &JsonPacker{}}
	key := "c:" + strconv.FormatInt(time.Now().Unix(), 10)

	var result Foo
	err := fooCache.Get(key, &result)
	assert.Equal(t, 0, result.Id)
	assert.True(t, errors.Is(err, error_code.NotFound))

	res := &Foo{Id: 10}
	_ = fooCache.Get(key, res)
	assert.Equal(t, 10, res.Id)

	res.Id = 100
	exists, _ := fooCache.Has(key)
	assert.False(t, exists)

	err = fooCache.Set(key, res, 60)
	assert.Nil(t, err)

	exists, _ = fooCache.Has(key)
	assert.True(t, exists)

	err = fooCache.Get(key, &result)
	assert.Nil(t, err)

	assert.Equal(t, 100, result.Id)
}

func TestCache_Run(t *testing.T) {
	driver := go_zero_driver.NewGoZeroDriver(redis.MustNewRedis(redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
	}))

	fooCache := &Cache[*Foo]{Driver: driver, Packer: &JsonPacker{}, Prefix: "c2:"}
	key := strconv.FormatInt(time.Now().Unix(), 10)

	var result Foo
	_ = fooCache.Run(key, &result, 60, func(foo *Foo) error {
		foo.Id = 1
		return nil
	})

	assert.Equal(t, 1, result.Id)

	var result2 Foo
	_ = fooCache.Run(key, &result2, 60, func(foo *Foo) error {
		foo.Id = 2
		return nil
	})

	assert.Equal(t, 1, result2.Id)
}

func TestCache_RunAHead(t *testing.T) {
	driver := go_zero_driver.NewGoZeroDriver(redis.MustNewRedis(redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
	}))

	fooCache := &Cache[*Foo]{Driver: driver, Packer: &JsonPacker{}, Prefix: "c3:"}
	key := strconv.FormatInt(time.Now().Unix(), 10)

	data := &CacheAHead[*Foo]{
		Data:        &Foo{},
		ExpiredTime: 100,
	}

	_ = fooCache.RunAHead(key, data, func(foo *Foo) error {
		foo.Id = 1
		return nil
	})

	assert.Equal(t, 1, data.Data.Id)

	res := &CacheAHead[*Foo]{
		Data:        &Foo{},
		ExpiredTime: 100,
	}

	_ = fooCache.RunAHead(key, res, func(foo *Foo) error {
		foo.Id = 2
		return nil
	})

	assert.Equal(t, 1, data.Data.Id)
}
