package httpclient

import (
    "encoding/json"
    "net/http"
    "time"
)

const DefaultTimeout = 10 * time.Second

type Client struct {
    client *http.Client
}

func New() *Client {
    return &Client{
        client: &http.Client{
            Timeout: DefaultTimeout,
        },
    }
}

func (c *Client) GetJSON(url string, target interface{}) error {
    response, err := c.client.Get(url)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    return json.NewDecoder(response.Body).Decode(target)
}
