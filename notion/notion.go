// Package Notion for interacting with notion api
package notion

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
)

const testDbID = "632816ea9920439dbd1019c91fc5fa15"

func getAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("NOTION_API_KEY")
}

func addRowToDB() (*notionapi.Page, error) {
	client := notionapi.NewClient(notionapi.Token(getAPIKey()))
	request := notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(testDbID),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: "test name from golang"}},
				},
			},
			"Tags": notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: "done", Color: notionapi.ColorPurple},
				},
			},
		},
	}
	return client.Page.Create(context.Background(), &request)
}
