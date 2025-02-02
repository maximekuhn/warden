package queue

import (
	"errors"
	"log/slog"

	"github.com/maximekuhn/warden/internal/domain/async"
)

var ErrQueueFull = errors.New("queue is full")

// EventsQueue implements a simple EventBus, where each event type has its own
// FIFO queue. It handles dispatching events to correct listeners, without
// blocking the callers of Publish* functions.
//
// It rejects an event publication if the corresponding queue has reached
// maximum capacity
type EventsQueue struct {
	// FIFO queue for StartServerEvent
	start chan async.StartServerEvent

	// FIFO queue for ServerStartedEvent
	started chan async.ServerStartedEvent

	logger *slog.Logger
}

func NewEventsQeue(queueSize uint, l *slog.Logger) *EventsQueue {
	start := make(chan async.StartServerEvent, queueSize)
	started := make(chan async.ServerStartedEvent, queueSize)
	return &EventsQueue{
		start:   start,
		started: started,
		logger:  l,
	}
}

// StartListeners setup the provided listeners and returns immediatly
func (q *EventsQueue) StartListeners(
	start *async.StartServerEventListener,
	started *async.ServerStartedEventListener,
) {
	go func() {
		select {
		case evt := <-q.start:
			go start.Execute(evt)
		case evt := <-q.started:
			go started.Execute(evt)
		}
	}()
}

func (q *EventsQueue) PublishStartServerEvent(evt async.StartServerEvent) error {
	select {
	case q.start <- evt:
		q.logger.Info("sent StartServerEvent to listener")
		return nil
	default:
		q.logger.Info("queue for StartServerEvent is full")
		return ErrQueueFull
	}
}

func (q *EventsQueue) PublishServerStartedEvent(evt async.ServerStartedEvent) error {
	select {
	case q.started <- evt:
		q.logger.Info("sent ServerStartedEvent to listener")
		return nil
	default:
		q.logger.Info("queue for ServerStartedEvent is full")
		return ErrQueueFull
	}
}
