package airbot

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

const (
	configDir   = "config"
	projectID   = "rising-artifact-182801"
	locationID  = "global"
	keyRingID   = "airbot"
	cryptoKeyID = "secrets"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", filepath.Join(configDir, "service-account.json"))
}

func TestDecryptSecrets(t *testing.T) {
	ciphertext, err := GetCiphertext(configDir)
	if err != nil {
		t.Error("Error occurred")
	}
	if len(ciphertext) == 0 {
		t.Error("No ciphertext")
	}

	ctx := context.Background()
	_, err = DecryptSecrets(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
	if err != nil {
		t.Error("Error occurred", err)
	}
}
