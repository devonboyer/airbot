package airbot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

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

func GetCiphertext(dir string) ([]byte, error) {
	ciphertext, err := ioutil.ReadFile(filepath.Join(dir, "secrets.encrypted"))
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(ciphertext), nil
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
		fmt.Println("request failed")
		return nil, err
	}

	plaintext, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		fmt.Println("decode failed")
		return nil, err
	}

	var secrets = &Secrets{}
	if err := json.Unmarshal(plaintext, secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}
