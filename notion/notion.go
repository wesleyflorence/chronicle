// Package Notion for interacting with notion api
package notion

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jomei/notionapi"
)

func AppendMedicineEntry(client *notionapi.Client, medicinePageID, medicine, note string) (*notionapi.Page, error) {
	// Get Block form Medicine Page
	medsMap := make(map[string]string)
	children, err := client.Block.GetChildren(context.Background(), notionapi.BlockID(
		medicinePageID), &notionapi.Pagination{PageSize: 10})
	if err != nil {
		log.Fatal("Error Calling GetChildren")
		return nil, err
	}

	// Populate the Map using Medicine Child Page Titles
	for _, child := range children.Results {
		if block, ok := child.(*notionapi.ChildPageBlock); ok {
			fmt.Println(block.ChildPage.Title)
			medsMap[block.ChildPage.Title] = block.ID.String()
		}
	}
	if _, ok := medsMap[medicine]; !ok {
		return nil, fmt.Errorf("%s not found in the Medicine Page", medicine)
	}

	// Get the Blocks from the provided medicine's Page
	medPage, err := client.Block.GetChildren(context.Background(), notionapi.BlockID(medsMap[medicine]),
		&notionapi.Pagination{PageSize: 10})
	if err != nil {
		log.Fatal("Error Calling GetChildren2")
		return nil, err
	}

	// Get the Database ID
	var childDbID string
	for _, block := range medPage.Results {
		if block, ok := block.(*notionapi.ChildDatabaseBlock); ok {
			childDbID = block.ID.String()
		}
	}

	// Query DB to get the row count to inc the Dose
	query := notionapi.DatabaseQueryRequest{PageSize: 100}
	var cursor *notionapi.Cursor
	hasMore := true
	dose := 1

	for hasMore {
		if cursor != nil {
			query.StartCursor = *cursor
		}
		queryResult, err := client.Database.Query(context.Background(),
			notionapi.DatabaseID(childDbID),
			&query)

		if err != nil {
			log.Fatal("Error Making Query to get Count")
			return nil, err
		}
		dose = dose + len(queryResult.Results)
		cursor = &queryResult.NextCursor
		hasMore = queryResult.HasMore
	}

	// Append the Row
	currentTime := notionapi.Date(time.Now())
	props := notionapi.Properties{
		"Date": notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &currentTime,
			},
		},
		"Dose": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: strconv.Itoa(dose)}},
			},
		},
		"Status": notionapi.StatusProperty{
			Status: notionapi.Status{
				Name:  "Taken",
				Color: notionapi.ColorGreen,
			},
		},
		"Note": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: note}},
			},
		},
	}

	request := buildPageCreateRequest(childDbID, &props)
	return client.Page.Create(context.Background(), &request)
}

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
		"Notes": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: note}},
			},
		},
	}
	request := buildPageCreateRequest(dbID, &props)
	return client.Page.Create(context.Background(), &request)
}

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

func lookupSizeColor(size string) notionapi.Color {
	switch size {
	case "Tiny":
		return notionapi.ColorPink
	case "Small":
		return notionapi.ColorRed
	case "Medium":
		return notionapi.ColorOrange
	default:
		return notionapi.ColorYellow
	}
}
