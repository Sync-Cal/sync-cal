package synccal

import (
	"bytes"

	STask "github.com/Sync-Cal/sync-cal/setup_task"
)

// ====================

func GenerateTerraform(setupTasks []STask.AtomicSetupTask) string {
	var buffer bytes.Buffer

	for _, task := range setupTasks {
		buffer.WriteString(task.GetTerraform())
	}

	return buffer.String()
}

// ================================================================
//								SETUPSOURCES
// ================================================================

type SetupSource interface {
	ToAtomic() []STask.SetupTask
	GetUuid() string
}

type SetupUrlSource struct {
	SetupSource
	Uuid string `json:"uuid" binding:"required"`
	Url  string `json:"url" binding:"required"`
}

func (s SetupUrlSource) GetUuid() string { return s.Uuid }
