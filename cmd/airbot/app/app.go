package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/botengine"
	"github.com/devonboyer/airbot/messenger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

func NewAirbotCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "airbot",
		Short: "Bot that responds to simple commands",
		Run: func(cmd *cobra.Command, args []string) {
			// Setup logger.
			if env == "development" {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetFormatter(&log.JSONFormatter{})
			}

			log.Info("Starting airbot")

			// Get storage client.
			ctx := context.Background()
			storage, err := airbot.NewStorageClient(ctx)
			if err != nil {
				log.WithError(err).Panic("Could not create storage client")
			}
			defer storage.Close()

			// Get ciphertext
			ciphertext, err := storage.Get(ctx, storageBucketName, "secrets.encrypted")
			if err != nil {
				log.WithError(err).Panic("Could not read ciphertext")
			}
			log.Info("Retrieved ciphertext")

			// Decrypt secrets
			secrets, err := airbot.DecryptSecrets(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
			if err != nil {
				log.WithError(err).Panic("Could not decrypt secrets")
			}
			log.Info("Decrypted secrets")

			hc := &http.Client{Timeout: 30 * time.Second}

			// Get messenger client
			messengerClient := messenger.New(
				secrets.Messenger.AccessToken,
				secrets.Messenger.VerifyToken,
				secrets.Messenger.AppSecret,
				messenger.WithLogger(log.StandardLogger()),
				messenger.WithHTTPClient(hc),
			)

			// Get airtable client
			airtableClient := airtable.New(
				secrets.Airtable.APIKey,
				airtable.WithHTTPClient(hc),
			)

			chatService := messenger.NewChatService(messengerClient)

			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

			stopCh := make(chan struct{})

			go func() {
				select {
				case sig := <-signalCh:
					log.Infof("Got %s signal. Aborting...", sig)
					close(stopCh)
				}
			}()

			// Create and setup bot.
			bot := botengine.New()
			bot.ChatService = chatService

			// Run bot.
			go runBot(stopCh, bot, airtableClient)

			setupRoutes(messengerClient, chatService)

			log.Info("Starting appengine server")

			// Run appengine server.
			appengine.Main()
		},
	}
}

func setupRoutes(client *messenger.Client, evh messenger.EventHandler) {
	http.HandleFunc("/webhook", client.WebhookHandler(evh))
}

func runBot(stopCh chan struct{}, bot *botengine.Bot, client *airtable.Client) {
	// Setup shows handlers
	shows := airbot.NewShowsBase(client)
	bot.HandleFunc("shows today", shows.TodayHandler())
	bot.HandleFunc("shows tomorrow", shows.TomorrowHandler())

	log.Info("Starting bot")

	// Run bot.
	errsChan := bot.Run()
	defer bot.Stop()

	for {
		select {
		case err := <-errsChan:
			log.WithError(err).Error("bot error")
		case <-stopCh:
			log.Info("Shutting down")
			return
		}
	}
}
