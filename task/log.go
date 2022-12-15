package task

import (
	"fmt"
	"time"

	"github.com/Sync-Cal/sync-cal/common"
	"google.golang.org/api/calendar/v3"
)

// ===================================================
//						PUBLIC
// ===================================================

type LogTask struct {
	TaskBase
	Message string `json:"message"`
}

func (t LogTask) ToAtomic(services *map[string]*calendar.Service, sourceStore common.CalendarSourceStore) *AtomicTask {
	return &AtomicTask{Id: t.Uuid, TaskData: AtomicLogData{Message: t.Message}}
}

// ===================================================
//						ATOMIC
// ===================================================

type AtomicLogData struct {
	AtomicTaskData
	Message string
	Time    time.Duration
}

func (task AtomicLogData) Execute() error {
	fmt.Println("STARTING " + task.Message)
	time.Sleep(time.Second)
	fmt.Println("DONE " + task.Message)
	return nil
}
