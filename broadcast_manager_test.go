package libbroadcast

import (
	"testing"
)

func TestBroadcastManager(t *testing.T) {
	broadcastManager := Global()

	if len(broadcastManager.broadcasters) != 0 {
		t.Fatal("Expected BroadcastManager.broadcasters to be empty")
	}

	b1 := broadcastManager.GetBroadcaster("test1")
	b2 := broadcastManager.GetBroadcaster("test2")

	if len(broadcastManager.broadcasters) != 2 {
		t.Fatal("Expected len(BroadcastManager.broadcasters) to be 2")
	}

	if b1 != broadcastManager.broadcasters["test1"] {
		t.Fatal("Expected item at key \"test1\" to equal b1")
	}

	if b2 != broadcastManager.broadcasters["test2"] {
		t.Fatal("Expected item at key \"test2\" to equal b2")
	}

	broadcastManager.Destroy("test1")

	if len(broadcastManager.broadcasters) != 1 {
		t.Fatal("Expected len(BroadcastManager.broadcasters) to be 1")
	}

	if _, ok := broadcastManager.broadcasters["test1"]; ok {
		t.Fatal("Expected key \"test1\" to not exist")
	}

	if _, ok := broadcastManager.broadcasters["test2"]; !ok {
		t.Fatal("Expected key \"test2\" to exist")
	}

	broadcastManager.Destroy("test2")

	if len(broadcastManager.broadcasters) != 0 {
		t.Fatal("Expected BroadcastManager.broadcasters to be empty")
	}
}
