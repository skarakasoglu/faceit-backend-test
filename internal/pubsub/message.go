package pubsub

import "time"

type Message struct {
	topic     string
	body      interface{}
	createdAt time.Time
}

func NewMessage(msg interface{}, topic string) *Message {
	return &Message{
		topic:     topic,
		body:      msg,
		createdAt: time.Now(),
	}
}

func (m *Message) Topic() string {
	return m.topic
}

func (m *Message) Body() interface{} {
	return m.body
}

func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}
