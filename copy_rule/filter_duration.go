package copyrule

import (
	"time"

	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type FilterDuration struct {
	CopyRule
	Max int
	Min int
}

func (r FilterDuration) ToAtomic(sourceStore common.CalendarSourceStore) AtomicCopyRule {
	return AtomicFilterDuration{Max: r.Max, Min: r.Min}
}

// ===================================================
//						ATOMIC
// ===================================================

type AtomicFilterDuration struct {
	AtomicCopyRule
	Max int
	Min int
}

func (f AtomicFilterDuration) Apply(event *calendar.Event) (bool, error) {
	RFC3339local := "2006-01-02T15:04:05Z"
	loc, err := time.LoadLocation(event.Start.TimeZone)

	if err != nil {
		return true, err
	}

	
	
	start, _ := time.ParseInLocation(RFC3339local, event.Start.DateTime, loc)
	end, _ := time.ParseInLocation(RFC3339local, event.End.DateTime, loc)
	
	durationMinutes := end.Sub(start).Minutes()

	if durationMinutes <= float64(f.Max) && durationMinutes >= float64(f.Min) {
		return true, nil
	}
	return false, nil

}
