package setuptask

type ManagedSubscribe struct {
	SetupTask
	Uuid        string                     `json:"uuid" binding:"required"`
	Source      managedSubsribeSource      `json:"source" binding:"required"`
	Destination managedSubsribeDestination `json:"destination" binding:"required"`
}

type managedSubsribeSource struct {
	User string `json:"user" binding:"required"`
	Uuid string `json:"uuid" binding:"required"`
}
type managedSubsribeDestination struct {
	User string `json:"user" binding:"required"`
}

func (t ManagedSubscribe) ToAtomic() []AtomicSetupTask {

	return []AtomicSetupTask{
		AtomicAclRule{
			Provider:     t.Source.User,
			ResourceName: t.Uuid,
			CalendarId:   t.Source.Uuid,
			User:         "toon@techwolf.ai",
			Role:         "writer",
		},
		AtomicCalendarListEntry{
			Provider:     t.Destination.User,
			CalendarId:   t.Source.Uuid,
			ResourceName: t.Uuid,
		},
	}
}
