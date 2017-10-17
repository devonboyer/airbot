package secrets

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	configDir   = "../config"
	projectID   = "rising-artifact-182801"
	locationID  = "global"
	keyRingID   = "airbot"
	cryptoKeyID = "secrets"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", filepath.Join(configDir, "service-account.json"))
}

func TestDecrypt(t *testing.T) {
	ciphertext, err := getCiphertext(configDir)
	if err != nil {
		t.Error("Error occurred", err)
	}

	ctx := context.Background()
	_, err = Decrypt(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
	if err != nil {
		t.Error("Error occurred", err)
	}
}

func getCiphertext(dir string) ([]byte, error) {
	ciphertext, err := ioutil.ReadFile(filepath.Join(dir, "secrets.encrypted"))
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(ciphertext), nil
}
