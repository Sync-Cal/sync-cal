package copyrule

import (
	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type AddAttendees struct {
	CopyRule
	Attendees    []string `json:"attendees"`
	ManagedUuids []string `json:"managed_uuids"`
}

func (r AddAttendees) ToAtomic(sourceStore common.CalendarSourceStore) AtomicCopyRule {
	var ids []string

	for _, uuid := range r.ManagedUuids {
		cal, _ := sourceStore.GetManaged(uuid)
		ids = append(ids, cal.CalendarId)
	}

	return AtomicAddAttendees{Attendees: r.Attendees, AttendingIds: ids}
}

// ===================================================
//						ATOMIC
// ===================================================

type AtomicAddAttendees struct {
	AtomicCopyRule
	Attendees    []string
	AttendingIds []string
}

func (f AtomicAddAttendees) Apply(event *calendar.Event) (bool, error) {
	for _, email := range f.Attendees {
		att := calendar.EventAttendee{Email: email}
		event.Attendees = append(event.Attendees, &att)
	}
	for _, id := range f.AttendingIds {
		att := calendar.EventAttendee{Email: id}
		event.Attendees = append(event.Attendees, &att)
	}
	return true, nil
}
