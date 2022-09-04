package pubsub

type MockBroker struct {
	PublishMock          func(string, interface{})
	SubscribeMock        func(Subscriber, string)
	UnsubscribeMock      func(Subscriber, string)
	RemoveSubscriberMock func(Subscriber)
}

func (m *MockBroker) Subscribe(s Subscriber, topic string) {
	m.SubscribeMock(s, topic)
}

func (m *MockBroker) Unsubscribe(s Subscriber, topic string) {
	m.UnsubscribeMock(s, topic)
}

func (m *MockBroker) Publish(topic string, msg interface{}) {
	m.PublishMock(topic, msg)
}

func (m *MockBroker) RemoveSubscriber(s Subscriber) {
	m.RemoveSubscriberMock(s)
}
