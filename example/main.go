package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tenz-io/notionapi"

	"github.com/tenz-io/gokit/httpcli"
	"github.com/tenz-io/gokit/logger"
)

var (
	envToken      = os.Getenv("NOTION_TOKEN")
	envPageID     = os.Getenv("NOTION_PAGE_ID")
	envDatabaseID = os.Getenv("NOTION_DB_ID")
	newPageID     = ""
)

func init() {
	logger.ConfigureWithOpts(
		logger.WithLoggerLevel(logger.DebugLevel),
		logger.WithSetAsDefaultLvl(true),
		logger.WithConsoleEnabled(true),
		logger.WithCallerEnabled(true),
		logger.WithCallerSkip(1),
	)
	logger.ConfigureTrafficWithOpts(
		logger.WithTrafficEnabled(true),
	)
}

func main() {
	defer func() {
		time.Sleep(100 * time.Millisecond)
	}()

	loadEnvVars()

	notionClient := newNotionClient(envToken)

	// test page
	testGetPage(notionClient, envPageID)

	// test create database
	if envDatabaseID == "" {
		testCreateDatabase(notionClient, envPageID)
	}

	// test get database
	testGetDatabase(notionClient, envDatabaseID)

	// test update database
	testUpdateDatabase(notionClient, envDatabaseID)

	// test create page in database
	testCreatePageInDatabase(notionClient, envDatabaseID)

	// test query database
	testQueryDatabase(notionClient, envDatabaseID)
}

// load env vars
func loadEnvVars() {
	// load os env from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envToken = os.Getenv("NOTION_TOKEN")
	envPageID = os.Getenv("NOTION_PAGE_ID")
	envDatabaseID = os.Getenv("NOTION_DB_ID")

	fmt.Sprintln("env vars loaded:")
	fmt.Sprintln("NOTION_TOKEN:", envToken)
	fmt.Sprintln("NOTION_PAGE_ID:", envPageID)
	fmt.Sprintln("NOTION_DB_ID:", envDatabaseID)

	if envToken == "" {
		panic("NOTION_TOKEN is required")
	}

	if envPageID == "" {
		panic("NOTION_PAGE_ID is required")
	}

}

// newNotionClient returns a new notion client
func newNotionClient(token string) *notionapi.Client {
	hc := &http.Client{}
	interceptor := httpcli.NewInterceptorWithOpts(
		httpcli.WithEnableTraffic(true),
	)
	interceptor.Apply(hc)

	return notionapi.NewClient(notionapi.Token(token),
		notionapi.WithHTTPClient(hc),
		notionapi.WithVersion("2022-06-28"),
	)
}

// testGetPage
func testGetPage(notionClient *notionapi.Client, pageID string) {
	page, err := notionClient.Page.Get(context.Background(), notionapi.PageID(pageID))
	if err != nil {
		panic(err)
		return
	}
	logger.WithFields(logger.Fields{
		"page": page,
	}).Infof("page")
}

// testCreateDatabase creates a new database
func testCreateDatabase(notionClient *notionapi.Client, pageId string) {
	db, err := notionClient.Database.Create(context.Background(), &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(pageId),
		},
		Title: []notionapi.RichText{
			{
				Type: notionapi.ObjectTypeText,
				Text: &notionapi.Text{
					Content: "My Database",
				},
				PlainText: "My Database",
			},
		},
		Properties: notionapi.PropertyConfigs{
			"Name": notionapi.TitlePropertyConfig{
				Type: notionapi.PropertyConfigTypeTitle,
			},
			"Description": notionapi.RichTextPropertyConfig{
				Type: notionapi.PropertyConfigTypeRichText,
			},
			"Status": notionapi.SelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeSelect,
				Select: notionapi.Select{
					Options: []notionapi.Option{
						{
							Name:  "TODO",
							Color: "gray",
						},
						{
							Name:  "DOING",
							Color: "blue",
						},
						{
							Name:  "DONE",
							Color: "green",
						},
					},
				},
			},
			"Date": notionapi.DatePropertyConfig{
				Type: notionapi.PropertyConfigTypeDate,
			},
			"Tags": notionapi.MultiSelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeMultiSelect,
				MultiSelect: notionapi.Select{
					Options: []notionapi.Option{
						{
							Name: "Tag 1",
						},
						{
							Name: "Tag 2",
						},
					},
				},
			},
		},
	})

	if err != nil {
		panic(err)
		return
	}
	logger.WithFields(logger.Fields{
		"ID": db.ID,
		"db": db,
	}).Infof("db")

	envDatabaseID = string(db.ID)
}

// testGetDatabase
func testGetDatabase(notionClient *notionapi.Client, dbID string) {
	db, err := notionClient.Database.Get(context.Background(), notionapi.DatabaseID(dbID))
	if err != nil {
		panic(err)
		return
	}
	logger.WithFields(logger.Fields{
		"db": db,
	}).Infof("db")
}

// testUpdateDatabase
func testUpdateDatabase(notionClient *notionapi.Client, dbID string) {
	db, err := notionClient.Database.Update(context.Background(), notionapi.DatabaseID(dbID), &notionapi.DatabaseUpdateRequest{
		Title: []notionapi.RichText{
			{
				Type: notionapi.ObjectTypeText,
				Text: &notionapi.Text{
					Content: "My Database Updated",
				},
				PlainText: "My Database Updated",
			},
		},
	})

	if err != nil {
		panic(err)
		return
	}
	logger.WithFields(logger.Fields{
		"db": db,
	}).Infof("db")
}

// testQueryDatabase
func testQueryDatabase(notionClient *notionapi.Client, dbID string) {
	queryResp, err := notionClient.Database.Query(context.Background(), notionapi.DatabaseID(dbID), &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Name",
			RichText: &notionapi.TextFilterCondition{
				Contains: "Test",
			},
		},
		PageSize: 10,
	})

	if err != nil {
		panic(err)
		return
	}
	logger.WithFields(logger.Fields{
		"pages":          len(queryResp.Results),
		"query_response": queryResp,
	}).Infof("db query result")
}

// testCreatePageInDatabase
func testCreatePageInDatabase(notionClient *notionapi.Client, dbID string) {
	timeObj, err := time.Parse(time.RFC3339, "2020-12-08T12:00:00Z")
	if err != nil {
		panic(err)
		return
	}

	dateObj := notionapi.Date(timeObj)

	page, err := notionClient.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(dbID),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: "Test Page",
						},
						PlainText: "Test Page",
					},
				},
			},
			"Description": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: "This is a test page",
						},
						PlainText: "This is a test page",
					},
				},
			},
			"Status": notionapi.SelectProperty{
				Select: notionapi.Option{
					Name:  "TODO",
					Color: "gray",
				},
			},
			"Date": notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: &dateObj,
				},
			},
			"Tags": notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{
						Name: "Tag 1",
					},
				},
			},
		},
		Children: notionapi.Blocks{
			&notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							Type: notionapi.ObjectTypeText,
							Text: &notionapi.Text{
								Content: "This is a test paragraph 1",
							},
						},
					},
				},
			},
			&notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							Type: notionapi.ObjectTypeText,
							Text: &notionapi.Text{
								Content: "This is a test paragraph 2",
								Link: &notionapi.Link{
									Url: "https://example.com",
								},
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		panic(err)
		return
	}
	logger.WithFields(logger.Fields{
		"ID":   page.ID,
		"page": page,
	}).Infof("page")
	newPageID = string(page.ID)
}
