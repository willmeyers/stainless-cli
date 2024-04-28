package stainless

import (
	"fmt"
	"net/http"
)

type Client struct {
	// Authenication credentials string used to populate `cookies` headers in requests made to Stainless API
	credentials string

	// HTTP client
	httpClient *http.Client
}

func New(options ...func(*Client) error) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{},
	}
	for _, option := range options {
		err := option(client)
		if err != nil {
			optionErr := fmt.Errorf("could not create new client using option %T", option)
			return nil, optionErr
		}
	}

	return client, nil
}

func WithCredentials(credentials string) func(*Client) error {
	return func(client *Client) error {
		client.credentials = credentials
		return nil
	}
}
