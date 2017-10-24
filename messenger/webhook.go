package messenger

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EventHandler interface {
	// HandleEvent may be called from multiple goroutines. Note that no effort is made to buffer events.
	HandleEvent(*WebhookEvent)
}

func (c *Client) WebhookHandler(handler EventHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			c.logger.Printf("messenger: received webhook verification request")

			// Handle verification request.
			if req.FormValue("hub.verify_token") == c.verifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		case "POST":
			c.logger.Printf("messenger: received webhook event")

			// Handle event.
			defer req.Body.Close()

			// Read body
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Validate event.
			if !verifySignature(c.appSecret, body, req.Header.Get("X-Hub-Signature")[5:]) {
				c.logger.Printf("messenger: invalid request signature")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var ev = &WebhookEvent{}
			err = json.NewDecoder(req.Body).Decode(ev)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// FIXME: Probably a bad idea to call handler synchronously
			handler.HandleEvent(ev)

			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// https://developers.facebook.com/docs/messenger-platform/webhook#security
func verifySignature(appSecret string, bytes []byte, expectedSignature string) bool {
	if expectedSignature == "" {
		return false
	}
	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write(bytes)
	if fmt.Sprintf("%x", mac.Sum(nil)) != expectedSignature {
		return false
	}
	return true
}
