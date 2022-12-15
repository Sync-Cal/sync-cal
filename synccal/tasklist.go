package synccal

import (
	"encoding/json"
	"errors"

	"github.com/Sync-Cal/sync-cal/common"
	Task "github.com/Sync-Cal/sync-cal/task"
	"github.com/Sync-Cal/sync-cal/utils"
	"google.golang.org/api/calendar/v3"
)

type TaskList struct {
	Tasks    []Task.Task       `json:"-"`
	RawTasks []json.RawMessage `json:"tasks"`
}

func (tl *TaskList) UnmarshalJSON(b []byte) error {
	type tasklist TaskList
	err := json.Unmarshal(b, (*tasklist)(tl))
	if err != nil {
		return err
	}

	for _, raw := range tl.RawTasks {
		var t utils.JustType
		err = json.Unmarshal(raw, &t)
		if err != nil {
			return err
		}

		var i interface{}
		switch t.Type {
		case "copy":
			i = &Task.CopyTask{}
		default:
			return errors.New("unknown Task type")
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		tl.Tasks = append(tl.Tasks, i.(Task.Task))
	}
	return nil
}

// Parameters:
//   - `services` : {"personal@gmail.com": Service}
//   - `managedCalendars` : {"personal": ManagedCalendar}
func (taskList *TaskList) Execute(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) {
	tasksMap := map[string]*Task.AtomicTask{}

	for _, task := range taskList.Tasks {
		tasksMap[task.GetUuid()] = task.ToAtomic(services, sourceStore)
	}

	for _, task := range taskList.Tasks {
		toTask := tasksMap[task.GetUuid()]
		for _, dependency := range task.GetDependencies() {
			fromTask := tasksMap[dependency.Id]
			link := Task.AtomicTaskLink{FromTask: fromTask, ToTask: toTask}
			toTask.BackLinks = append(toTask.BackLinks, &link)
			fromTask.ForwardLinks = append(fromTask.ForwardLinks, &link)
			toTask.RemainingBlocks += 1
		}
	}

	rootTasks := []*Task.AtomicTask{}
	for _, task := range tasksMap {
		if len(task.BackLinks) == 0 {
			rootTasks = append(rootTasks, task)
		}
	}

	runAtomicTasks(rootTasks)
}

func runAtomicTasks(tasks []*Task.AtomicTask) {

	for _, task := range tasks {
		go task.Execute()
	}
}
