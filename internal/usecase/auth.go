package usecase

import (
	"context"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
)

var _ Auth = &AuthUseCase{}

type AuthUseCase struct {
	oauth OAuthWebAPI
}

func NewAuth(oauth OAuthWebAPI) *AuthUseCase {
	return &AuthUseCase{oauth: oauth}
}

// Login implements Auth.
func (uc *AuthUseCase) Login(ctx context.Context, code string) (*entity.User, error) {
	token, err := uc.oauth.GetToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from code %q: %w", code, err)
	}

	user, err := uc.oauth.GetUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}
