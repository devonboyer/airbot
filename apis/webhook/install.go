package webhook

import (
	"net/http"

	"github.com/devonboyer/airbot/messenger"
)

func Install(appSecret, verifyToken string, eventsCh chan<- *messenger.Event) {
	http.Handle("/webhook", messenger.NewWebhookHandler(appSecret, verifyToken, eventsCh))
}
