package main

import (
	"context"
	_ "embed"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func ProtectedRange() {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var (
		spreadsheetID = "1KL-lrhs-Wu9kRAppBxAHUUFr7OCfNYla8Z7W-0tX4Mo"
		// folderID      = "1c04RznMaAumXl9OfVkstH4ZAIG3ULOgR"
		// fileID = "1vV43-P9I2dXZ3KkuqviBW9IRImEdIWA8qVTOG7Tm4G0"
	)

	// // Call the Sheets API to retrieve the spreadsheet.
	// spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve spreadsheet: %v", err)
	// }

	// // Loop through each sheet in the spreadsheet.
	// for _, sheet := range spreadsheet.Sheets {
	// 	// Loop through each protected range in the sheet.
	// 	for _, protectedRange := range sheet.ProtectedRanges {
	// 		fmt.Printf("Protected range on sheet %s: %#v\n", sheet.Properties.Title, protectedRange)
	// 	}
	// }

	// Replace with the ID of the protected range to update.
	var protectedRangeId int64 = 1017731739

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
	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Fatalf("Failed to update protected range: %v", err)
	}
}
