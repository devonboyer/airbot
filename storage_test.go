package airbot

import (
	"context"
	"os"
	"path/filepath"
	"testing"
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
	storage, err := NewStorage(ctx)
	if err != nil {
		t.Error("Error occurred", err)
	}
	defer storage.Close()

	data, err := storage.Get(ctx, storageBucketName, "secrets.encrypted")
	if err != nil {
		t.Error("Error occurred", err)
	}
	if len(data) == 0 {
		t.Error("No data")
	}
}
