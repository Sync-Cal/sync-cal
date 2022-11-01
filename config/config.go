package config

type Config struct {
	Users           map[string]string         `json:"users"`
	ExternalSources map[string]ExternalSource `json:"externalSources"`
	Setup           []map[string]interface{}  `json:"setup"`
	Jobs            []Job                     `json:"jobs"`
}

type Job struct {
	Cron  string `json:"cron"`
	Tasks []map[string]interface{} `json:"tasks"`
}

type CopyTask struct {
	Source      TaskSource
	Destination TaskSource
}

type SubscriptionTask struct {
	Source      TaskSource
	Destination TaskSource
}

type TaskSource struct {
	Type string
	Name string
}

type ExternalSource struct {
	Url string `json:"url"`
}
