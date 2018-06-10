package airbot

import (
	"bytes"
	"context"
	"io"

	gcs "cloud.google.com/go/storage"
)

type StorageClient struct {
	client *gcs.Client
}

func NewStorageClient(ctx context.Context) (*StorageClient, error) {
	client, err := gcs.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &StorageClient{client}, nil
}

func (s *StorageClient) Get(ctx context.Context, bucket, object string) ([]byte, error) {
	reader, err := s.client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

func (s *StorageClient) Close() error {
	return s.client.Close()
}
