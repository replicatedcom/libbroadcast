package libbroadcast

import (
	"errors"
	"sync"
)

var (
	ChannelExistsError   = errors.New("broadcast channel exists")
	ChannelNotFoundError = errors.New("broadcast channel not found")
)

type Broadcaster struct {
	channels map[string]broadcastChan
	mu       sync.Mutex
}

func newBroadcaster() *Broadcaster {
	return &Broadcaster{
		channels: make(map[string]broadcastChan),
	}
}

func (broadcaster *Broadcaster) CreateChan(id string) (broadcastChan, error) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()

	if _, ok := broadcaster.channels[id]; ok {
		return nil, ChannelExistsError
	}
	broadcaster.channels[id] = newBroadcastChan()
	return broadcaster.channels[id], nil
}

func (broadcaster *Broadcaster) GetChan(id string) (broadcastChan, error) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()

	if ch, ok := broadcaster.channels[id]; ok {
		return ch, nil
	}
	return nil, ChannelNotFoundError
}

func (broadcaster *Broadcaster) RemoveChan(id string) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()

	if ch, ok := broadcaster.channels[id]; ok {
		close(ch)
		delete(broadcaster.channels, id)
	}
}

func (broadcaster *Broadcaster) Send(data interface{}) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()

	for _, ch := range broadcaster.channels {
		go func(c broadcastChan) {
			defer func() {
				// Prevent error sending on closed channel race condition
				recover()
			}()
			c <- data
		}(ch)
	}
}
