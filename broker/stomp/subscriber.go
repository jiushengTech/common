package stomp

import (
	"sync"

	stompV3 "github.com/go-stomp/stomp/v3"
	"github.com/jiushengTech/common/broker"
)

type subscriber struct {
	sync.RWMutex

	b *stompBroker

	options broker.SubscribeOptions
	topic   string
	sub     *stompV3.Subscription
	closed  bool
}

func (s *subscriber) Options() broker.SubscribeOptions {
	s.RLock()
	defer s.RUnlock()

	return s.options
}

func (s *subscriber) Topic() string {
	s.RLock()
	defer s.RUnlock()

	return s.topic
}

func (s *subscriber) Unsubscribe(removeFromManager bool) error {
	s.Lock()
	defer s.Unlock()

	s.closed = true

	var err error
	if s.sub != nil {
		err = s.sub.Unsubscribe()
	}

	if s.b != nil && s.b.subscribers != nil && removeFromManager {
		_ = s.b.subscribers.RemoveOnly(s.topic)
	}

	return err
}

func (s *subscriber) IsClosed() bool {
	s.RLock()
	defer s.RUnlock()

	return s.closed
}
