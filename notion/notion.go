// Package notion for interacting with notion api
package notion

import (
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
