package airbot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/messenger"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type Secrets struct {
	Airtable struct {
		APIKey string `json:"api_key"`
	} `json:"airtable"`
	Messenger struct {
		AccessToken string `json:"access_token"`
		VerifyToken string `json:"verify_token"`
		AppSecret   string `json:"app_secret"`
	} `json:"messenger"`
}

func DecryptSecrets(ctx context.Context, projectID, locationID, keyRingID, cryptoKeyID string, ciphertext []byte) (*Secrets, error) {
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, locationID, keyRingID, cryptoKeyID)

	req := &cloudkms.DecryptRequest{
		Ciphertext: string(ciphertext),
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parentName, req).Do()
	if err != nil {
		return nil, errors.Wrap(err, "decrypt request failed")
	}

	plaintext, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}

	var secrets = &Secrets{}
	if err := json.Unmarshal(plaintext, secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *Secrets) NewMessengerClient(httpClient *http.Client) *messenger.Client {
	return messenger.New(
		s.Messenger.AccessToken,
		messenger.WithHTTPClient(httpClient),
	)
}

func (s *Secrets) NewAirtableClient(httpClient *http.Client) *airtable.Client {
	return airtable.New(
		s.Airtable.APIKey,
		airtable.WithHTTPClient(httpClient),
	)
}

func MustReadSecrets(dir string) *Secrets {
	file, err := os.Open(filepath.Join(dir, "secrets.json"))
	if err != nil {
		panic(err)
	}
	var secrets = &Secrets{}
	if err := json.NewDecoder(file).Decode(secrets); err != nil {
		panic(err)
	}
	return secrets
}
