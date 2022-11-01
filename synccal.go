package synccal

import (
	// "golang-calendar-sync/config"

	"fmt"
	"github.com/Sync-Cal/sync-cal/config"
	"google.golang.org/api/calendar/v3"
)

func executeCopyTask(
	task map[string]interface{},
	services map[string]*calendar.Service,
	tasksCache *map[string]interface{},
) error {
	t, err := config.ParseCopyTask(task)
	if err != nil {
		return err
	}

	if t.Source.Type == "external" && t.Destination.Type == "internal" {
		
	}

	return nil
}

func ExecuteTasks(tasks []map[string]interface{}, state map[string]interface{}, services map[string]*calendar.Service) error {
	tasksCache := make(map[string]interface{})
	for _, task := range tasks {
		switch task["type"] {
		case "copy":
			executeCopyTask(task, services, &tasksCache)
		default:
			return fmt.Errorf("Task type %s is not supported", task["type"])
		}
	}

	return nil
}

func ExecuteJob() {}

func ExecuteSetup(tasks []map[string]interface{}, services map[string]*calendar.Service) {

}
