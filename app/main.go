package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
)

func main() {
	logrus.Info("starting airbot")

	appengine.Main()
}
