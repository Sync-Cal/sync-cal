package setuptask

import "fmt"

type CalendarProvider struct {
	SetupTask
	Alias string
}

func (_ CalendarProvider) GetType() string { return "provider" }

func (provider CalendarProvider) GetTerraform() string {
	str := `
	provider "gcal" {
		alias        = "%s"
		client_id    = var.provider_client_id
		access_token = var.provider_access_token_%s
	  }

	variable "provider_access_token_%s" {
		type = string
	}
	`
	return fmt.Sprintf(str, provider.Alias, provider.Alias, provider.Alias)
}
