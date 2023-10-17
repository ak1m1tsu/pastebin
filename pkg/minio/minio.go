package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	c       *minio.Client
	timeout time.Duration
}

func New(endpoint, accessKey, secretKey string, timeout time.Duration) (*Minio, error) {
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &Minio{
		c:       c,
		timeout: timeout,
	}, nil
}

func (m *Minio) GetObject(ctx context.Context, bucket, object string) (*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	obj, err := m.c.GetObject(reqCtx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object with id %q from minio bucket %q: %w", object, bucket, err)
	}

	return obj, nil
}

func (m *Minio) GetBucketObjects(ctx context.Context, bucket string) ([]*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	objects := make([]*minio.Object, 0)

	for objInfo := range m.c.ListObjects(reqCtx, bucket, minio.ListObjectsOptions{WithMetadata: true}) {
		if objInfo.Err != nil {
			return nil, fmt.Errorf("failed to list objects from minio bucket %q: %w", bucket, objInfo.Err)
		}

		object, err := m.c.GetObject(ctx, bucket, objInfo.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get object with id %q from minio bucket %q: %w", objInfo.Key, bucket, err)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func (m *Minio) UploadObject(ctx context.Context, bucket, object string, size int64, reader io.Reader) error {
	reqCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	exists, errBucketExists := m.c.BucketExists(ctx, bucket)
	if errBucketExists != nil || !exists {
		err := m.c.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %q: %w", bucket, err)
		}
	}

	_, err := m.c.PutObject(reqCtx, bucket, object, reader, size,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return fmt.Errorf("failed to upload object %q to minio bucket %q: %w", object, bucket, err)
	}

	return nil
}

func (m *Minio) DeleteObject(ctx context.Context, bucket, object string) error {
	reqCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	err := m.c.RemoveObject(reqCtx, bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object %q from minio bucket %q: %w", object, bucket, err)
	}

	return nil
}
