package events

import (
	"strconv"
	"sync"
	"time"
)

const DefaultSubscriberBuffer = 64

type Bus struct {
	mu          sync.Mutex
	buffer      int
	nextID      uint64
	sequences   map[string]uint64
	subscribers map[string]map[uint64]chan Event
	closed      bool
}

func NewBus(buffer int) *Bus {
	if buffer < 1 {
		buffer = DefaultSubscriberBuffer
	}
	return &Bus{
		buffer: buffer, sequences: make(map[string]uint64),
		subscribers: make(map[string]map[uint64]chan Event),
	}
}

func (b *Bus) Publish(event Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.sequences[event.RunID]++
	event.ID = strconv.FormatUint(b.sequences[event.RunID], 10)
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	for id, channel := range b.subscribers[event.RunID] {
		select {
		case channel <- event:
		default:
			delete(b.subscribers[event.RunID], id)
			close(channel)
		}
	}
}
