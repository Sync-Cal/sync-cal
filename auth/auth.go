package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (token *AuthToken) Refresh(clientId string, clientSecret string) error {
	res, err := http.PostForm(
		"https://www.googleapis.com/oauth2/v4/token",
		url.Values{
			"client_id":     {clientId},
			"client_secret": {clientSecret},
			"refresh_token": {token.RefreshToken},
			"grant_type":    {"refresh_token"},
		},
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(string(body)), &jsonMap)
	token.AccessToken = fmt.Sprintf("%v", jsonMap["access_token"])
	return nil
}

func (t AuthToken) ToOAuth2() oauth2.Token {
	return oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, token AuthToken) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	var client *http.Client
	t := token.ToOAuth2()
	client = config.Client(context.Background(), &t)
	return client
}

func GetServices(tokensMap map[string]AuthToken, credentialsFilePath string) map[string]*calendar.Service {
	ctx := context.Background()
	b, err := os.ReadFile(credentialsFilePath)

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	services := map[string]*calendar.Service{}

	for email, token := range tokensMap {
		cl := getClient(config, token)
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(cl))
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		services[email] = srv
	}

	return services
}
