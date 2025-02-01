package async

// EventBus is a simple FIFO queue wrapper
// to pass messages between multiple components
type EventBus struct {
	in  chan Event
	out chan<- Event
}

func NewEventBus(in chan Event, out chan<- Event) *EventBus {
	return &EventBus{
		in:  in,
		out: out,
	}
}

// Start triggers a new go routine to run the event bus.
// It returns immediatly after the go routine spawned.
func (b *EventBus) Start() {
	go b.run()
}

func (b *EventBus) run() {
	for {
		evt := <-b.in
		b.out <- evt
	}
}

func (b *EventBus) Publish(evt Event) {
	b.in <- evt
}
