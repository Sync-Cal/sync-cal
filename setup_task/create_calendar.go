package setuptask

import (
	"fmt"
	"github.com/Sync-Cal/sync-cal/utils"
)

// ===================================================
//						PUBLIC
// ===================================================

type CreateCalendar struct {
	SetupTask
	// setupTask
	User    string `json:"user" binding:"required"`
	Uuid    string `json:"uuid" binding:"required"`
	Details CalendarDetails
}

type CalendarDetails struct {
	Summary string `json:"summary"`
}

func (t CreateCalendar) ToAtomic() []AtomicSetupTask {
	atomicTask := AtomicCreateCalendar{
		Provider:     t.User,
		ResourceName: t.Uuid,
		Summary:      t.Details.Summary,
	}
	return []AtomicSetupTask{atomicTask}
}

// ===================================================
//						ATOMIC
// ===================================================


type AtomicCreateCalendar struct {
	SetupTask
	Provider     string
	ResourceName string

	Summary string
}

func (_ AtomicCreateCalendar) GetType() string { return "create" }

func (calendar AtomicCreateCalendar) GetTerraform() string {
	str := `
	resource "gcal_calendar" "%s" {
		provider = gcal.%s
		summary = "%s"
	}
	`

	str = fmt.Sprintf(str, calendar.ResourceName, calendar.Provider, calendar.Summary)
	str += utils.TerraformOutput("gcal_calendar", calendar.ResourceName)

	return str
}
