package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type Storage struct {
	client     *minio.Client
	bucketName string
}

func NewStorage(client *minio.Client, bucketName string) *Storage {
	return &Storage{client: client, bucketName: bucketName}
}

func (s *Storage) Upload(
	ctx context.Context,
	objectKey string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	_, err := s.client.PutObject(
		ctx,
		s.bucketName,
		objectKey,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("upload object %q: %w", objectKey, err)
	}

	return nil
}

func (s *Storage) Delete(
	ctx context.Context,
	objectKey string,
) error {
	err := s.client.RemoveObject(
		ctx,
		s.bucketName,
		objectKey,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("remove object %q: %w", objectKey, err)
	}
	return nil
}
