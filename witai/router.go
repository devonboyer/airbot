package witai

import (
	"github.com/devonboyer/airbot/botengine"
)

// Router implements botengine.Router
type Router struct {
	Client    *Client
	Processor IntentProcessor
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Match(msg *botengine.Message) botengine.Handler {
	return nil
}
