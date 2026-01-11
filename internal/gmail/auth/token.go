package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	scopes []string = []string{
		gmail.MailGoogleComScope,
		gmail.GmailAddonsCurrentMessageReadonlyScope,
		gmail.GmailAddonsCurrentMessageMetadataScope,
	}
)

const (
	credentialsEnvKey = "GEEMAIL_API_CREDENTIALS"
	tokenFile         = "geemail.json"
	tokenFileDir      = ".config"
)

func NewHTTPClient(ctx context.Context) (*http.Client, error) {
	credentials := os.Getenv(credentialsEnvKey)
	if credentials == "" {
		return nil, fmt.Errorf("no credentials found for geemail")
	}
	config, err := google.ConfigFromJSON([]byte(credentials), scopes...)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	tok, err := tokenFromEnv()
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		if err := saveToken(tok); err != nil {
			log.Fatalf("error saving token: %v\n", err)
		}
	}
	return config.Client(context.Background(), tok), nil
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	if err := open(authURL); err != nil {
		log.Fatalf("cannot open URL to redeem token: %v", err)
	}

	authCode := callback()
	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Fatalf("unable to retrieve token from web: %v", err)
	}

	return tok
}

func tokenFromEnv() (*oauth2.Token, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get user home dir: %w", err)
	}
	fpath := path.Join(home, tokenFileDir, tokenFile)
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close() //nolint:errcheck
	var token oauth2.Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func saveToken(token *oauth2.Token) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home dir: %w", err)
	}
	fpath := path.Join(home, tokenFileDir, tokenFile)
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close() //nolint:errcheck
	return json.NewEncoder(f).Encode(token)
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
