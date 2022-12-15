package calendardestination

import (
	"fmt"
	// "sync"

	"github.com/Sync-Cal/sync-cal/common"
	"github.com/Sync-Cal/sync-cal/utils"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type Managed struct {
	CalendarDestination
	Uuid string `json:"uuid" validate:"required"`
}

func (s Managed) ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) AtomicCalendarDestination {
	cal, _ := sourceStore.GetManaged(s.Uuid)
	service := (*services)[cal.ProviderEmail]
	return atomicManaged{ExternalId: cal.CalendarId, Service: service}
}

// ===================================================
//						ATOMIC
// ===================================================

type atomicManaged struct {
	AtomicCalendarDestination
	ExternalId string
	Service    *calendar.Service
}

func (dest atomicManaged) UpsertEvents(events []calendar.Event) error {
	currentEvents, _ := utils.GetAllGCalEvents(dest.Service, dest.ExternalId)
	currentEventsMap := map[string]calendar.Event{}
	for _, e := range currentEvents {
		currentEventsMap[e.Id] = e
	}

	insertEvent := func(event calendar.Event) {
		var err error
		for i := 1; i <= 3; i++ { // Retry 3 times
			_, err := dest.Service.Events.Insert(dest.ExternalId, &event).Do()
			if err == nil {
				break
			}
		}
		if err != nil {
			if err != nil {
				fmt.Printf("FAILED to insert: %s\n", event.Summary)
				fmt.Println(err)
			} else {
				fmt.Printf("successfully inserted: %s\n", event.Summary)
			}
			fmt.Println()
		}
	}

	updateEvent := func(event calendar.Event) {
		var err error
		for i := 1; i <= 3; i++ { // Retry 3 times
			_, err := dest.Service.Events.Insert(dest.ExternalId, &event).Do()
			if err == nil {
				break
			}
		}
		if err != nil {
			if err != nil {
				fmt.Printf("FAILED to update: %s\n", event.Summary)
				fmt.Println(err)
			} else {
				fmt.Printf("successfully updated: %s\n", event.Summary)
			}
			fmt.Println()
		}
	}

	for _, newEvent := range events {

		// Check if event is new/updated
		currentEvent, ok := currentEventsMap[newEvent.Id]

		action := "insert"

		if ok {
			if currentEvent.Updated == newEvent.Updated {
				action = ""
			} else {
				action = "update"
			}
		}

		// Execute event
		switch action {
		case "insert":
			insertEvent(newEvent)
		case "update":
			updateEvent(newEvent)
		}

	}
	return nil
}
