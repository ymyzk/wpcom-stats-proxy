package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Client struct {
	config     *Configuration
	HTTPClient *http.Client
}

func NewClient(conf *Configuration) *Client {
	config := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://public-api.wordpress.com/oauth2/authorize",
			TokenURL: "https://public-api.wordpress.com/oauth2/token",
		},
	}
	token := oauth2.Token{AccessToken: conf.Token}
	return &Client{
		config:     conf,
		HTTPClient: config.Client(oauth2.NoContext, &token),
	}
}

func (c *Client) getPostStats(stats *wordpressPostStats, id int) error {
	url := fmt.Sprintf("https://public-api.wordpress.com/rest/v1.1/sites/%s/stats/post/%d", c.config.Site, id)
	res, _ := c.HTTPClient.Get(url)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// TODO: Better error handling
		return errors.New(fmt.Sprintf("request failed with status code: %d", res.StatusCode))
	}

	return json.NewDecoder(res.Body).Decode(&stats)
}
