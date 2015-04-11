package eventmanager

var globalBroadcastManager *BroadcastManager

type BroadcastManager struct {
	broadcasters map[string]*Broadcaster
}

func init() {
	globalBroadcastManager = &BroadcastManager{
		broadcasters: make(map[string]*Broadcaster),
	}
}

func GetBroadcastManager() *BroadcastManager {
	return globalBroadcastManager
}

func (broadcastManager *BroadcastManager) CreateBroadcaster(name string) *Broadcaster {
	// don't allow duplicated names to be created (since that will be a bad situation)
	if b, ok := broadcastManager.broadcasters[name]; ok {
		return b
	}

	broadcaster := newBroadcaster()
	broadcastManager.broadcasters[name] = broadcaster
	return broadcaster
}

func (broadcastManager *BroadcastManager) GetBroadcaster(name string) *Broadcaster {
	return broadcastManager.CreateBroadcaster(name)
}

func (broadcastManager *BroadcastManager) Destroy(name string) {
	delete(broadcastManager.broadcasters, name)
}
