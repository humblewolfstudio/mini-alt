package types

type GlobalEvent struct {
	Event
	Bucket string
}

func (e GlobalEvent) GetBase() Event {
	return e.Event
}
