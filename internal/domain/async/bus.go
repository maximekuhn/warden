package async

// EventBus allows publishing events, such as starting or stopping a minecraft
// server.
//
// Where and how events are published is the implementation's responsibility.
//
// Publish* functions must be non-blocking for the caller.
type EventBus interface {
	// A non-nil error indicates that the event could not be published.
	PublishStartServerEvent(evt StartServerEvent) error

	// A non-nil error indicates that the event could not be published.
	PublishServerStartedEvent(evt ServerStartedEvent) error

	// A non-nil error indicates that the event could not be published.
	PublishStopServerEvent(evt StopServerEvent) error
}
