package calendarsource

import (
	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type Managed struct {
	CalendarSource
	Uuid string `json:"uuid" validate:"required"`
}

func (s Managed) ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) AtomicCalendarSource {
	cal, _ := sourceStore.GetManaged(s.Uuid)
	service := (*services)[cal.ProviderEmail]
	return atomicManaged{ExternalId: cal.CalendarId, Service: service}
}

// ===================================================
//						ATOMIC
// ===================================================

type atomicManaged struct {
	AtomicCalendarSource
	ExternalId string
	Service    *calendar.Service
}
