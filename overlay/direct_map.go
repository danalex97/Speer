package overlay

import (
	"math/rand"
	"strconv"
	"sync"
)

type DirectMap struct {
	*sync.RWMutex

	chanMap map[string]*DirectChan
	chanCtr int
}

func NewDirectMap() *DirectMap {
	return &DirectMap{
		RWMutex: new(sync.RWMutex),
		chanMap: make(map[string]*DirectChan),
		chanCtr: 0,
	}
}

func (mp *DirectMap) NewId() string {
	mp.Lock()
	defer mp.Unlock()

	id := strconv.Itoa(mp.chanCtr)
	mp.chanCtr += 1

	return id
}

func (mp *DirectMap) Join(id string) string {
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

func (mp *DirectMap) Chan(id string) *DirectChan {
	mp.RLock()
	defer mp.RUnlock()

	if channel, ok := mp.chanMap[id]; ok {
		return channel
	} else {
		return nil
	}
}

func (mp *DirectMap) RegisterChan(id string, channel *DirectChan) {
	mp.Lock()
	defer mp.Unlock()

	mp.chanMap[id] = channel
}
