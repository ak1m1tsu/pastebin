package usecase

import (
	"context"

	"github.com/romankravchuk/pastebin/internal/entity"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Pastes --output ./mocks --outpkg mocks
type Pastes interface {
	Create(context.Context, *entity.Paste) error
	Get(context.Context, string) (*entity.Paste, error)
	Delete(context.Context, string) error
	Update(context.Context, *entity.Paste) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesRepo --output ./mocks --outpkg mocks
type PastesRepo interface {
	Create(context.Context, *entity.Paste) error
	Get(context.Context, string) (*entity.Paste, error)
	Delete(context.Context, string) error
	Update(context.Context, *entity.Paste) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesBlobStorage --output ./mocks --outpkg mocks
type PastesBlobStorage interface {
	Create(ctx context.Context, p *entity.Paste) error
	Get(ctx context.Context, userID, hash string) (entity.File, error)
	Delete(ctx context.Context, hash string) error
	Update(ctx context.Context, p *entity.Paste) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesCache --output ./mocks --outpkg mocks
type PastesCache interface {
	Create(context.Context, *entity.Paste) error
	Get(context.Context, string) (*entity.Paste, bool, error)
	Delete(context.Context, string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name OAuthWebAPI --output ./mocks --outpkg mocks
type OAuthWebAPI interface {
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*entity.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Auth --output ./mocks --outpkg mocks
type Auth interface {
	Login(ctx context.Context, code string) (*entity.User, error)
}
