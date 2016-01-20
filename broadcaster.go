package libbroadcast

import (
	"log"
	"time"
)

type Broadcaster struct {
	name      string
	listeners map[string]chan interface{}
}

func newBroadcaster(name string) *Broadcaster {
	return &Broadcaster{
		name:      name,
		listeners: make(map[string]chan interface{}),
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
	broadcaster.listeners[id] = listener
	return listener
}

func (broadcaster *Broadcaster) RemoveListener(id string) {
	delete(broadcaster.listeners, id)
}

func (broadcaster *Broadcaster) HasListeners() bool {
	return len(broadcaster.listeners) > 0
}

func (broadcaster *Broadcaster) Send(data interface{}) {
	for id, listener := range broadcaster.listeners {
		go func(i string, l chan interface{}) {
			select {
			case l <- data:
				return
			case <-time.After(30 * time.Minute):
				log.Printf("Broadcaster %q event send timed out for channel %s", broadcaster.name, i)
			}
		}(id, listener)
	}
}
