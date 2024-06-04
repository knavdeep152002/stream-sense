package concurrency

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
	"sync"
	"time"
)

type Publisher[M any] struct {
	observerPool *Pool[chan M]
	producer     chan M
	running      bool
}

func randomString() (h string, err error) {
	bytes := make([]byte, 8)
	if _, err = rand.Read(bytes); err != nil {
		return
	}
	timeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeBytes, uint64(time.Now().UnixNano()))

	h = hex.EncodeToString(bytes) + hex.EncodeToString(timeBytes)
	return
}

func (p *Publisher[M]) RegisterObserver(observer chan M) (id string, err error) {
	id, err = randomString()
	if err != nil {
		return
	}
	log.Println("register", id, observer)
	p.observerPool.Set(id, observer)
	return
}

func (p *Publisher[M]) DeregisterObserver(id string) {
	p.observerPool.Unset(id)
}

func (p *Publisher[M]) runPublishingLoop() {
	log.Println("running")
	for p.running {
		log.Println("running")
		if data, ok := <-p.producer; ok {
			var wg sync.WaitGroup
			p.observerPool.Map(func(o chan M) {
				log.Println("publishing", data, o)
				wg.Add(1)
				go func() {
					defer wg.Done()
					select {
					case o <- data:
					case <-time.After(time.Second * 1):
					}
				}()
			})
			wg.Wait()
		} else {
			p.running = false
			log.Println("closing")
		}
	}
	log.Println("closed")
}

func NewPublisher[M any](producer chan M) (p *Publisher[M]) {
	log.Println("new publisher")
	p = &Publisher[M]{
		observerPool: NewPool[chan M](),
		producer:     producer,
		running:      true,
	}
	log.Println("running")
	go p.runPublishingLoop()
	return
}
