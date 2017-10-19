package messenger

import "net/http"

func WebhookHandler(verifyToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			if req.FormValue("hub.verify_token") == verifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
