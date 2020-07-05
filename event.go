package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, save the token and return the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// token.json stores the user's access and refresh tokens
	// created automatically when the authorization flow completes for the first time
	tokFile := "token.json"
	tok, err := getTokenFromFile(tokFile) // Returns an error if the file is absent
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then copy and paste the "+
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
func getTokenFromFile(file string) (*oauth2.Token, error) {
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

func main() {
	jsonKey, err := ioutil.ReadFile("credentials.json")
	
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Get a oauth2 configuration object with the jsonKey
	config, err := google.ConfigFromJSON(jsonKey, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	event := &calendar.Event{
		Summary: "Vehicle Papers Renewal",
		Start: &calendar.EventDateTime{
			Date:     "2020-07-22",
			TimeZone: "Africa/Lagos",
		},
		End: &calendar.EventDateTime{
			Date:     "2020-10-22",
			TimeZone: "Africa/Lagos",
		},
	}

	// calender identifier: can be the emailaddress of the calender to be used
	//  or a special keyword, "primary" which will use the primary calender of the logged i user
	calendarID := "primary"
	event, err = srv.Events.Insert(calendarID, event).Do()

	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}
