package pubsub

import "sync"

type Subscribers map[string]*Subscriber

type Broker struct {
	subscribers Subscribers
	topics      map[string]Subscribers
	mtx         sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: Subscribers{},
		topics:      map[string]Subscribers{},
	}
}

func (b *Broker) Subscribe(s *Subscriber, topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}

	b.topics[topic][s.id] = s
}

func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	delete(b.topics[topic], s.id)

	s.RemoveTopic(topic)
}

func (b *Broker) Publish(topic string, msg interface{}) {
	b.mtx.RLock()
	topics := b.topics[topic]
	b.mtx.RUnlock()

	for _, s := range topics {
		message := NewMessage(msg, topic)
		if !s.active {
			return
		}

		go func(s *Subscriber) {
			s.Signal(message)
		}(s)
	}
}

func (b *Broker) RemoveSubscriber(s *Subscriber) {
	for topic := range s.topics {
		b.Unsubscribe(s, topic)
	}

	b.mtx.Lock()
	delete(b.subscribers, s.id)
	b.mtx.Unlock()

	s.Destruct()
}
