package libbroadcast

import (
	"log"
)

type Broadcaster struct {
	listeners map[string]chan interface{}
}

func newBroadcaster() *Broadcaster {
	return &Broadcaster{
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

func (broadcaster *Broadcaster) Send(data interface{}) {
	log.Printf("Broadcasting data to all listeners (there are %d)", len(broadcaster.listeners))

	for id, listener := range broadcaster.listeners {
		go func(i string, l chan interface{}) {
			log.Printf(" -> Delivering to listener %s", i)
			l <- data
			log.Printf(" -> Delivered to listener %s", i)
		}(id, listener)
	}
}
