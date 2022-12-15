package calendarsource

import (
	"github.com/Sync-Cal/sync-cal/common"
	"github.com/Sync-Cal/sync-cal/utils"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type Url struct {
	CalendarSource
	Uuid string `json:"uuid" validate:"required"`
}

func (s Url) ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) AtomicCalendarSource {
	source, _ := sourceStore.GetUrl(s.Uuid)
	return atomicUrl{Url: source.Url}
}

// ===================================================
//						ATOMIC
// ===================================================

type atomicUrl struct {
	AtomicCalendarSource
	Url string
}

func (source atomicUrl) GetEvents() ([]calendar.Event, error) {
	return utils.GetAllGcalEventsFromUrl(source.Url)
}
