package sps

import (
	"time"
)

const (
	jobsBufferLen = 100
	cleanPeriod   = time.Minute
)

// Topic represents a single thread safe topic.
type Topic struct {
	jobs chan func()

	sequence int
	messages [][]byte
	// Subscribers has a form map[name]{Last polled sequence}
	subscribers map[string]int

	done chan struct{}
}

// NewTopic creates a new instance of the topic.
func NewTopic() *Topic {
	t := &Topic{
		jobs:        make(chan func(), jobsBufferLen),
		subscribers: make(map[string]int),
		done:        make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(cleanPeriod)
		for {
			select {
			case j, ok := <-t.jobs:
				if ok {
					j()
				}
			case <-ticker.C:
				t.clean()
			case _, ok := <-t.done:
				if !ok {
					ticker.Stop()
					// Finish all enqueued jobs to prevent routines lock.
					for j := range t.jobs {
						j()
					}
					return
				}
			}

		}
	}()

	return t

}

// Subscribe creates a new subscription for the topic.
func (t *Topic) Subscribe(name string) {
	if t.isClosed() {
		return
	}

	t.jobs <- func() {
		if _, ok := t.subscribers[name]; ok {
			return
		}
		t.subscribers[name] = t.sequence
	}

}

// Unsubscribe removes subscription from the topic.
func (t *Topic) Unsubscribe(name string) {
	if t.isClosed() {
		return
	}

	t.jobs <- func() {
		delete(t.subscribers, name)
	}
}

// Publish adds new message to the topic pool.
func (t *Topic) Publish(msg []byte) {
	if t.isClosed() {
		return
	}

	t.jobs <- func() {
		t.messages = append(t.messages, msg)
		t.sequence++
	}
}

// Poll returns list of unread messages or error if no subscription found.
func (t *Topic) Poll(name string) ([][]byte, error) {
	if t.isClosed() {
		return nil, nil
	}

	type response struct {
		data [][]byte
		err  error
	}

	result := make(chan response)

	t.jobs <- func() {
		var res response
		defer func() { result <- res }()

		lastID, ok := t.subscribers[name]
		if !ok {
			res.err = ErrNoSubscriptionFound
			return
		}

		if len(t.messages) == 0 {
			return
		}

		length := (t.sequence - lastID)
		offset := len(t.messages) - length

		if offset < 0 {
			// Negative offset means that sequence logic is broken
			// and system has an inconsistent state.
			res.err = ErrSequenceInconsistentState
			return
		}

		res.data = make([][]byte, length)

		copy(res.data, t.messages[offset:])

		t.subscribers[name] = t.sequence
	}

	res := <-result
	return res.data, res.err
}

// Close closes topic instance.
func (t *Topic) Close() {
	if t.isClosed() {
		return
	}
	close(t.done)
}

func (t *Topic) isClosed() bool {
	select {
	case _, ok := <-t.done:
		return !ok
	default:
		return false
	}
}

func (t *Topic) clean() {
	min := t.sequence
	for _, seq := range t.subscribers {
		if seq < min {
			min = seq
		}
	}

	offset := len(t.messages) - (t.sequence - min)

	t.messages = shift(t.messages, offset)
}

func shift(data [][]byte, offset int) [][]byte {
	if offset <= 0 {
		return data
	}

	if offset > len(data) {
		offset = len(data)
	}

	copy(data[:], data[offset:])
	for k, n := len(data)-offset, len(data); k < n; k++ {
		data[k] = nil
	}
	return data[:len(data)-offset]
}
