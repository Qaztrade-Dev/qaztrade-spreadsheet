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
)

// func main() {
// 	ctx := context.Background()

// 	creds, err := google.CredentialsFromJSON(ctx, data, drive.DriveScope)
// 	if err != nil {
// 		log.Fatalf("Failed to parse credentials file: %v", err)
// 	}

// 	srv, err := drive.NewService(ctx, option.WithCredentialsJSON(credentials))
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Sheets client: %v", err)
// 	}

// var (
// 	spreadsheetID = "1bv_mj8-xnNzBGYmF2YqbEwNPz2IyOuZVaD4E4203trc"
// 	folderID      = "1c04RznMaAumXl9OfVkstH4ZAIG3ULOgR"
// )

// Replace with the new name for the copied file.
// newFileName := "COPIED_SPREADSHEET_NAME"

// Call the Files.Copy method to create a copy of the spreadsheet file.
// copy := &drive.File{
// 	Title: newFileName,
// 	Parents: []*drive.ParentReference{
// 		{
// 			Id: folderID,
// 		},
// 	},
// }
// copiedFile, err := srv.Files.Copy(spreadsheetID, copy).Do()
// if err != nil {
// 	log.Fatalf("Failed to copy file: %v", err)
// }

// fmt.Printf("Copied file ID: %s\n", copiedFile.Id)

// // Add a new permission for the owner of the service account to the copied file.
// permission := &drive.Permission{
// 	Type:  "user",
// 	Role:  "owner",
// 	Value: "ali.tlekbai@gmail.com",
// }
// res, err := srv.Permissions.Insert(copiedFile.Id, permission).Do()
// if err != nil {
// 	log.Fatalf("Failed to create permission: %v", err)
// }

// fmt.Printf("Added owner permission to file: %s\n", copiedFile.Id)

// spreadsheetClient, err := cli.NewSpreadsheetClient(ctx, spreadsheetID)
// if err != nil {
// 	log.Fatalf("NewSpreadsheetClient error: %v", err)
// }

// err = spreadsheetClient.SubmitRow(ctx, payload)
// if err != nil {
// 	log.Fatalf("SubmitRow error: %v", err)
// }
// }

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

// func main() {
// 	ctx := context.Background()
// 	b, err := os.ReadFile("client_secret.json")
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// If modifying these scopes, delete your previously saved token.json.
// 	config, err := google.ConfigFromJSON(b, drive.DriveScope)
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	client := getClient(config)

// 	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Sheets client: %v", err)
// 	}

// 	var (
// 		fileID = "1vV43-P9I2dXZ3KkuqviBW9IRImEdIWA8qVTOG7Tm4G0"
// 		// spreadsheetID = "1bv_mj8-xnNzBGYmF2YqbEwNPz2IyOuZVaD4E4203trc"
// 		// folderID      = "1c04RznMaAumXl9OfVkstH4ZAIG3ULOgR"
// 	)

// 	// Replace with the new name for the copied file.
// 	// newFileName := "COPIED_SPREADSHEET_NAME"

// 	// // Call the Files.Copy method to create a copy of the spreadsheet file.
// 	// copy := &drive.File{
// 	// 	Title: newFileName,
// 	// 	Parents: []*drive.ParentReference{
// 	// 		{
// 	// 			Id: folderID,
// 	// 		},
// 	// 	},
// 	// }
// 	// copiedFile, err := srv.Files.Copy(spreadsheetID, copy).Do()
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to copy file: %v", err)
// 	// }

// 	// fmt.Printf("Copied file ID: %s\n", copiedFile.Id)
// }
