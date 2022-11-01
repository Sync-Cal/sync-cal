package main

import (
	"encoding/json"
	"fmt"
	"github.com/apognu/gocal"
	"golang-calendar-sync/config"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getProgramConfig() *config.Config {
	buf, _ := ioutil.ReadFile("../config.json")

	c := &config.Config{}
	json.Unmarshal(buf, c)
	// if err != nil {
	//     // return nil, fmt.Errorf("in file %q: %w", filename, err)
	// }

	return c
}

func getIcalTextFromUrl(url string) (string, error) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("Error %f", err)
	}

	return string(body), nil
}

func getParsedCalFromUrl(url string) (*gocal.Gocal, error) {
	icalText, err := getIcalTextFromUrl(url)
	if err == nil {
		cal := gocal.NewParser(strings.NewReader(icalText))
		err = cal.Parse()
		if err == nil {
			return cal, nil
		}

	}
	return nil, err
}

func convertDate(date *time.Time) *calendar.EventDateTime {
	return &calendar.EventDateTime{DateTime: date.Format(time.RFC3339)}

}

func gocalToGoogle(event *gocal.Event) *calendar.Event {

	fmt.Println()
	return &calendar.Event{
		Summary:     event.Summary,
		Description: event.Description,
		Start:       convertDate(event.Start),
		End:         convertDate(event.End),
		ICalUID:     event.Uid,
		Location:    event.Location,
		// Created:          event.Created.String(),
		RecurringEventId: event.RecurrenceID,

		// Start: &calendar.EventDateTime{
		// 	DateTime: "2022-03-10T09:00:00-07:00",
		// 	TimeZone: "America/Los_Angeles",
		// },
		// End: &calendar.EventDateTime{
		// 	DateTime: "2022-03-10T17:00:00-07:00",
		// 	TimeZone: "America/Los_Angeles",
		// },

		// Recurrence: event.RecurrenceRule,
	}
}
