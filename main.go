package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting airbot")

	secrets := Secrets{}
	if err := readSecrets(&secrets); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func readSecrets(secrets *Secrets) error {
	secrets.Airtable.APIKey = os.Getenv("AIRTABLE_API_KEY")
	if secrets.Airtable.APIKey == "" {
		return errors.New("AIRTABLE_API_KEY must be set")
	}

	secrets.Airtable.BaseID = os.Getenv("AIRTABLE_BASE_ID")
	if secrets.Airtable.BaseID == "" {
		return errors.New("AIRTABLE_BASE_ID must be set")
	}

	secrets.Messenger.AccessToken = os.Getenv("MESSENGER_ACCESS_TOKEN")
	if secrets.Messenger.AccessToken == "" {
		return errors.New("MESSENGER_ACCESS_TOKEN must be set")
	}

	secrets.Messenger.VerifyToken = os.Getenv("MESSENGER_VERIFY_TOKEN")
	if secrets.Messenger.AccessToken == "" {
		return errors.New("MESSENGER_VERIFY_TOKEN must be set")
	}

	secrets.Messenger.AppSecret = os.Getenv("MESSENGER_APP_SECRET")
	if secrets.Messenger.AccessToken == "" {
		return errors.New("MESSENGER_APP_SECRET must be set")
	}

	return nil
}
