package main

import (
	"context"
	"net/http"
	"os"

	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/cli"
	"github.com/devonboyer/airbot/messenger"
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

	if env == "development" {
		cli.Run(nil)
	} else {
		// Get messenger client
		client := messenger.New(
			secrets.Messenger.AccessToken,
			secrets.Messenger.VerifyToken,
			secrets.Messenger.AppSecret,
		)

		// Run bot
		source := airbot.NewMessengerSource(client)
		bot := airbot.NewBot(secrets, source)
		bot.Run()
		defer bot.Stop()
		defer source.Stop()

		setupRoutes(client)

		appengine.Main()
	}
}

func setupRoutes(client *messenger.Client) {
	http.HandleFunc("/webhook", client.WebhookHandler())
}
