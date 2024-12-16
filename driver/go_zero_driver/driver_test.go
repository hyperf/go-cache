package go_zero_driver

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"testing"
	"time"
)

func TestGoZeroDriver(t *testing.T) {
	red := redis.MustNewRedis(redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
	})

	driver := NewGoZeroDriver(red)

	key := strconv.FormatInt(time.Now().Unix(), 10)
	exists, _ := driver.Has(key)
	assert.False(t, exists)

	str, err := driver.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, "", str)

	err = driver.Set(key, "HelloWorld", 60)
	assert.Nil(t, err)

	str, err = driver.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, "HelloWorld", str)
}
