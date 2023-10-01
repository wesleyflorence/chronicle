package notion

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
)

func setUp() (*notionapi.Client, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	digestionDbID := os.Getenv("DIGESTION_DB")
	return notionapi.NewClient(notionapi.Token(notionAPIKey)), digestionDbID
}

func TestAppendDigestionEntry(t *testing.T) {
	client, dbID := setUp()
	tests := []struct {
		name    string
		bristol int
		size    string
		notes   string
		wantErr bool
	}{
		{
			name:    "test1",
			bristol: 4,
			size:    "Large",
			notes:   "Sent from Chronicle",
			wantErr: false,
		},
		{
			name:    "test2",
			bristol: 5,
			size:    "Small",
			notes:   "Sent from Chronicle",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendDigestionEntry(client, dbID, tt.bristol, tt.size, tt.notes)
			if err != nil {
				t.Errorf("AppendDigestionEntry() error = %v", err)
				return
			}
			currentTime := time.Now().Format("2006-01-02")
			if !got.CreatedTime.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour)) {
				t.Errorf("AppendDigestionEntry() = %v, want created time to be %v", got.CreatedTime, currentTime)
			}
		})
	}
}
