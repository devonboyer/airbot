package messenger

import "net/http"

func (c *Client) WebhookHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			// Handle verification request.
			if req.FormValue("hub.verify_token") == c.verifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		case "POST":
			// Handle events.

			// process this
			// send to receiver channel
		}
	}
}
