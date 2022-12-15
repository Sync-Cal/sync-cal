package common

import "fmt"

type CalendarSourceStore struct {
	Managed   map[string]ManagedCalendar
	Unmanaged map[string]UnmanagedCalendar
	Url       map[string]UrlCalendar
}

func (s CalendarSourceStore) GetManaged(uuid string) (*ManagedCalendar, error) {
	c, ok := s.Managed[uuid]
	if ok {
		return &c, nil
	}
	return nil, fmt.Errorf("Managed calendar " + uuid + " does not exist!")
}

func (s CalendarSourceStore) GetUrl(uuid string) (*UrlCalendar, error) {
	c, ok := s.Url[uuid]
	if ok {
		return &c, nil
	}
	return nil, fmt.Errorf("Url calendar " + uuid + " does not exist!")
}
