package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	envAccessToken        = "ACCESS_TOKEN"
	envRefreshToken       = "REFRESH_TOKEN"
	googleCredentialsFile = "googlecredentials.json"
	tokenEnvFile          = "usertoken.env"
	scopeDriveFile        = "https://www.googleapis.com/auth/drive.file"
)

type userState struct {
	state        string
	codeVerifier string
	code         string
}

var userStates = make(map[string]userState)

var callbackReceivedChan = make(chan string)

func main() {
	// Start callback server
	router := gin.Default()

	router.GET("/oauth/callback", func(ctx *gin.Context) {
		code := ctx.Query("code")
		state := ctx.Query("state")

		us := userStates[state]
		us.code = code
		userStates[state] = us

		callbackReceivedChan <- state
		ctx.Header("Location", "https://t.me/LoQueGastoTestBot")
		ctx.Status(http.StatusFound)
	})

	go router.Run(":8080")

	scopes := []string{scopeDriveFile}
	credentials := getCredentials(googleCredentialsFile)
	oauthConfig, err := google.ConfigFromJSON(credentials, scopes...)
	if err != nil {
		log.Fatalf("Error creating Google Config:\n%+v\n", err)
	}

	//_ = godotenv.Load(tokenEnvFile)
	if os.Getenv(envRefreshToken) == "" {
		loginURL, state := getLoginURL(oauthConfig)

		openBrowser(loginURL)

		waitForCallback()
		token := exchangeCodeWithToken(oauthConfig, state)

		saveToken(token)

		_ = os.Setenv(envAccessToken, token.AccessToken)
		_ = os.Setenv(envRefreshToken, token.RefreshToken)
	}

	ctx := context.Background()
	googleClient := oauthConfig.Client(ctx, &oauth2.Token{
		AccessToken:  os.Getenv(envAccessToken),
		RefreshToken: os.Getenv(envRefreshToken),
	})

	createSpreadsheet(googleClient)
	//docID := createDoc(googleClient)
	//readDoc(googleClient, docID)
	//getAllFiles(googleClient)
}

func getCredentials(fileName string) []byte {
	credentials, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading credentials:\n%+v\n", err)
	}

	return credentials
}
func getLoginURL(c *oauth2.Config) (string, string) {
	state := uuid.New().String()
	codeVerifier := base64.RawURLEncoding.EncodeToString([]byte(uuid.New().String()))
	codeVerifierHash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(codeVerifierHash[:])

	userStates[state] = userState{
		state:        state,
		codeVerifier: codeVerifier,
	}

	authURL := c.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	return authURL, state
}
func openBrowser(url string) {
	err := exec.Command("cmd", "/c", "start", strings.Replace(url, "&", "^&", -1)).Run()
	if err != nil {
		log.Fatalf("Error opening browser:\n%+v\n", err)
	}
}
func waitForCallback() {
	log.Println("Esperando code en callback")
	<-callbackReceivedChan
}
func exchangeCodeWithToken(c *oauth2.Config, state string) *oauth2.Token {
	log.Println("Intercambiando code por token")
	ctx := context.Background()

	us := userStates[state]
	code := us.code
	codeVerifier := us.codeVerifier
	token, err := c.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		log.Fatalf("Error exchangig code:\n%+v\n", err)
	}

	return token
}
func saveToken(token *oauth2.Token) {
	log.Println("Guardando token")
	envFile, _ := os.OpenFile(tokenEnvFile, os.O_WRONLY|os.O_CREATE, 0644)
	defer func() { _ = envFile.Close() }()

	_, _ = fmt.Fprintf(envFile, "ACCESS_TOKEN=%s\nREFRESH_TOKEN=%s", token.AccessToken, token.RefreshToken)

	err := bufio.NewWriter(envFile).Flush()
	if err != nil {
		log.Fatalln(err)
	}
	_ = envFile.Close()
}

func createSpreadsheet(googleClient *http.Client) {
	sheetsSvc, err := sheets.NewService(context.Background(), option.WithHTTPClient(googleClient))
	if err != nil {
		log.Fatalln(err)
	}

	sprdsht, err := sheetsSvc.Spreadsheets.Create(
		&sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: "test-oauth-spreadsheet",
			},
		}).Do()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Spreadsheet created with ID: %s\n", sprdsht.SpreadsheetId)
}
func createDoc(googleClient *http.Client) string {
	docsSvc, err := docs.NewService(context.Background(), option.WithHTTPClient(googleClient))
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := docsSvc.Documents.Create(
		&docs.Document{
			Title: "test-oauth-doc",
		}).Do()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Doc created with ID: %s\n", doc.DocumentId)

	return doc.DocumentId
}
func readDoc(googleClient *http.Client, docID string) {
	docsSvc, err := docs.NewService(context.Background(), option.WithHTTPClient(googleClient))
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := docsSvc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Doc title: %s\n", doc.Title)
	log.Println("Doc body:")
	for _, l := range doc.Body.Content {
		if l.Paragraph == nil {
			continue
		}
		log.Printf("%s", l.Paragraph.Elements[0].TextRun.Content)
	}
}
func getAllFiles(googleClient *http.Client) {
	driveSvc, err := drive.NewService(context.Background(), option.WithHTTPClient(googleClient))
	if err != nil {
		log.Fatalln(err)
	}

	filesList, err := driveSvc.Files.List().Do()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Items:")
	for _, item := range filesList.Items {
		log.Printf(item.Title)
	}
}
