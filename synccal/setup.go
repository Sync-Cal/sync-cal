package synccal

import (
	"encoding/json"
	"errors"

	"github.com/Sync-Cal/sync-cal/common"
	STask "github.com/Sync-Cal/sync-cal/setup_task"
	"github.com/Sync-Cal/sync-cal/utils"
)

type Setup struct {
	Users      []string          `json:"users" binding:"required"`
	Tasks      []STask.SetupTask `json:"-"`
	Sources    []SetupSource     `json:"-"`
	RawSources []json.RawMessage `json:"sources" binding:"required"`
	RawTasks   []json.RawMessage `json:"tasks" binding:"required"`
}

func (setup Setup) GetEmptyTerraform() string {
	atomicTasks := []STask.AtomicSetupTask{}

	for _, u := range setup.Users {
		atomicTasks = append(atomicTasks, STask.CalendarProvider{Alias: u})
	}

	return GenerateTerraform(atomicTasks)
}

func (setup Setup) GetTerraform() string {
	atomicTasks := []STask.AtomicSetupTask{}

	for _, u := range setup.Users {
		atomicTasks = append(atomicTasks, STask.CalendarProvider{Alias: u})
	}

	for _, task := range setup.Tasks {
		atomicTasks = append(atomicTasks, task.ToAtomic()...)
	}

	return GenerateTerraform(atomicTasks)
}

// Parameters:
//   - `userMap` : {"personal": "personal@gmail.com", "work": "name@work.com"}
//   - `uuidMap` : {"personal_calendar": "6maqsgash83tedb793skafqg7s@group.calendar.google.com"}
func (setup Setup) GetManagedCalendars(userMap map[string]string, uuidMap map[string]string) []common.ManagedCalendar {
	atomicTasks := []STask.AtomicSetupTask{}

	for _, task := range setup.Tasks {
		atomicTasks = append(atomicTasks, task.ToAtomic()...)
	}

	calendars := []common.ManagedCalendar{}
	for _, aTask := range atomicTasks {
		if aTask.GetType() == "create" {
			c := aTask.(STask.AtomicCreateCalendar)
			mC := common.ManagedCalendar{
				ProviderEmail: userMap[c.Provider],
				CalendarId:    uuidMap[c.ResourceName],
				Uuid:          c.ResourceName,
			}
			calendars = append(calendars, mC)
		}
	}
	return calendars
}

func (s *Setup) UnmarshalJSON(b []byte) error {
	type setup Setup
	err := json.Unmarshal(b, (*setup)(s))
	if err != nil {
		return err
	}

	// TASKS
	for _, raw := range s.RawTasks {
		var t utils.JustType
		err = json.Unmarshal(raw, &t)
		if err != nil {
			return err
		}

		var i interface{}
		switch t.Type {
		case "create":
			i = &STask.CreateCalendar{}
		case "managed_subscribe":
			i = &STask.ManagedSubscribe{}
		default:
			return errors.New("unknown task type: " + t.Type)
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		s.Tasks = append(s.Tasks, i.(STask.SetupTask))
	}

	// SOURCES
	for _, raw := range s.RawSources {
		var t utils.JustType
		err = json.Unmarshal(raw, &t)
		if err != nil {
			return err
		}

		var i interface{}
		switch t.Type {
		case "url":
			i = &SetupUrlSource{}
		default:
			return errors.New("unknown source type: " + t.Type)
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		s.Sources = append(s.Sources, i.(SetupSource))
	}

	return nil
}
