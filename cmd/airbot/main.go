package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/devonboyer/airbot/cmd/airbot/app"
)

func main() {
	if err := app.NewAirbotCommand().Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}
