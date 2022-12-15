package copyrule

import (
	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

type CopyRule interface {
	ToAtomic(sourceStore common.CalendarSourceStore) AtomicCopyRule
}

type AtomicCopyRule interface {
	Apply(event *calendar.Event) (bool, error)
}
