package main

import (
	"context"
	_ "embed"
	"log"

	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CopyTemplate() {
	ctx := context.Background()
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var (
		spreadsheetID = "1bv_mj8-xnNzBGYmF2YqbEwNPz2IyOuZVaD4E4203trc"
		folderID      = "1c04RznMaAumXl9OfVkstH4ZAIG3ULOgR"
	)

	// Replace with the new name for the copied file.
	newFileName := "COPIED_SPREADSHEET_NAME1"

	// Call the Files.Copy method to create a copy of the spreadsheet file.
	copy := &drive.File{
		Title: newFileName,
		Parents: []*drive.ParentReference{
			{
				Id: folderID,
			},
		},
	}
	copiedFile, err := srv.Files.Copy(spreadsheetID, copy).Do()
	if err != nil {
		log.Fatalf("Failed to copy file: %v", err)
	}

	fmt.Printf("Copied file ID: %s\n", copiedFile.Id)

	// Call the Sheets API to retrieve the spreadsheet.
	spreadsheet, err := spreadsheetsSrv.Spreadsheets.Get(copiedFile.Id).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve spreadsheet: %v", err)
	}

	// Loop through each sheet in the spreadsheet.
	for _, sheet := range spreadsheet.Sheets {
		// Loop through each protected range in the sheet.
		for _, protectedRange := range sheet.ProtectedRanges {
			AddServiceAccount(spreadsheetsSrv, copiedFile.Id, protectedRange.ProtectedRangeId)
		}
	}
}

func main() {
	CopyTemplate()
}

func AddServiceAccount(srv *sheets.Service, spreadsheetID string, protectedRangeId int64) {
	// Replace with the ID of the protected range to update.

	// Replace with the email address of the user to add to the protected range.
	newUserEmail := "sheets@secret-beacon-380907.iam.gserviceaccount.com"

	// Define the new editors to be added to the protected range.
	newEditors := &sheets.Editors{
		Users: []string{newUserEmail},
	}

	// Call the Sheets API to modify the protected range's editors.
	requests := []*sheets.Request{
		{
			UpdateProtectedRange: &sheets.UpdateProtectedRangeRequest{
				ProtectedRange: &sheets.ProtectedRange{
					ProtectedRangeId: protectedRangeId,
					Editors:          newEditors,
				},
				Fields: "editors",
			},
		},
	}
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Fatalf("Failed to update protected range: %v", err)
	}
}

/*

	_, err := c.service.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				CreateDeveloperMetadata: &sheets.CreateDeveloperMetadataRequest{
					DeveloperMetadata: &sheets.DeveloperMetadata{
						Location: &sheets.DeveloperMetadataLocation{
							Spreadsheet: true,
						},
						Visibility:    "DOCUMENT",
						MetadataKey:   "token",
						MetadataValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWQiOiIxSTd0WUFoVWpQSkdhTVU3X1hiaEMwOHJRdzU1SVJjN2JFdGcxbWdtUlBLZyJ9.7yAZuGAm7_WSkGJURMSn5aS8UacVAY-CPx-vOO0rPDE",
					},
				},
			},
		},
	}).Do()

	resp, err := c.service.Spreadsheets.DeveloperMetadata.Search(spreadsheetID, &sheets.SearchDeveloperMetadataRequest{
		DataFilters: []*sheets.DataFilter{
			{
				DeveloperMetadataLookup: &sheets.DeveloperMetadataLookup{
					MetadataKey: "token",
				},
			},
		},
	}).Do()

*/
