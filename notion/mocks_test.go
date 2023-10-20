package notion

import (
	"context"
	"time"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/mock"
)

type MockPageService struct {
	mock.Mock
}

func (m *MockPageService) Create(ctx context.Context, request *notionapi.PageCreateRequest) (*notionapi.Page, error) {
	args := m.Called(ctx, request)
	var page *notionapi.Page
	if args.Get(0) != nil {
		page = args.Get(0).(*notionapi.Page)
	}
	return page, args.Error(1)
}

func (m *MockPageService) Get(ctx context.Context, pageID notionapi.PageID) (*notionapi.Page, error) {
	args := m.Called(ctx, pageID)
	return args.Get(0).(*notionapi.Page), args.Error(1)
}

func (m *MockPageService) Update(ctx context.Context, pageID notionapi.PageID, request *notionapi.PageUpdateRequest) (*notionapi.Page, error) {
	args := m.Called(ctx, pageID, request)
	return args.Get(0).(*notionapi.Page), args.Error(1)
}

type MockBlockService struct {
	mock.Mock
}

func (m *MockBlockService) AppendChildren(ctx context.Context, blockID notionapi.BlockID, request *notionapi.AppendBlockChildrenRequest) (*notionapi.AppendBlockChildrenResponse, error) {
	args := m.Called(ctx, blockID, request)
	return args.Get(0).(*notionapi.AppendBlockChildrenResponse), args.Error(1)
}

func (m *MockBlockService) Get(ctx context.Context, blockID notionapi.BlockID) (notionapi.Block, error) {
	args := m.Called(ctx, blockID)
	return args.Get(0).(notionapi.Block), args.Error(1)
}

func (m *MockBlockService) GetChildren(ctx context.Context, blockID notionapi.BlockID, pagination *notionapi.Pagination) (*notionapi.GetChildrenResponse, error) {
	args := m.Called(ctx, blockID, pagination)
	return args.Get(0).(*notionapi.GetChildrenResponse), args.Error(1)
}

func (m *MockBlockService) Update(ctx context.Context, id notionapi.BlockID, request *notionapi.BlockUpdateRequest) (notionapi.Block, error) {
	args := m.Called(ctx, id, request)
	return args.Get(0).(notionapi.Block), args.Error(1)
}

func (m *MockBlockService) Delete(ctx context.Context, blockID notionapi.BlockID) (notionapi.Block, error) {
	args := m.Called(ctx, blockID)
	return args.Get(0).(notionapi.Block), args.Error(1)
}

type MockDatabaseService struct {
	mock.Mock
}

func (m *MockDatabaseService) Create(ctx context.Context, request *notionapi.DatabaseCreateRequest) (*notionapi.Database, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*notionapi.Database), args.Error(1)
}

func (m *MockDatabaseService) Query(ctx context.Context, databaseID notionapi.DatabaseID, request *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
	args := m.Called(ctx, databaseID, request)
	return args.Get(0).(*notionapi.DatabaseQueryResponse), args.Error(1)
}

func (m *MockDatabaseService) Get(ctx context.Context, databaseID notionapi.DatabaseID) (*notionapi.Database, error) {
	args := m.Called(ctx, databaseID)
	return args.Get(0).(*notionapi.Database), args.Error(1)
}

func (m *MockDatabaseService) Update(ctx context.Context, databaseID notionapi.DatabaseID, request *notionapi.DatabaseUpdateRequest) (*notionapi.Database, error) {
	args := m.Called(ctx, databaseID, request)
	return args.Get(0).(*notionapi.Database), args.Error(1)
}

func MockNotionMedicinePage() *notionapi.Page {
	createdTime, _ := time.Parse(time.RFC3339, "2023-10-20T06:42:00Z")
	lastEditedTime, _ := time.Parse(time.RFC3339, "2023-10-20T06:42:00Z")

	return &notionapi.Page{
		Object:         "page",
		ID:             "userid",
		CreatedTime:    createdTime,
		LastEditedTime: lastEditedTime,
		CreatedBy: notionapi.User{
			Object: "user",
			ID:     "userid",
		},
		LastEditedBy: notionapi.User{
			Object: "user",
			ID:     "userid",
		},
		Archived: false,
		Properties: map[string]notionapi.Property{
			"Chronicle": &notionapi.CheckboxProperty{
				Checkbox: true,
			},
			"Date": &notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: &notionapi.Date{},
				},
			},
			"Dose": &notionapi.NumberProperty{
				Number: 1,
			},
			"Note": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{Content: "Test Note"}},
				},
			},
			"Status": &notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: "Taken",
				},
			},
		},
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: "dbId",
		},
		URL: "https://www.notion.so/medId",
	}
}

func MockBlockChildren() *notionapi.GetChildrenResponse {
	createdTime, _ := time.Parse(time.RFC3339, "2023-08-22T05:47:00Z")
	lastEditedTime, _ := time.Parse(time.RFC3339, "2023-08-22T06:04:00Z")

	child := &notionapi.ChildPageBlock{
		BasicBlock: notionapi.BasicBlock{
			Object:         "block",
			ID:             "id",
			Type:           "child_page",
			CreatedTime:    (*time.Time)(&createdTime),
			LastEditedTime: (*time.Time)(&lastEditedTime),
			// ... other necessary fields
		},
		ChildPage: struct {
			Title string "json:\"title\""
		}{
			Title: "Ondansetron",
		},
	}

	return &notionapi.GetChildrenResponse{
		Object: "list",
		Results: []notionapi.Block{
			child,
		},
	}
}

func MockDatabaseChildren() *notionapi.GetChildrenResponse {
	createdTime, _ := time.Parse(time.RFC3339, "2023-08-22T05:47:00Z")
	lastEditedTime, _ := time.Parse(time.RFC3339, "2023-08-22T06:04:00Z")

	child := &notionapi.ChildDatabaseBlock{
		BasicBlock: notionapi.BasicBlock{
			Object:         "block",
			ID:             "id",
			Type:           "child_page",
			CreatedTime:    (*time.Time)(&createdTime),
			LastEditedTime: (*time.Time)(&lastEditedTime),
		},
		ChildDatabase: struct {
			Title string "json:\"title\""
		}{
			Title: "Taken",
		},
	}

	return &notionapi.GetChildrenResponse{
		Object: "list",
		Results: []notionapi.Block{
			child,
		},
	}
}

func MockDatabaseQueryResponse() *notionapi.DatabaseQueryResponse {
	createdTime, _ := time.Parse(time.RFC3339, "2023-10-05T17:49:00Z")
	lastEditedTime, _ := time.Parse(time.RFC3339, "2023-10-05T17:49:00Z")

	page := notionapi.Page{
		Object:         "page",
		ID:             "id",
		CreatedTime:    createdTime,
		LastEditedTime: lastEditedTime,
		CreatedBy: notionapi.User{
			Object: "user",
			ID:     "id",
		},
		LastEditedBy: notionapi.User{
			Object: "user",
			ID:     "id",
		},
		Archived: false,
		Properties: notionapi.Properties{
			"Chronicle": &notionapi.CheckboxProperty{
				Checkbox: true,
			},
			"Date": &notionapi.DateProperty{
				Date: &notionapi.DateObject{},
			},
			"Dose": &notionapi.NumberProperty{
				Number: 12,
			},
			"Note": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: "Some note content",
						},
					},
				},
			},
			"Status": &notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: "Taken",
				},
			},
		},
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: "dbId",
		},
		URL:       "https://www.notion.so/id",
		PublicURL: "",
	}

	return &notionapi.DatabaseQueryResponse{
		Object:  "list",
		Results: []notionapi.Page{page},
		HasMore: false,
	}
}
