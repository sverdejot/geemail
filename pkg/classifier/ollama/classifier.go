package ollama

import (
    "fmt"
    "log"

    geemail "github.com/sverdejot/geemail/pkg"
    client "github.com/xyproto/ollamaclient/v2"
)

type classifier struct {
    c *client.Config
}

func New() (*classifier, error) {
    c := client.New("geemail")

    if err := c.PullIfNeeded(); err != nil {
        return nil, fmt.Errorf("cannot create ollama client: %v", err)
    }

    return &classifier{c}, nil
}

func (c *classifier) Classify(ct geemail.Content) string {
    out, err := c.c.GetOutput(ct.String())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(out)
    return out
}
