// Package notion for interacting with notion api
package notion

import (
	"context"
	"time"

	"github.com/jomei/notionapi"
)

// AppendDigestionEntry adds a digestion entry to a Notion database.
//
// It creates a new page with the provided digestion entry details,
// including the Bristol scale, size, notes, and other properties.
func AppendDigestionEntry(client *notionapi.Client, dbID string, bristol int, size, note string) (*notionapi.Page, error) {
	currentTime := notionapi.Date(time.Now())
	props := notionapi.Properties{
		"Week": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: weekRelativeToChemoStart(time.Now())}},
			},
		},
		"Type": notionapi.MultiSelectProperty{
			MultiSelect: []notionapi.Option{
				{Name: "Poop", Color: notionapi.ColorBlue},
			},
		},
		"Date": notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &currentTime,
			},
		},
		"Bristol (1-7)": notionapi.NumberProperty{
			Number: float64(bristol),
		},
		"Size": notionapi.SelectProperty{
			Select: notionapi.Option{
				Name:  size,
				Color: lookupSizeColor(size),
			},
		},
		"Note": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: note}},
			},
		},
		"Chronicle": notionapi.CheckboxProperty{
			Checkbox: true,
		},
	}
	request := buildPageCreateRequest(dbID, &props)
	return client.Page.Create(context.Background(), &request)
}

func lookupSizeColor(size string) notionapi.Color {
	switch size {
	case "Tiny":
		return notionapi.ColorPink
	case "Small":
		return notionapi.ColorRed
	case "Medium":
		return notionapi.ColorOrange
	case "Large":
		return notionapi.ColorYellow
	default:
		return notionapi.ColorBlue
	}
}
