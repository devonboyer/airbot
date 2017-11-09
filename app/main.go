package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/devonboyer/airbot/botengine"

	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/airtable"
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

	logrus.Info("Starting airbot")

	// Get storage client.
	ctx := context.Background()
	storage, err := airbot.NewStorageClient(ctx)
	if err != nil {
		logrus.WithError(err).Panic("Could not create storage client")
	}
	defer storage.Close()

	// Get ciphertext
	ciphertext, err := storage.Get(ctx, storageBucketName, "secrets.encrypted")
	if err != nil {
		logrus.WithError(err).Panic("Could not read ciphertext")
	}
	logrus.Info("Retrieved ciphertext")

	// Decrypt secrets
	secrets, err := airbot.DecryptSecrets(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
	if err != nil {
		logrus.WithError(err).Panic("Could not decrypt secrets")
	}
	logrus.Info("Decrypted secrets")

	hc := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Get messenger client
	messengerClient := messenger.New(
		secrets.Messenger.AccessToken,
		secrets.Messenger.VerifyToken,
		secrets.Messenger.AppSecret,
		messenger.WithLogger(logrus.StandardLogger()),
		messenger.WithHTTPClient(hc),
	)

	// Get airtable client
	airtableClient := airtable.New(
		secrets.Airtable.APIKey,
		airtable.WithHTTPClient(hc),
	)

	chatService := messenger.NewChatService(messengerClient)

	// Create and setup bot.
	bot := botengine.New()
	bot.ChatService = chatService

	installNotFoundHandler(bot)

	// Run bot.
	logrus.Info("Starting bot")

	go runBot(bot, airtableClient)

	setupRoutes(messengerClient, chatService)

	logrus.Info("Starting appengine server")

	// Run appengine server.
	appengine.Main()
}

func setupRoutes(client *messenger.Client, evh messenger.EventHandler) {
	http.HandleFunc("/webhook", client.WebhookHandler(evh))
}

func runBot(bot *botengine.Bot, client *airtable.Client) {
	// Setup shows handlers
	shows := airbot.NewShowsBase(client)
	bot.HandleFunc("shows today", shows.TodayHandler())
	bot.HandleFunc("shows tomorrow", shows.TomorrowHandler())

	// Run bot.
	errsChan := bot.Run()
	defer bot.Stop()

	for {
		select {
		case err := <-errsChan:
			logrus.WithError(err).Error("bot error")
		}
	}
}

func installNotFoundHandler(bot *botengine.Bot) {
	bot.NotFoundHandler = botengine.HandlerFunc(func(w io.Writer, msg *botengine.Message) {
		fmt.Fprintf(w, "I don't understand \"%s\".", msg.Body)
	})
}
