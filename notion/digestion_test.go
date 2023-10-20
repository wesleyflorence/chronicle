package notion

import (
	"fmt"
	"testing"
	"time"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppendDigestionEntry(t *testing.T) {
	currentTime := notionapi.Date(time.Now())
	dbID := "testDbId"

	mockPageService := new(MockPageService)
	client := &notionapi.Client{
		Page: mockPageService,
	}

	mockedResponse := &notionapi.Page{
		Object: "page",
		ID:     "id",
		Properties: notionapi.Properties{
			"Bristol (1-7)": &notionapi.NumberProperty{
				Number: 4,
			},
			"Chronicle": &notionapi.CheckboxProperty{
				Checkbox: true,
			},
			"Date": &notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: &currentTime,
				},
			},
			"Note": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{Content: "Test Note"}},
				},
			},
			"Size": &notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: "Medium",
				},
			},
		},
	}

	testCases := []struct {
		name         string
		bristol      int
		size         string
		note         string
		expectError  bool
		expectedPage *notionapi.Page
	}{
		{
			name:         "Successful entry",
			bristol:      4,
			size:         "Medium",
			note:         "Test Note",
			expectError:  false,
			expectedPage: mockedResponse,
		},
		{
			name:        "Error from client",
			bristol:     4,
			size:        "Medium",
			note:        "Test Note",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectError {
				mockPageService.On("Create", mock.Anything, mock.AnythingOfType("*notionapi.PageCreateRequest")).Return(nil, fmt.Errorf("mock error")).Once()
			} else {
				mockPageService.On("Create", mock.Anything, mock.AnythingOfType("*notionapi.PageCreateRequest")).Return(mockedResponse, nil).Once()
			}
			page, err := AppendDigestionEntry(client, dbID, tc.bristol, tc.size, tc.note)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPage, page)
			}
			// Assert that PageService was called
			mockPageService.AssertExpectations(t)
		})
	}
}
