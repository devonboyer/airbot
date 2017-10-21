package messenger

import "net/http"

func (c *Client) WebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			if req.FormValue("hub.verify_token") == c.verifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
