package main

import (
	"context"
	"net/http"
	"os"

	"github.com/devonboyer/airbot"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
)

var configDir, projectID, locationID, keyRingID, cryptoKeyID, storageBucketName string

func init() {
	configDir = "config"
	projectID = os.Getenv("PROJECT_ID")
	locationID = os.Getenv("KMS_LOCATION_ID")
	keyRingID = os.Getenv("KMS_KEYRING_ID")
	cryptoKeyID = os.Getenv("KMS_CRYPTOKEY_ID")
	storageBucketName = os.Getenv("STORAGE_BUCKET_NAME")
}

func main() {
	logrus.Info("Starting airbot")

	// Get ciphertext
	ciphertext, err := airbot.GetCiphertext(configDir)
	if err != nil {
		logrus.WithError(err).Panic("Could not read ciphertext")
	}

	// Decrypt secrets
	ctx := context.Background()
	secrets, err := airbot.DecryptSecrets(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
	if err != nil {
		logrus.WithError(err).Panic("Could not decrypt secrets")
	}
	logrus.Info("Decrypted secrets")

	// Setup webhook
	http.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			if req.FormValue("hub.verify_token") == secrets.Messenger.VerifyToken {
				w.Write([]byte(req.FormValue("hub.challenge")))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	logrus.Info("Starting appengine server")

	appengine.Main()
}
