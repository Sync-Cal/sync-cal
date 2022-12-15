package utils

import (
	"fmt"
	"github.com/apognu/gocal"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetAllGCalEvents(service *calendar.Service, calendarId string) ([]calendar.Event, error) {
	events := []calendar.Event{}

	response, err := service.Events.List(calendarId).Do()

	for true {
		if err != nil {
			return nil, err
		}

		for _, ev := range response.Items {
			events = append(events, *ev)
		}

		nextPageToken := response.NextPageToken
		if nextPageToken == "" {
			break
		}

		response, err = service.Events.List(calendarId).PageToken(nextPageToken).Do()
	}

	return events, nil
}

func GetIcalTextFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("Error %f", err)
	}

	return string(body), nil
}

func GetParsedCalFromUrl(url string) (*gocal.Gocal, error) {
	icalText, err := GetIcalTextFromUrl(url)
	if err == nil {
		cal := gocal.NewParser(strings.NewReader(icalText))
		err = cal.Parse()
		if err == nil {
			return cal, nil
		}

	}
	return nil, err
}

func ConvertDate(date *time.Time) *calendar.EventDateTime {
	return &calendar.EventDateTime{DateTime: date.Format(time.RFC3339)}

}

func GocalToGoogle(event *gocal.Event) *calendar.Event {
	return &calendar.Event{
		Summary:     event.Summary,
		Description: event.Description,
		Start:       ConvertDate(event.Start),
		End:         ConvertDate(event.End),
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

func GetAllGcalEventsFromUrl(url string) ([]calendar.Event, error) {
	gocal, err := GetParsedCalFromUrl(url)
	if err != nil {
		return nil, err
	}

	var events []calendar.Event

	for _, ev := range gocal.Events {
		events = append(events, *GocalToGoogle(&ev))
	}

	return events, nil
}

func TerraformOutput(resourceType string, resourceName string) string {
	str := `output "%s_%s" {
		value = %s.%s
	}`

	return fmt.Sprintf(str, resourceType, resourceName, resourceType, resourceName)
}

type JustType struct {
	Type string `json:"type"`
}
