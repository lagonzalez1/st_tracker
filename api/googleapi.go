package api

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func ConnectSheetsAPI() {
	ctx := context.Background()

	cred_json, err := os.ReadFile("/Users/luisgonzalez/Documents/Keys/fleet-rhino-387703-83400b19b794.json")
	if err != nil {
		fmt.Printf("Unable to open file fleet-rhino-387703-83400b19b794.json")
	}

	creds, err := google.CredentialsFromJSON(ctx, cred_json, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to validate creds: %v", err)
	}
	sheetsService, err := sheets.NewService(ctx, option.WithCredentials(creds))

	if err != nil {
		fmt.Printf("Unable to create Sheets service: %v", err)
	}
	// Use the Sheets service to fetch data
	spreadsheetId := "1JKHLw46NozTj-9XhA84B1x75iVGa592hFkpmGdkVAwI"
	readRange := "Sheet1!B12:F59" // Adjust the range as needed

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve data from sheet: %v", err)
	}

	// Print the data
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			fmt.Println(row)
		}
	}
}
