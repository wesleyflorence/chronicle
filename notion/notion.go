// Package notion for interacting with notion api
package notion

import (
	"strconv"
	"time"

	"github.com/jomei/notionapi"
)

func buildPageCreateRequest(parentID string, properties *notionapi.Properties) notionapi.PageCreateRequest {
	return notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(parentID),
		},
		Properties: *properties,
	}
}

// Return the week relative to August 21st
func weekRelativeToChemoStart(date time.Time) string {
	startDate := time.Date(date.Year(), time.August, 21, 0, 0, 0, 0, time.UTC)
	daysDiff := date.Sub(startDate).Hours() / 24
	weeksDiff := int(daysDiff)/7 + 1
	return strconv.Itoa(weeksDiff)
}
