package airbot

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	configDir         = "config"
	storageBucketName = "storage-rising-artifact-182801"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", filepath.Join(configDir, "service-account.json"))
}

func TestStorage(t *testing.T) {
	// Get a file from cloud storage.
	ctx := context.Background()
	storage, err := NewStorageClient(ctx)
	require.NoError(t, err)
	defer storage.Close()

	data, err := storage.Get(ctx, storageBucketName, "secrets.encrypted")
	require.NoError(t, err)
	require.NotZero(t, len(data))
}
