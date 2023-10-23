package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	c *minio.Client
}

func New(endpoint, accessKey, secretKey string) (*Minio, error) {
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &Minio{c: c}, nil
}

func (m *Minio) GetObject(ctx context.Context, bucket, object string) (*minio.Object, error) {
	obj, err := m.c.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object with id %q from minio bucket %q: %w", object, bucket, err)
	}

	return obj, nil
}

func (m *Minio) GetBucketObjects(ctx context.Context, bucket string) ([]*minio.Object, error) {
	objects := make([]*minio.Object, 0)

	for objInfo := range m.c.ListObjects(ctx, bucket, minio.ListObjectsOptions{WithMetadata: true}) {
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
	exists, errBucketExists := m.c.BucketExists(ctx, bucket)
	if errBucketExists != nil || !exists {
		err := m.c.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %q: %w", bucket, err)
		}
	}

	_, err := m.c.PutObject(ctx, bucket, object, reader, size,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return fmt.Errorf("failed to upload object %q to minio bucket %q: %w", object, bucket, err)
	}

	return nil
}

func (m *Minio) DeleteObject(ctx context.Context, bucket, object string) error {
	err := m.c.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object %q from minio bucket %q: %w", object, bucket, err)
	}

	return nil
}
