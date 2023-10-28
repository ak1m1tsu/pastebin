package usecase

import (
	"context"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
)

var _ Auth = &AuthUseCase{}

type AuthUseCase struct {
	users UsersRepo
	oauth AuthWebAPI
}

func NewAuth(users UsersRepo, oauth AuthWebAPI) *AuthUseCase {
	return &AuthUseCase{
		oauth: oauth,
		users: users,
	}
}

// Token godoc.
func (uc *AuthUseCase) Token(ctx context.Context, req entity.CreateTokenRequest) (*entity.TokenCredentails, error) {
	token, err := uc.oauth.GetToken(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from code %q: %w", req.Code, err)
	}

	apiUser, err := uc.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := uc.users.GetByEmail(ctx, apiUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email %q: %w", apiUser.Email, err)
	}

	return &entity.TokenCredentails{
		UserID:      user.ID,
		Email:       user.Email,
		AccessToken: token.AccessToken,
		Type:        token.TokenType,
		ExpireAt:    token.Expiry,
	}, nil
}

// CreateUser godoc.
func (uc *AuthUseCase) CreateUser(ctx context.Context, req entity.CreateTokenRequest) (*entity.User, error) {
	token, err := uc.oauth.GetToken(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from code %q: %w", req.Code, err)
	}

	apiUser, err := uc.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user := &entity.User{
		Email:    apiUser.Email,
		Username: apiUser.Username,
		Avatar:   apiUser.Avatar,
	}
	if err = user.AccessToken.GenerateFrom([]byte(token.AccessToken)); err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	err = uc.users.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
