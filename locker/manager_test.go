package locker

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type Foo struct {
	Id   int
	Name string
}

func TestManager_Lock(t *testing.T) {
	res := &Foo{}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		mu := Manager.Get("test")
		mu.Lock()
		time.Sleep(1 * time.Second)
		res.Name += "2"
		mu.Unlock()
	}()

	time.Sleep(500 * time.Millisecond)

	mu := Manager.Get("test")
	mu.Lock()
	res.Name += "1"
	mu.Unlock()

	wg.Wait()

	assert.Equal(t, "21", res.Name)
}

func TestManager_TryLock(t *testing.T) {
	res := &Foo{}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		mu := Manager.Get("test2")

		mu.Lock()
		time.Sleep(1 * time.Second)
		res.Name += "2"
		mu.Unlock()
	}()

	time.Sleep(500 * time.Millisecond)

	mu := Manager.Get("test2")
	locked := mu.TryLock()
	if locked {
		mu.Unlock()
	}

	assert.False(t, locked)

	wg.Wait()

	assert.Equal(t, "2", res.Name)
}
