package secrets

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
	ciphertext := getCiphertext(t, configDir)

	ctx := context.Background()
	_, err := Decrypt(ctx, projectID, locationID, keyRingID, cryptoKeyID, ciphertext)
	require.NoError(t, err)
}

func getCiphertext(t *testing.T, dir string) []byte {
	data, err := ioutil.ReadFile(filepath.Join(dir, "secrets.encrypted"))
	require.NoError(t, err)
	require.NotZero(t, len(data))
	return bytes.TrimSpace(data)
}
