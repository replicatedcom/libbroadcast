package libbroadcast

import (
	"sync"
)

var globalBroadcastManager *BroadcastManager

type BroadcastManager struct {
	broadcasters map[string]*Broadcaster
	mu           sync.Mutex
}

func init() {
	globalBroadcastManager = &BroadcastManager{
		broadcasters: make(map[string]*Broadcaster),
	}
}

func Global() *BroadcastManager {
	return globalBroadcastManager
}

func (broadcastManager *BroadcastManager) GetBroadcaster(name string) *Broadcaster {
	broadcastManager.mu.Lock()
	defer broadcastManager.mu.Unlock()

	// don't allow duplicated names to be created (since that will be a bad situation)
	if b, ok := broadcastManager.broadcasters[name]; ok {
		return b
	}

	broadcaster := newBroadcaster()
	broadcastManager.broadcasters[name] = broadcaster
	return broadcaster
}

func (broadcastManager *BroadcastManager) CreateBroadcaster(name string) *Broadcaster {
	return broadcastManager.GetBroadcaster(name)
}

func (broadcastManager *BroadcastManager) Destroy(name string) {
	broadcastManager.mu.Lock()
	defer broadcastManager.mu.Unlock()

	delete(broadcastManager.broadcasters, name)
}
