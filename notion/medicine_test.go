package notion

import (
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppendMedicineEntry(t *testing.T) {
	mockPageService := new(MockPageService)
	mockBlockService := new(MockBlockService)
	mockDatabaseService := new(MockDatabaseService)
	client := &notionapi.Client{
		Page:     mockPageService,
		Block:    mockBlockService,
		Database: mockDatabaseService,
	}

	mockedResponse := MockNotionMedicinePage()

	testCases := []struct {
		name         string
		medicine     string
		size         int
		note         string
		expectError  bool
		expectedPage *notionapi.Page
	}{
		{
			name:         "Successful entry",
			medicine:     "Ondansetron",
			size:         -1,
			note:         "Test Note",
			expectError:  false,
			expectedPage: mockedResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mocking GetChildren Block and Database methods
			mockBlockService.On("GetChildren", mock.Anything, mock.AnythingOfType("notionapi.BlockID"), mock.AnythingOfType("*notionapi.Pagination")).Return(MockBlockChildren(), nil).Once()
			mockBlockService.On("GetChildren", mock.Anything, mock.AnythingOfType("notionapi.BlockID"), mock.AnythingOfType("*notionapi.Pagination")).Return(MockDatabaseChildren(), nil).Once()

			// Mocking Query method
			mockDatabaseService.On("Query", mock.Anything, mock.AnythingOfType("notionapi.DatabaseID"), mock.AnythingOfType("*notionapi.DatabaseQueryRequest")).Return(MockDatabaseQueryResponse(), nil).Once()

			// Mocking Create method
			mockPageService.On("Create", mock.Anything, mock.AnythingOfType("*notionapi.PageCreateRequest")).Return(MockNotionMedicinePage(), nil).Once()

			// Call AppendMedicineEntry
			page, err := AppendMedicineEntry(client, "testMedicinePageID", tc.medicine, tc.size, tc.note)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPage, page)
			}

			// Assert that all expected calls were made
			mockBlockService.AssertExpectations(t)
			mockDatabaseService.AssertExpectations(t)
			mockPageService.AssertExpectations(t)
		})
	}
}
