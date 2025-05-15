package events

type Event struct {
	Topic   string
	Key     string
	Payload []byte
}
