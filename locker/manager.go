package locker

import "sync"

type LockerManager struct {
	lockers *sync.Map
}

func (c *LockerManager) Get(id string) *sync.Mutex {
	val, ok := c.lockers.Load(id)
	if ok {
		return val.(*sync.Mutex)
	}

	res := &sync.Mutex{}

	c.lockers.Store(id, res)

	return res
}

var Manager = &LockerManager{
	lockers: new(sync.Map),
}