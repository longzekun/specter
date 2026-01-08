package core

const (
	eventBufSize = 5
)

type Event struct {
	EventType uint32
}

type eventBroker struct {
	stop        chan struct{}
	publish     chan Event
	subscribe   chan chan Event
	unsubscribe chan chan Event
}

func (broker *eventBroker) Start() {
	subscribes := map[chan Event]struct{}{}
	for {
		select {
		case <-broker.stop:
			for sub := range subscribes {
				close(sub)
			}
			return
		case sub := <-broker.subscribe:
			subscribes[sub] = struct{}{}
		case unsub := <-broker.unsubscribe:
			delete(subscribes, unsub)
		case event := <-broker.publish:
			for sub := range subscribes {
				sub <- event
			}
		}
	}
}

func (broker *eventBroker) Stop() {
	close(broker.stop)
}

func (broker *eventBroker) Subscribe() chan Event {
	events := make(chan Event, eventBufSize)
	broker.subscribe <- events
	return events
}

func (broker *eventBroker) Unsubscribe(events chan Event) {
	broker.unsubscribe <- events
	close(events)
}

func (broker *eventBroker) Publish(event Event) {
	broker.publish <- event
}

func newEventBroker() *eventBroker {
	broker := &eventBroker{
		stop:        make(chan struct{}),
		publish:     make(chan Event, eventBufSize),
		subscribe:   make(chan chan Event, eventBufSize),
		unsubscribe: make(chan chan Event, eventBufSize),
	}

	go broker.Start()

	return broker
}

var (
	EventBroker = newEventBroker()
)
