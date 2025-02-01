package async

// EventBus allows publishing events, such as starting or stopping a minecraft
// server.
//
// Where and how events are published is the implementation's responsibility.
//
// Publish* functions must be non-blocking for the caller.
type EventBus interface {
	PublishStartServerEvent(evt StartServerEvent)

	PublishServerStartedEvent(evt ServerStartedEvent)
}
