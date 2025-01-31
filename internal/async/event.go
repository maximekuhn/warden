package async

type Event int

const (
	_ Event = iota
	EventStartServer
	EventStopServer
)
