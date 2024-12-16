package cache

import (
	"errors"
	"github.com/hyperf/go-cache/driver/null_driver"
	"github.com/hyperf/go-cache/error_code"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache_Get(t *testing.T) {
	res := &Cache[*Foo]{Driver: null_driver.NullDriver{}, Packer: &JsonPacker{}}

	var result Foo
	err := res.Get("xx", &result)
	assert.Equal(t, 0, result.Id)
	assert.True(t, errors.Is(err, error_code.NotFound))

	result2 := &Foo{Id: 10}
	_ = res.Get("xx", result2)
	assert.Equal(t, 10, result2.Id)
}
