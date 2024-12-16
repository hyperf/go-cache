package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Foo struct {
	Id int `json:"id"`
}

func TestNewFamilyCode(t *testing.T) {
	packer := &JsonPacker{}
	res, _ := packer.Pack(&Foo{Id: 1})
	assert.Equal(t, "{\"id\":1}", res)

	var foo Foo
	_ = packer.UnPack(res, &foo)
	assert.Equal(t, Foo{Id: 1}, foo)
}
