package concurrency

import (
	"sync"
)

type Pool[M any] struct {
	ObjectMap      map[string]M
	objectMapMutex sync.RWMutex
}

func (p *Pool[M]) Get(id string) (M, bool) {
	p.objectMapMutex.RLock()
	defer p.objectMapMutex.RUnlock()
	o, ok := p.ObjectMap[id]
	return o, ok
}

func (p *Pool[M]) GetOrCreate(id string, create func() (M, error)) (o M, err error) {
	var ok bool
	o, ok = p.Get(id)
	if !ok {
		p.objectMapMutex.Lock()
		defer p.objectMapMutex.Unlock()
		o, ok = p.ObjectMap[id]
		if !ok {
			o, err = create()
			if err != nil {
				return
			}
			p.ObjectMap[id] = o
		}
	}
	return
}

func (p *Pool[M]) Set(id string, o M) {
	p.objectMapMutex.Lock()
	defer p.objectMapMutex.Unlock()
	p.ObjectMap[id] = o
}

func (p *Pool[M]) Unset(id string) {
	p.objectMapMutex.Lock()
	defer p.objectMapMutex.Unlock()
	delete(p.ObjectMap, id)
}

func (p *Pool[M]) Map(f func(M)) {
	p.objectMapMutex.RLock()
	defer p.objectMapMutex.RUnlock()

	for _, o := range p.ObjectMap {
		f(o)
	}
}

func NewPool[M any]() *Pool[M] {
	return &Pool[M]{
		ObjectMap: make(map[string]M),
	}
}
