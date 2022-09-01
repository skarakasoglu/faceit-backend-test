package pubsub

import (
	"github.com/google/uuid"
	"sync"
)

const messageBuffer = 100

type Subscriber struct {
	id       string
	messages chan *Message
	topics   map[string]bool
	active   bool
	mtx      sync.RWMutex
}

func NewSubscriber() *Subscriber {
	id := uuid.New().String()
	return &Subscriber{
		id:       id,
		messages: make(chan *Message, messageBuffer),
		topics:   map[string]bool{},
		active:   true,
	}
}

func (s *Subscriber) AddTopic(topic string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.topics[topic] = false
}

func (s *Subscriber) Signal(msg *Message) {
	if s.active {
		s.messages <- msg
	}
}

func (s *Subscriber) Destruct() {
	s.active = false
	close(s.messages)
}

func (s *Subscriber) Messages() <-chan *Message {
	return s.messages
}
