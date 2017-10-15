package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
)

func main() {
	logrus.Info("starting airbot")

	http.HandleFunc("/", handle)

	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello world!")
}
