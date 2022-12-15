package task

import (
	"encoding/json"
	"errors"

	"google.golang.org/api/calendar/v3"

	caldest "github.com/Sync-Cal/sync-cal/calendar_destination"
	calsource "github.com/Sync-Cal/sync-cal/calendar_source"
	"github.com/Sync-Cal/sync-cal/common"
	copyrule "github.com/Sync-Cal/sync-cal/copy_rule"
	"github.com/Sync-Cal/sync-cal/utils"
)

// ===================================================
//						PUBLIC
// ===================================================

type CopyTask struct {
	TaskBase
	Source      calsource.CalendarSource    `json:"-"`
	Destination caldest.CalendarDestination `json:"-"`
	Rules       []copyrule.CopyRule         `json:"-"`

	RawSource      json.RawMessage   `json:"source"`
	RawDestination json.RawMessage   `json:"destination"`
	RawRules       []json.RawMessage `json:"rules"`
}

func (t *CopyTask) UnmarshalJSON(b []byte) error {
	type copytask CopyTask
	err := json.Unmarshal(b, (*copytask)(t))
	if err != nil {
		return err
	}

	// Destination
	var rD utils.JustType
	json.Unmarshal(t.RawDestination, &rD)
	var destination caldest.CalendarDestination

	switch rD.Type {
	case "managed":
		destination = &caldest.Managed{}
	default:
		return errors.New("unknown Destination type")
	}

	json.Unmarshal(t.RawDestination, destination)
	t.Destination = destination

	// Source
	var rS utils.JustType
	json.Unmarshal(t.RawSource, &rS)
	var source calsource.CalendarSource

	switch rS.Type {
	case "managed":
		source = &calsource.Managed{}
	case "url":
		source = &calsource.Url{}
	default:
		return errors.New("unknown Source type")
	}

	json.Unmarshal(t.RawSource, source)
	t.Source = source

	// Rules
	for _, raw := range t.RawRules {
		var tt utils.JustType
		err = json.Unmarshal(raw, &tt)

		if err != nil {
			return err
		}

		var i interface{}
		switch tt.Type {
		case "add_attendees":
			i = &copyrule.AddAttendees{}
		case "filter_duration":
			i = &copyrule.FilterDuration{}
		default:
			return errors.New("unknown rule type")
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		t.Rules = append(t.Rules, i.(copyrule.CopyRule))
	}

	return nil
}

func (t CopyTask) ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) *AtomicTask {

	atomicRules := []copyrule.AtomicCopyRule{}
	for _, r := range t.Rules {
		atomicRules = append(atomicRules, r.ToAtomic(sourceStore))
	}

	td := AtomicCopyData{
		Source:      t.Source.ToAtomic(services, sourceStore),
		Destination: t.Destination.ToAtomic(services, sourceStore),
		Rules:       atomicRules,
	}
	return &AtomicTask{Id: t.Uuid, TaskData: td}
}

// ===================================================
//						ATOMIC
// ===================================================

type AtomicCopyData struct {
	AtomicTaskData
	Rules       []copyrule.AtomicCopyRule
	Source      calsource.AtomicCalendarSource
	Destination caldest.AtomicCalendarDestination
}

func (task AtomicCopyData) Execute() error {
	events, err := task.Source.GetEvents()
	if err != nil {
		return err
	}

	// Apply rules
	newEvents := []calendar.Event{}
	for _, ev := range events {
		hasBeenFiltered := false

		for _, rule := range task.Rules {
			keep, err := rule.Apply(&ev)
			if err != nil {
				return err
			}
			if !keep {
				hasBeenFiltered = true
				break
			}
		}
		if !hasBeenFiltered {
			newEvents = append(newEvents, ev)
		}
	}

	err = task.Destination.UpsertEvents(newEvents)
	return err
}
