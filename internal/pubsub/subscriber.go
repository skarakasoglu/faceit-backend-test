package pubsub

import (
	"github.com/google/uuid"
	"sync"
)

const messageBuffer = 100

type Subscriber interface {
	AddTopic(string)
	RemoveTopic(string)
	Signal(msg *Message)
	Destruct()
	Messages() <-chan *Message
	Topics() map[string]bool
	Id() string
	Active() bool
}

type subscriber struct {
	id       string
	messages chan *Message
	topics   map[string]bool
	active   bool
	mtx      sync.RWMutex
}

func NewSubscriber() *subscriber {
	id := uuid.New().String()
	return &subscriber{
		id:       id,
		messages: make(chan *Message, messageBuffer),
		topics:   map[string]bool{},
		active:   true,
	}
}

func (s *subscriber) AddTopic(topic string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.topics[topic] = true
}

func (s *subscriber) RemoveTopic(topic string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.topics[topic] = false
}

func (s *subscriber) Signal(msg *Message) {
	if s.active {
		s.messages <- msg
	}
}

func (s *subscriber) Destruct() {
	s.active = false
	close(s.messages)
}

func (s *subscriber) Messages() <-chan *Message {
	return s.messages
}

func (s *subscriber) Topics() map[string]bool {
	return s.topics
}

func (s *subscriber) Id() string {
	return s.id
}

func (s *subscriber) Active() bool {
	return s.active
}
