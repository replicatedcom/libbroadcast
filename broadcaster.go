package libbroadcast

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

func (broadcaster *Broadcaster) HasListeners() bool {
	return len(broadcaster.listeners) > 0
}

func (broadcaster *Broadcaster) Send(data interface{}) {
	for id, listener := range broadcaster.listeners {
		go func(i string, l chan interface{}) {
			l <- data
		}(id, listener)
	}
}
