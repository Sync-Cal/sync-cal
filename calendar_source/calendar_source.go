package calendarsource

import (
	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)


type CalendarSource interface {
	ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) AtomicCalendarSource
}

type AtomicCalendarSource interface {
	GetEvents() ([]calendar.Event, error)
}
