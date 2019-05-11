package overlay

import (
	"math/rand"
	"strconv"
	"sync"
)

type DirectMap interface {
	Bootstrap

	Chan(id string) DirectConnector
}

type ChanMap struct {
	*sync.RWMutex

	chanMap map[string]DirectConnector
}

func NewChanMap() DirectMap {
	return &ChanMap{
		RWMutex: new(sync.RWMutex),
		chanMap: make(map[string]DirectConnector),
	}
}

func (mp *ChanMap) NewId() string {
	mp.Lock()
	defer mp.Unlock()

	id := strconv.Itoa(len(mp.chanMap))
	mp.chanMap[id] = newDirectChan(id, mp)

	return id
}

func (mp *ChanMap) Join(id string) string {
	i := rand.Intn(len(mp.chanMap))
	for k := range mp.chanMap {
		if i == 0 {
			if k == id && len(mp.chanMap) > 1 {
				return mp.Join(id)
			}
			if k != id {
				return k
			}
		}
		i--
	}
	return ""
}

func (mp *ChanMap) Chan(id string) DirectConnector {
	mp.RLock()
	defer mp.RUnlock()

	if channel, ok := mp.chanMap[id]; ok {
		return channel
	} else {
		return nil
	}
}
