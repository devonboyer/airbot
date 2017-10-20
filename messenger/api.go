package messenger

import (
	"fmt"
	"net/http"
)

const apiVersion = "2.6"

type Client struct {
	accessToken string
	verifyToken string
	appSecret   string
	basePath    string
	hc          *http.Client
}

func New(accessToken, verifyToken, appSecret string) *Client {
	return &Client{
		accessToken: accessToken,
		verifyToken: verifyToken,
		appSecret:   appSecret,
		basePath:    fmt.Sprintf("https://graph.facebook.com/v%s/me/messages", apiVersion),
		hc:          http.DefaultClient,
	}
}
