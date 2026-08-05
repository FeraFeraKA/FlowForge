package events

import "sync"

func (b *Bus) Subscribe(runID string) (<-chan Event, func()) {
	b.mu.Lock()
	if b.closed {
		channel := make(chan Event)
		close(channel)
		b.mu.Unlock()
		return channel, func() {}
	}
	b.nextID++
	subscriberID := b.nextID
	channel := make(chan Event, b.buffer)
	if b.subscribers[runID] == nil {
		b.subscribers[runID] = make(map[uint64]chan Event)
	}
	b.subscribers[runID][subscriberID] = channel
	b.mu.Unlock()

	var once sync.Once
	return channel, func() {
		once.Do(func() { b.unsubscribe(runID, subscriberID) })
	}
}

func (b *Bus) unsubscribe(runID string, subscriberID uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	channel, exists := b.subscribers[runID][subscriberID]
	if !exists {
		return
	}
	delete(b.subscribers[runID], subscriberID)
	close(channel)
}
