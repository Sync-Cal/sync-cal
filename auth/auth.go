package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	// "golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, token *oauth2.Token) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	// tok := getToken(email, clientId, clientSecret)
	// fmt.Println(tok)

	var client *http.Client
	client = config.Client(context.Background(), token)
	return client
}

func getTokens(clientId string, clientSecret string) (map[string]oauth2.Token, error) {
	f, err := os.Open("../data/tokens.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &map[string]oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	for email, t := range *tok {
		res, err := http.PostForm(
			"https://www.googleapis.com/oauth2/v4/token",
			url.Values{
				"client_id":     {clientId},
				"client_secret": {clientSecret},
				"refresh_token": {t.RefreshToken},
				"grant_type":    {"refresh_token"},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("ERROR: %f", err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("ERROR: %f", err)
		}

		jsonMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(string(body)), &jsonMap)

		(*tok)[email] = *&oauth2.Token{
			AccessToken:  fmt.Sprintf("%v", jsonMap["access_token"]),
			TokenType:    t.TokenType,
			RefreshToken: t.RefreshToken,
			Expiry:       t.Expiry,
		}

	}

	return *tok, err
}

func GetServices(emails []string, clientId string, clientSecret string) map[string]*calendar.Service {
	ctx := context.Background()
	b, err := os.ReadFile("../data/credentials.json")

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	log.Println(config)

	// config := &oauth2.Config{ClientID: clientId, ClientSecret: clientSecret}
	tokens, _ := getTokens(clientId, clientSecret)
	services := map[string]*calendar.Service{}

	for _, email := range emails {
		if token, ok := tokens[email]; ok {
			fmt.Println(email, token.AccessToken)
			cl := getClient(config, &token)
			srv, err := calendar.NewService(ctx, option.WithHTTPClient(cl))
			if err != nil {
				log.Fatalf("Unable to retrieve Calendar client: %v", err)
			}
			services[email] = srv
		} else {
			fmt.Println("No token for %s", email)
		}

	}

	return services
}
