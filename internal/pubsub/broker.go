package pubsub

import "sync"

type Subscribers map[string]Subscriber

type Broker interface {
	Subscribe(s Subscriber, topic string)
	Unsubscribe(s Subscriber, topic string)
	Publish(topic string, msg interface{})
	RemoveSubscriber(s Subscriber)
}

type broker struct {
	subscribers Subscribers
	topics      map[string]Subscribers
	mtx         sync.RWMutex
}

func NewBroker() *broker {
	return &broker{
		subscribers: Subscribers{},
		topics:      map[string]Subscribers{},
	}
}

func (b *broker) Subscribe(s Subscriber, topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}

	b.topics[topic][s.Id()] = s
}

func (b *broker) Unsubscribe(s Subscriber, topic string) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	delete(b.topics[topic], s.Id())

	s.RemoveTopic(topic)
}

func (b *broker) Publish(topic string, msg interface{}) {
	b.mtx.RLock()
	topics := b.topics[topic]
	b.mtx.RUnlock()

	for _, s := range topics {
		message := NewMessage(msg, topic)
		if !s.Active() {
			return
		}

		go func(s Subscriber) {
			s.Signal(message)
		}(s)
	}
}

func (b *broker) RemoveSubscriber(s Subscriber) {
	for topic := range s.Topics() {
		b.Unsubscribe(s, topic)
	}

	b.mtx.Lock()
	delete(b.subscribers, s.Id())
	b.mtx.Unlock()

	s.Destruct()
}
