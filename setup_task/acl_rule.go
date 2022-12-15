package setuptask

import (
	"fmt"

	"github.com/Sync-Cal/sync-cal/utils"
)

// ===================================================
//						PUBLIC
// ===================================================

type AclRule struct {
	SetupTask
}

// TODO to atomic

// ===================================================
//						ATOMIC
// ===================================================

type AtomicAclRule struct {
	SetupTask
	Provider     string
	ResourceName string

	CalendarId string
	User       string
	Role       string
}

func (_ AtomicAclRule) GetType() string { return "acl_rule" }

func (aclRule AtomicAclRule) GetTerraform() string {
	str := `
	resource "gcal_aclrule" "%s" {
		provider = gcal.%s
		calendar_id = gcal_calendar.%s.id
		user = "%s"
		role = "%s"
	}
	`
	str = fmt.Sprintf(str, aclRule.ResourceName, aclRule.Provider, aclRule.CalendarId, aclRule.User, aclRule.Role)
	str += utils.TerraformOutput("gcal_aclrule", aclRule.ResourceName)
	return str
}
