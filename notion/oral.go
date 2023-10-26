// Package notion for interacting with notion api
package notion

import (
	"context"
	"time"

	"github.com/jomei/notionapi"
)

func AppendOralEntry(client *notionapi.Client, dbID string, sensation, activity, note string) (*notionapi.Page, error) {
	currentTime := notionapi.Date(time.Now())
	props := notionapi.Properties{
		"Week": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: weekRelativeToChemoStart(time.Now())}},
			},
		},
		"Sensation": notionapi.MultiSelectProperty{
			MultiSelect: []notionapi.Option{
				{
					Name:  sensation,
					Color: lookupSensationColor(sensation)},
			},
		},
		"Activity": notionapi.MultiSelectProperty{
			MultiSelect: []notionapi.Option{
				{
					Name:  activity,
					Color: lookupActivityColor(activity)},
			},
		},
		"Date": notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &currentTime,
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

func lookupSensationColor(sensation string) notionapi.Color {
	switch sensation {
	case "neutral":
		return notionapi.ColorBlue
	case "yuck":
		return notionapi.ColorYellow
	case "acidic":
		return notionapi.ColorRed
	default:
		return notionapi.ColorGray
	}
}

func lookupActivityColor(activity string) notionapi.Color {
	switch activity {
	case "Brush Rembrandt":
		return notionapi.ColorPink
	case "Brush Arm and Hammer":
		return notionapi.ColorOrange
	case "Floss":
		return notionapi.ColorGreen
	case "Mouth Wash":
		return notionapi.ColorBlue
	default:
		return notionapi.ColorGray
	}
}
