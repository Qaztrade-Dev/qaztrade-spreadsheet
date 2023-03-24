package main

import (
	"context"
	_ "embed"
	"log"

	"fmt"

	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

//go:embed credentials.json
var credentials []byte

func main() {
	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var (
		// spreadsheetID = "1bv_mj8-xnNzBGYmF2YqbEwNPz2IyOuZVaD4E4203trc"
		// folderID      = "1c04RznMaAumXl9OfVkstH4ZAIG3ULOgR"
		fileID = "1vV43-P9I2dXZ3KkuqviBW9IRImEdIWA8qVTOG7Tm4G0"
	)

	permission := &drive.Permission{
		Type:     "anyone",
		Role:     "writer",
		WithLink: true,
	}
	res, err := srv.Permissions.Insert(fileID, permission).Do()
	if err != nil {
		log.Fatalf("Failed to create permission: %v", err)
	}
	fmt.Printf("%#v\n", res)

	fmt.Printf("Added owner permission to file: %s\n", fileID)
}
