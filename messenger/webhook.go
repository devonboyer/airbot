package messenger

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

// EventHandler responds to a webhook event.
//
// An EventHandler must be safe for concurrent use by multiple
// goroutines.
type EventHandler interface {
	HandleEvent(*Event)
}

func (c *Client) WebhookHandler(evh EventHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			glog.Info("Received webhook verification request")

			// Handle verification request.
			if req.FormValue("hub.verify_token") == c.verifyToken {
				setContentType(w.Header(), "text/plain; charset=utf-8")
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
				glog.Error("Could not ready body")
				handleError(w, http.StatusInternalServerError)
				return
			}

			glog.Infof("Received webhook event, %v", string(body))

			// Verify event signature.
			if !c.skipVerify && !verifySignature(c.appSecret, body, req.Header.Get("X-Hub-Signature")[5:]) {
				glog.Error("Invalid request signature")
				handleError(w, http.StatusUnauthorized)
				return
			}

			var ev = &Event{}
			if err := json.Unmarshal(body, ev); err != nil {
				glog.Errorf("Could not unmarshal event, %v", err)
				handleError(w, http.StatusInternalServerError)
				return
			}

			// Check the webhook event is from a Page subscription
			switch ev.Object {
			case "page":
				evh.HandleEvent(ev)
				w.WriteHeader(http.StatusOK)
			default:
				handleError(w, http.StatusNotFound)
			}
		default:
			handleError(w, http.StatusMethodNotAllowed)
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
