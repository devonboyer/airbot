package airbot

import (
	"bytes"
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type Storage struct {
	client *storage.Client
}

func NewStorage(ctx context.Context) (*Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) Get(ctx context.Context, bucket, object string) ([]byte, error) {
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

func (s *Storage) Close() error {
	return s.client.Close()
}
