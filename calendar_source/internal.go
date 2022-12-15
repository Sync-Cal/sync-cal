package calendarsource

import "google.golang.org/api/calendar/v3"

type Internal struct {
	Uuid   string
	events []*calendar.Event
}

func (source Internal) GetEvents() ([]*calendar.Event, error) {
	return source.events, nil // TODO
}