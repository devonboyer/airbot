package messenger

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) WebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			// Handle verification request.
			if req.FormValue("hub.verify_token") == c.verifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		case "POST":
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
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var res = &ReceiveBody{}
			err = json.NewDecoder(req.Body).Decode(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO: What now?

			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

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
