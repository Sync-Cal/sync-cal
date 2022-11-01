package config

import (
	// "encoding/json"
	"encoding/json"
	// "fmt"
	"io/ioutil"
)

func ParseConfig(filename string) *Config {
	buf, _ := ioutil.ReadFile(filename)
	c := &Config{}
	json.Unmarshal(buf, c)
	return c
}

func ParseCopyTask(task map[string]interface{}) (*CopyTask, error) {
	jsonString, _ := json.Marshal(task)

	t := &CopyTask{}
	err := json.Unmarshal(jsonString, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
