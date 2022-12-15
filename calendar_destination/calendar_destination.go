package calendardestination

import (
	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

type CalendarDestination interface {
	ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) AtomicCalendarDestination
}

type AtomicCalendarDestination interface {
	UpsertEvents([]calendar.Event) error
}
