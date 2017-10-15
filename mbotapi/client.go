package mbotapi

import "net/http"

type Client struct {
	VerifyToken string
}

func New(verifyToken string) *Client {
	return &Client{
		VerifyToken: verifyToken,
	}
}

func (c *Client) SetWebhook(pattern string) *http.ServeMux {
	mux := http.NewServeMux()
	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			if req.FormValue("hub.verify_token") == c.VerifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
	return mux
}
