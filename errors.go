package sps

// Package errors.
const (
	ErrNoSubscriptionFound       = Error("no subscription found")
	ErrNoTopicFound              = Error("no topic found")
	ErrSequenceInconsistentState = Error("sequence system has an inconsistent state")
)

// Error represents constant error.
type Error string

func (e Error) Error() string {
	return string(e)
}
