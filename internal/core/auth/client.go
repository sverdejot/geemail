package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getClientFromFile(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		log.Fatalf("cannot read token, please, regenerate it using 'generate-token': %v", err)
	}
	return config.Client(context.Background(), tok)
}

func NewHTTPClient(credentials []byte) (*http.Client, error) {
	config, err := google.ConfigFromJSON(credentials, scopes...)
	if err != nil {
		return nil, fmt.Errorf("failed creating default HTTP client: %w", err)
	}
	return getClientFromFile(config), nil
}
