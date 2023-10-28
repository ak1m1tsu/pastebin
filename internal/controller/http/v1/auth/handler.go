package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/log"
)

type handler struct {
	uc usecase.Auth
	l  *log.Logger
}

func New(mux chi.Router, uc usecase.Auth, l *log.Logger) {
	h := &handler{
		uc: uc,
		l:  l,
	}

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/", h.HandleRegistryUser)
		r.Post("/token", h.HandleGetToken)
	})
}

// HandleGetToken godoc
//
//	@summary	Получения авторизационных данных
//	@tags		auth
//	@produce	json
//
//	@accept		json
//
//	@param		code	body		entity.CreateTokenRequest	true	"Уникальный код, сгенерированный OAuth2 приложением"
//	@success	200		{object}	any{message=string}
//	@failure	400		{object}	any{message=string}
//	@failure	500		{object}	any{message=string}
//	@router		/token [post]
func (h *handler) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	var input entity.CreateTokenRequest

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		h.l.Error("failed to parse input data", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.BadRequest(w, r)

		return
	}

	ctx := context.TODO()

	token, err := h.uc.Token(ctx, input)
	if err != nil {
		h.l.Error("failed to generate token", err, log.FF{{Key: "input", Value: input}})

		response.InternalServerError(w, r)

		return
	}

	response.OK(w, r, render.M{
		"token":   token,
		"message": "ok",
	})
}

// HandleRegistryUser godoc
//
//	@summary	Регистрация нового пользователя с помощью Github OAuth 2.0
//	@tags		auth
//	@produce	json
//	@accept		json
//	@param		code	body		entity.CreateTokenRequest	true	"Уникальный код, сгенерированный OAuth2 приложением"
//	@success	200		{object}	any{message=string,data=any{user=entity.UserResponse}}
func (h *handler) HandleRegistryUser(w http.ResponseWriter, r *http.Request) {
	var input entity.CreateTokenRequest

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		h.l.Error("failed to parse input data", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.BadRequest(w, r)

		return
	}

	ctx := context.TODO()

	user, err := h.uc.CreateUser(ctx, input)
	if err != nil {
		h.l.Error("failed to register user", err, log.FF{{Key: "input", Value: input}})

		response.InternalServerError(w, r)

		return
	}

	response.OK(w, r, render.M{
		"user": user,
	})
}
