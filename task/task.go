package task

import (
	"sync"

	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type TaskBase struct {
	Task
	Uuid         string       `json:"uuid" validate:"required"`
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Id string
}

type Task interface {
	ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) *AtomicTask
	GetUuid() string
	GetDependencies() []Dependency
}

func (task TaskBase) GetUuid() string {
	return task.Uuid
}

func (task TaskBase) GetDependencies() []Dependency {
	return task.Dependencies
}

// ===================================================
//						ATOMIC
// ===================================================

type AtomicTask struct {
	AtomicTaskInfo
	TaskData AtomicTaskData
	Id       string

	lock sync.Mutex
}

type AtomicTaskInfo struct {
	BackLinks       []*AtomicTaskLink
	ForwardLinks    []*AtomicTaskLink
	RemainingBlocks int // Number of backlinks that are still blocking this task

	Completed    bool // Whether the task is completed
	ErrorMessage bool // The error message occured during the task
}

type AtomicTaskData interface {
	Execute() error
}

type AtomicTaskLink struct {
	FromTask *AtomicTask
	ToTask   *AtomicTask
}

func (task *AtomicTask) markCompleted() {
	task.Completed = true

	for _, link := range task.ForwardLinks {
		// fmt.Printf("%s", link.toTask.TaskData)
		link.ToTask.notify()
	}
}

func (task *AtomicTask) notify() {
	task.lock.Lock()
	task.RemainingBlocks -= 1
	if task.RemainingBlocks == 0 {
		go task.Execute()
	}
	task.lock.Unlock()
}

func (task *AtomicTask) Execute() {
	task.TaskData.Execute()
	task.markCompleted()
}
