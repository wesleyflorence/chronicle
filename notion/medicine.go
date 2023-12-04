// Package notion for interacting with notion api
package notion

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jomei/notionapi"
)

// AppendMedicineEntry appends a medicine entry to a Notion database.
// It takes a Notion client, medicine page ID, medicine name, and note as input.
// It retrieves the medicine block ID, medicine database ID, and dose from the Notion client.
// Then, it appends a new row to the medicine database with the current date, dose, status, note, and chronicle properties.
// Finally, it returns the created page or an error if any occurred.
func AppendMedicineEntry(client *notionapi.Client, medicinePageID, medicine string, size int, note string) (*notionapi.Page, error) {
	medicineBlockID, err := getMedicineBlockID(client, medicinePageID, medicine)
	if err != nil {
		log.Println("Error retrieving Medicine Block ID")
		return nil, err
	}

	dbID, err := getChildDbID(client, medicineBlockID)
	if err != nil {
		log.Println("Error retrieving Medicine DB ID")
		return nil, err
	}

	dose, err := getDose(client, *dbID)
	if err != nil {
		log.Println("Error retrieving Dose")
		return nil, err
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
		"Chronicle": notionapi.CheckboxProperty{
			Checkbox: true,
		},
	}

	if medicine == "Milk of Magnesia" {
		props["mL"] = notionapi.NumberProperty{
			Number: float64(size)}
	}

	request := buildPageCreateRequest(dbID.String(), &props)
	return client.Page.Create(context.Background(), &request)
}

// Get Block form Medicine Page
func getMedicineBlockID(client *notionapi.Client, medicinePageID, medicine string) (*notionapi.BlockID, error) {
	medMap := make(map[string]string)
	children, err := client.Block.GetChildren(context.Background(), notionapi.BlockID(
		medicinePageID), &notionapi.Pagination{PageSize: 10})
	if err != nil {
		log.Println("Error Calling GetChildren")
		return nil, err
	}

	// Populate the Map using Medicine Child Page Titles
	for _, child := range children.Results {
		if block, ok := child.(*notionapi.ChildPageBlock); ok {
			medMap[block.ChildPage.Title] = block.ID.String()
		}
	}
	if _, ok := medMap[medicine]; !ok {
		return nil, fmt.Errorf("%s not found in the Medicine Page", medicine)
	}

	id := notionapi.BlockID(medMap[medicine])
	return &id, nil
}

func getChildDbID(client *notionapi.Client, blockID *notionapi.BlockID) (*notionapi.BlockID, error) {
	// Get the Blocks from the provided medicine's Page
	medicinePage, err := client.Block.GetChildren(context.Background(), *blockID,
		&notionapi.Pagination{PageSize: 100})
	if err != nil {
		log.Println("Error Calling GetChildren")
		return nil, err
	}

	// Get the Database ID
	for _, block := range medicinePage.Results {
		if block, ok := block.(*notionapi.ChildDatabaseBlock); ok {
			return &block.ID, nil
		}
	}
	return nil, fmt.Errorf("No DB found on page")
}

func getDose(client *notionapi.Client, childDbID notionapi.BlockID) (int, error) {
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
			log.Println("Error Making Query to get Count")
			return -1, err
		}
		dose = dose + len(queryResult.Results)
		cursor = &queryResult.NextCursor
		hasMore = queryResult.HasMore
	}
	return dose, nil
}
