package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"loquegasto-telegram/internal/defines"
	"net/http"
	"os"
	"time"
)

type SheetsService interface {
	AddRow(date time.Time, description string, amount float64, payerName string) (*sheets.AppendValuesResponse, error)
	GetSpreadsheetID() string
}

type sheetsService struct {
	repo          *sheets.Service
	spreadsheetID string
}

func NewSheetsService() SheetsService {
	var srv sheetsService

	ctx := context.Background()
	b := os.Getenv("SHEETS_CONFIG")

	config, err := google.ConfigFromJSON([]byte(b), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv.repo, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	srv.spreadsheetID = os.Getenv(defines.EnvSheetsSpreadsheetID)

	return &srv
}

func (s *sheetsService) AddRow(date time.Time, description string, amount float64, payerName string) (*sheets.AppendValuesResponse, error) {
	writeRange := "Gastos!A2:D"
	var values [][]interface{}
	values = append(values, []interface{}{date.Format("02/01/2006"), description, amount, payerName})
	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          writeRange,
		Values:         values,
	}

	resp, err := s.repo.Spreadsheets.Values.Append(s.spreadsheetID, writeRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil || resp.HTTPStatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed to add to Google Sheet (%d - %s)", resp.HTTPStatusCode, err.Error()))
	}

	return resp, nil
}

func (s *sheetsService) GetSpreadsheetID() string {
	return s.spreadsheetID
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromEnv()
	if err != nil {
		tok = getTokenFromWeb(config)
		if err := saveToken("sheetsToken.json", tok); err != nil {
			log.Fatalln(err)
		}
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

// Retrieves a token from a env.
func tokenFromEnv() (*oauth2.Token, error) {
	token := os.Getenv(defines.EnvSheetsToken)
	if token == "" {
		return nil, errors.New(defines.EnvSheetsToken + " empty")
	}
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(token), tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	return json.NewEncoder(f).Encode(token)
}
