package go_receive

import (
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}
