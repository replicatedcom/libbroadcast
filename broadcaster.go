package libbroadcast

import (
	"log"
	"sync"
	"time"
)

type Broadcaster struct {
	name      string
	listeners map[string]chan interface{}
	mutex     *sync.Mutex
}

func newBroadcaster(name string) *Broadcaster {
	return &Broadcaster{
		name:      name,
		listeners: make(map[string]chan interface{}),
		mutex:     &sync.Mutex{},
	}
}

func (broadcaster *Broadcaster) GetNextMessage(id string, onBroadcast func(interface{})) {
	listener := broadcaster.GetChan(id)
	go func() {
		data := <-listener
		broadcaster.RemoveListener(id)
		onBroadcast(data)
	}()
}

func (broadcaster *Broadcaster) GetChan(id string) chan interface{} {
	listener := make(chan interface{}, 1)

	broadcaster.mutex.Lock()
	defer broadcaster.mutex.Unlock()
	broadcaster.listeners[id] = listener

	return listener
}

func (broadcaster *Broadcaster) RemoveListener(id string) {
	broadcaster.mutex.Lock()
	defer broadcaster.mutex.Unlock()
	delete(broadcaster.listeners, id)
}

func (broadcaster *Broadcaster) HasListeners() bool {
	broadcaster.mutex.Lock()
	defer broadcaster.mutex.Unlock()
	return len(broadcaster.listeners) > 0
}

func (broadcaster *Broadcaster) Send(data interface{}) {
	broadcaster.mutex.Lock()
	defer broadcaster.mutex.Unlock()
	for id, listener := range broadcaster.listeners {
		// Don't do anything synchronously here becasue we are holding the lock
		go func(i string, l chan interface{}) {
			select {
			case l <- data:
				return
			case <-time.After(30 * time.Minute):
				if os.Getenv("DEBUG") != "" {
					log.Printf("Broadcaster %q event send timed out for channel %s", broadcaster.name, i)
				}
			}
		}(id, listener)
	}
}
