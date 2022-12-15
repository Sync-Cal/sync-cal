package setuptask

import (
	"fmt"

	"github.com/Sync-Cal/sync-cal/utils"
)

type AtomicCalendarListEntry struct {
	AtomicSetupTask
	Provider     string
	ResourceName string

	CalendarId string
}

func (_ AtomicCalendarListEntry) GetType() string { return "list_entry" }


func (calendarListEntry AtomicCalendarListEntry) GetTerraform() string {
	str := `
	resource "gcal_calendarlistentry" "%s" {
		provider = gcal.%s
		calendar_id = gcal_calendar.%s.id
	}
	`
	str = fmt.Sprintf(str, calendarListEntry.ResourceName, calendarListEntry.Provider, calendarListEntry.CalendarId)
	str += utils.TerraformOutput("gcal_calendarlistentry", calendarListEntry.ResourceName)
	return str
}
