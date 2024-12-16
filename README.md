# Golang 缓存组件

## 安装

```shell
go get -u github.com/hyperf/go-cache
```

## 使用

```go
package main

import (
	"errors"
	"fmt"
	"github.com/hyperf/go-cache/cache"
	"github.com/hyperf/go-cache/driver/go_zero_driver"
	"github.com/hyperf/go-cache/error_code"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Foo struct {
	Id int `json:"id"`
}

func main() {
	fooCache := &cache.Cache[*Foo]{
		Packer: &cache.JsonPacker{},
		Driver: go_zero_driver.NewGoZeroDriver(redis.MustNewRedis(redis.RedisConf{
			Host: "127.0.0.1:6379",
			Type: "node",
		})),
		Prefix: "c:",
	}

	res := &Foo{}
	err := fooCache.Get("xx", res)
	fmt.Println(errors.Is(err, error_code.NotFound)) // true
	fmt.Println(res.Id)                              // 0

	_ = fooCache.Set("xx", &Foo{Id: 2}, 30)

	err = fooCache.Get("xx", res)
	fmt.Println(err)    // nil
	fmt.Println(res.Id) // 2
}


```