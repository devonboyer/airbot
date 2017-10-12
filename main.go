package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting airbot")

	secrets := &Secrets{}

	// Read environment variables
	secrets.Airtable.APIKey = os.Getenv("AIRTABLE_API_KEY")
	if secrets.Airtable.APIKey == "" {
		logrus.Panic("AIRTABLE_API_KEY must be set")
	}

	secrets.Airtable.BaseID = os.Getenv("AIRTABLE_BASE_ID must be set")
	if secrets.Airtable.BaseID == "" {
		logrus.Panic("AIRTABLE_BASE_ID must be set")
	}

	secrets.Messenger.AccessToken = os.Getenv("MESSENGER_ACCESS_TOKEN must be set")
	if secrets.Messenger.AccessToken == "" {
		logrus.Panic("MESSENGER_ACCESS_TOKEN must be set")
	}

	secrets.Messenger.VerifyToken = os.Getenv("MESSENGER_VERIFY_TOKEN must be set")
	if secrets.Messenger.AccessToken == "" {
		logrus.Panic("MESSENGER_VERIFY_TOKEN must be set")
	}

	secrets.Messenger.AppSecret = os.Getenv("MESSENGER_APP_SECRET must be set")
	if secrets.Messenger.AccessToken == "" {
		logrus.Panic("MESSENGER_APP_SECRET must be set")
	}
}
