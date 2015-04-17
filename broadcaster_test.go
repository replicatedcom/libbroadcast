package libbroadcast

import (
	"sync"
	"testing"
	"time"
)

func TestBroadcasterOnNextMessage(t *testing.T) {
	broadcaster := Global().GetBroadcaster("test-broadcaster")
	ch, err := broadcaster.CreateChan("test-chan")
	if err != nil {
		t.Fatalf("Error creating broadcast channel: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 1; i < 6; i++ {
		ch.OnNextMessage(func(i int) func(interface{}, error) {
			return func(msg interface{}, err error) {
				if msg.(int) != i {
					t.Errorf("Expecting %d got %d", i, msg)
				}
				wg.Done()
			}
		}(i))
	}

	for i := 1; i < 6; i++ {
		broadcaster.Send(i)
	}

	wg.Wait()

	broadcaster.RemoveChan("test-chan")
}

func TestBroadcasterWaitForNextMessage(t *testing.T) {
	broadcaster := Global().GetBroadcaster("test-broadcaster")
	ch, err := broadcaster.CreateChan("test-chan")
	if err != nil {
		t.Fatalf("Error creating broadcast channel: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 1; i < 6; i++ {
		go func(i int) {
			msg, _ := ch.WaitForNextMessage()
			if msg.(int) != i {
				t.Errorf("Expecting %d got %d", i, msg)
			}
			wg.Done()
		}(i)
	}

	for i := 1; i < 6; i++ {
		broadcaster.Send(i)
	}

	wg.Wait()

	broadcaster.RemoveChan("test-chan")
}

func TestBroadcasterRace(t *testing.T) {
	broadcaster := Global().GetBroadcaster("test-broadcaster")
	ch, err := broadcaster.CreateChan("test-chan")
	if err != nil {
		t.Fatalf("Error creating broadcast channel: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		for i := 1; i < 6; i++ {
			msg, _ := ch.WaitForNextMessage()
			if msg.(int) != i {
				t.Errorf("Expecting %d got %d", i, msg)
			}
			wg.Done()
		}
	}()

	go func() {
		for i := 1; i < 6; i++ {
			broadcaster.Send(i)
		}
	}()

	wg.Wait()

	broadcaster.RemoveChan("test-chan")
}

func TestBroadcasterSelect(t *testing.T) {
	broadcaster := Global().GetBroadcaster("test-broadcaster")
	ch, err := broadcaster.CreateChan("test-chan")
	if err != nil {
		t.Fatalf("Error creating broadcast channel: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	timer := time.NewTimer(time.Millisecond * 100)

	for i := 1; i < 3; i++ {
		go func(i int) {
			select {
			case msg := <-ch:
				if msg.(int) != i {
					t.Errorf("Expecting %d got %d", i, msg)
				}
				wg.Done()

			case <-timer.C:
				if 2 != i {
					t.Errorf("Expecting 1 got %d", i)
				}
				wg.Done()
			}

		}(i)
	}

	broadcaster.Send(1)

	wg.Wait()

	broadcaster.RemoveChan("test-chan")
}
