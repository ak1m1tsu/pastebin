package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var _ usecase.AuthWebAPI = &GithubAPI{}

type GithubAPI struct {
	cfg          *oauth2.Config
	infoEndpoint string
}

func NewGithubAPI(clientID, clientSecret string) *GithubAPI {
	return &GithubAPI{
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"user:read"},
			Endpoint:     github.Endpoint,
		},
		infoEndpoint: "https://api.github.com/user",
	}
}

func (api *GithubAPI) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := api.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("GithubAPI.GetToken: %w", err)
	}

	return token, nil
}

func (api *GithubAPI) GetUserInfo(ctx context.Context, token *oauth2.Token) (*entity.APIUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api.infoEndpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("GithubAPI.GetUserInfo: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GithubAPI.GetUserInfo: %w", err)
	}
	defer resp.Body.Close()

	var user *entity.APIUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("GithubAPI.GetUserInfo: %w", err)
	}

	return user, nil
}
