package libbroadcast

type Broadcaster struct {
	listeners map[string]chan interface{}
}

func newBroadcaster() *Broadcaster {
	return &Broadcaster{
		listeners: make(map[string]chan interface{}),
	}
}

func (broadcaster *Broadcaster) Listen(id string, onBroadcast func(interface{})) {
	listener := make(chan interface{})
	broadcaster.listeners[id] = listener

	go func() {
		data := <-listener
		onBroadcast(data)
	}()
}

func (broadcaster *Broadcaster) ListenSynchronously(id string) interface{} {
	listener := make(chan interface{})
	broadcaster.listeners[id] = listener

	data := <-listener
	return data
}

func (broadcaster *Broadcaster) GetChan(id string) chan interface{} {
	listener := make(chan interface{})
	broadcaster.listeners[id] = listener
	return listener
}

func (broadcaster *Broadcaster) Send(data interface{}) {
	for _, listener := range broadcaster.listeners {
		listener <- data
	}
}

func (broadcaster *Broadcaster) RemoveListener(id string) {
	delete(broadcaster.listeners, id)
}
