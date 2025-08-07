package types

type ObjectEvent struct {
	Event
	Key string
}

func (e ObjectEvent) GetBase() Event {
	return e.Event
}
