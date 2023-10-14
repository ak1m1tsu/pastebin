package usecase

import (
	"context"

	"github.com/romankravchuk/pastebin/internal/entity"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Pastes --output ./mocks --outpkg mocks
type Pastes interface {
	Create(context.Context, entity.Paste) error
	Get(context.Context, string) (entity.Paste, error)
	Delete(context.Context, string) error
	Update(context.Context, *entity.Paste) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesRepo --output ./mocks --outpkg mocks
type PastesRepo interface {
	Store(context.Context, entity.Paste) error
	GetPaste(context.Context, string) (entity.Paste, error)
	DeletePaste(context.Context, string) error
	UpdatePaste(context.Context, *entity.Paste) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesBlobStorage --output ./mocks --outpkg mocks
type PastesBlobStorage interface {
	Create(context.Context, string, []byte) error
	Get(context.Context, string) ([]byte, error)
	Delete(context.Context, string) error
	Update(context.Context, string, []byte) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesCache --output ./mocks --outpkg mocks
type PastesCache interface {
	Create(context.Context, entity.Paste) error
	Get(context.Context, string) (entity.Paste, bool, error)
	Delete(context.Context, string) error
}
