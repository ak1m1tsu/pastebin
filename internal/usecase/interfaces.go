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
	Delete(ctx context.Context, userID, hash string) error
	Update(ctx context.Context, p *entity.Paste) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PastesCache --output ./mocks --outpkg mocks
type PastesCache interface {
	Create(context.Context, *entity.Paste) error
	Get(context.Context, string) (*entity.Paste, bool, error)
	Delete(context.Context, string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name AuthWebAPI --output ./mocks --outpkg mocks
type AuthWebAPI interface {
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*entity.APIUser, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Auth --output ./mocks --outpkg mocks
type Auth interface {
	Token(ctx context.Context, req entity.CreateTokenRequest) (*entity.TokenCredentails, error)
	CreateUser(ctx context.Context, req entity.CreateTokenRequest) (*entity.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name UsersRepo --output ./mocks --outpkg mocks
type UsersRepo interface {
	Create(ctx context.Context, u *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
