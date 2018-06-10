package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/apis/shows"
	"github.com/devonboyer/airbot/apis/webhook"
	"github.com/devonboyer/airbot/botengine"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/appengine"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetDefault("env", "development")

	viper.SetDefault("project_id", "rising-artifact-182801")
	viper.SetDefault("kms_location_id", "global")
	viper.SetDefault("kms_keyring_id", "airbot")
	viper.SetDefault("kms_cryptokey_id", "secrets")
	viper.SetDefault("storage_bucket_name", "storage-rising-artifact-182801")

	viper.BindEnv("env")
}

const programName = "airbot"

var rootCmd = &cobra.Command{
	Use:   programName,
	Short: "Bot that responds to simple commands",
	Run: func(cmd *cobra.Command, args []string) {
		initLogs()

		log.Info("Starting airbot")

		// Get storage client.
		ctx := context.Background()
		storage, err := airbot.NewStorageClient(ctx)
		if err != nil {
			log.WithError(err).Panic("Could not create storage client")
		}
		defer storage.Close()

		// Get ciphertext
		ciphertext, err := storage.Get(ctx, viper.GetString("storage_bucket_name"), "secrets.encrypted")
		if err != nil {
			log.WithError(err).Panic("Could not read ciphertext")
		}
		log.Info("Retrieved ciphertext")

		// Decrypt secrets
		secrets, err := airbot.DecryptSecrets(ctx,
			viper.GetString("project_id"),
			viper.GetString("kms_location_id"),
			viper.GetString("kms_keyring_id"),
			viper.GetString("kms_cryptokey_id"),
			ciphertext)
		if err != nil {
			log.WithError(err).Panic("Could not decrypt secrets")
		}
		log.Info("Decrypted secrets")

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

		var (
			httpClient      = &http.Client{Timeout: 30 * time.Second}
			airtableClient  = secrets.NewAirtableClient(httpClient)
			messengerClient = secrets.NewMessengerClient(httpClient)
		)

		svc := airbot.NewMessengerService(messengerClient)
		svc.Run()

		// Create and setup bot.
		bot := botengine.New()
		bot.ChatService = svc

		// Install apis
		shows.Install(bot, airtableClient)
		webhook.Install(secrets.Messenger.AppSecret, secrets.Messenger.VerifyToken, svc.Events())

		log.Info("Starting bot")
		go runBot(stopCh, bot)

		log.Info("Starting appengine server")
		appengine.Main()
	},
}

func runBot(stopCh chan struct{}, bot *botengine.Bot) {
	errsCh := bot.Run()
	defer bot.Stop()

	for {
		select {
		case err := <-errsCh:
			log.WithError(err).Error("bot error")
		case <-stopCh:
			log.Info("Shutting down")
			return
		}
	}
}

func initLogs() {
	if viper.GetString("env") == "development" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}
