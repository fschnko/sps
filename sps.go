package sps

import "sync"

// Databus thread safe messaging bus.
type Databus struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

// New creates a new instance of the Databus.
func New() *Databus {
	return &Databus{
		topics: make(map[string]*Topic),
	}
}

// Subscribe creates a new subscription for the topic.
// Creates a new topic if it doesn't exist.
func (b *Databus) Subscribe(topicName, subscriberName string) {
	t, ok := b.get(topicName)
	if !ok {
		t = b.create(topicName)
	}

	t.Subscribe(subscriberName)
}

// Unsubscribe removes the subscription from the topic.
func (b *Databus) Unsubscribe(topicName, subscriberName string) {
	t, ok := b.get(topicName)
	if ok {
		t.Unsubscribe(subscriberName)
	}
}

// Poll returns a list of unread messages.
// Returns an error if no topic or subscription found.
func (b *Databus) Poll(topicName, subscriberName string) ([][]byte, error) {
	t, ok := b.get(topicName)
	if !ok {
		return nil, ErrNoTopicFound
	}
	return t.Poll(subscriberName)
}

// Publish adds a new message to the topic pool.
func (b *Databus) Publish(topicName string, msg []byte) {
	t, ok := b.get(topicName)
	if ok {
		t.Publish(msg)
	}
}

// Close closes all topics and releases allocated resources.
func (b *Databus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for name, t := range b.topics {
		t.Close()
		delete(b.topics, name)
	}
}

func (b *Databus) get(topicName string) (*Topic, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	t, ok := b.topics[topicName]
	return t, ok
}

func (b *Databus) create(topicName string) *Topic {
	b.mu.Lock()
	defer b.mu.Unlock()

	t, ok := b.topics[topicName]
	if !ok {
		t = NewTopic()
		b.topics[topicName] = t
	}

	return t
}
