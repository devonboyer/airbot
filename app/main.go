package main

import (
	"context"
	"net/http"
	"os"

	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/messenger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
)

var version, env, configDir, projectID, locationID, keyRingID, cryptoKeyID, storageBucketName string

func init() {
	env = os.Getenv("ENV")
	projectID = os.Getenv("PROJECT_ID")
	locationID = os.Getenv("KMS_LOCATION_ID")
	keyRingID = os.Getenv("KMS_KEYRING_ID")
	cryptoKeyID = os.Getenv("KMS_CRYPTOKEY_ID")
	storageBucketName = os.Getenv("STORAGE_BUCKET_NAME")
}

func main() {
	// Setup logger.
	if env == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logger := logrus.WithField("version", version)

	logger.Info("Starting airbot")

	// Get storage client.
	ctx := context.Background()
	storage, err := airbot.NewStorageClient(ctx)
	if err != nil {
		logger.WithError(err).Panic("Could not create storage client")
	}
	defer storage.Close()

	// Get ciphertext
	ciphertext, err := storage.Get(ctx, storageBucketName, "secrets.encrypted")
	if err != nil {
		logrus.WithError(err).Panic("Could not read ciphertext")
	}
	logger.Info("Retrieved ciphertext")

	// Decrypt secrets
	secrets, err := airbot.DecryptSecrets(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
	if err != nil {
		logrus.WithError(err).Panic("Could not decrypt secrets")
	}
	logger.Info("Decrypted secrets")

	setupRoutes(secrets)
	setupBot(secrets)

	logrus.Info("Starting appengine server")

	appengine.Main()
}

func setupRoutes(secrets *airbot.Secrets) {
	mc := messenger.New(
		secrets.Messenger.AccessToken,
		secrets.Messenger.VerifyToken,
		secrets.Messenger.AppSecret,
	)

	r := mux.NewRouter()
	r.HandleFunc("/webhook", mc.WebhookHandler())
	http.Handle("/", r)
}

func setupBot(secrets *airbot.Secrets) {
	bot := airbot.NewBot()
	bot.Handle("shows today", func(s string) string {
		// TODO: Go get shows from airtable and reduce to a string
		return ""
	})
}
