package libbroadcast

import (
	"errors"
)

type broadcastChan chan interface{}

var (
	ChannelClosedError = errors.New("channel closed")
)

func newBroadcastChan() broadcastChan {
	return make(broadcastChan)
}

func (c broadcastChan) WaitForNextMessage() (interface{}, bool) {
	msg, ok := <-c
	return msg, ok
}

func (c broadcastChan) OnNextMessage(onBroadcast func(interface{}, error)) {
	go func() {
		data, ok := <-c
		if ok {
			onBroadcast(data, nil)
		} else {
			onBroadcast(nil, ChannelClosedError)
		}
	}()
}
