package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
    "runtime"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

const (
    tokenFile = "config/token.json"
)

func GenerateToken(credentials []byte) error {
    config, err := google.ConfigFromJSON(credentials, scopes...)
    if err != nil {
        return fmt.Errorf("cannot read config file: %w", err)
    }

    tok, err := tokenFromFile(tokenFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        return saveToken(tokenFile, tok)
    }
    return nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

    open(authURL)

    authCode := callback()
    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web: %v", err)
    }

    return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    if envFile := os.Getenv("TOKEN_FILE"); envFile != "" {
        file = envFile
    }
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var token oauth2.Token
    if err := json.NewDecoder(f).Decode(&token); err != nil {
        return nil, err
    }
    return &token, nil
}

func saveToken(path string, token *oauth2.Token) error {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return fmt.Errorf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
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
