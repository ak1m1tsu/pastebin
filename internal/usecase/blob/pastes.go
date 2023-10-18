package blob

import (
	"context"

	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/minio"
)

var _ usecase.PastesBlobStorage = &PastesBlobStorage{}

type PastesBlobStorage struct {
	*minio.Minio
}

// Create implements usecase.PastesBlobStorage.
func (*PastesBlobStorage) Create(ctx context.Context, bucket string, id string, data []byte) error {
	panic("unimplemented")
}

// Delete implements usecase.PastesBlobStorage.
func (*PastesBlobStorage) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// Get implements usecase.PastesBlobStorage.
func (*PastesBlobStorage) Get(ctx context.Context, id string) ([]byte, error) {
	panic("unimplemented")
}

// Update implements usecase.PastesBlobStorage.
func (*PastesBlobStorage) Update(ctx context.Context, id string, data []byte) error {
	panic("unimplemented")
}

func New(m *minio.Minio) *PastesBlobStorage {
	return &PastesBlobStorage{m}
}
