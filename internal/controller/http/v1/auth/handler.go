package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
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
		r.Get("/token", h.HandleGetToken)
	})
}

// HandleGetToken godoc
//
//	@summary	Получения авторизационных данных
//	@tags		auth
//	@produce	json
//	@param		code	query		string	true	"Уникальный код, сгенерированный OAuth2 приложением"
//	@success	200		{object}	any{message=string}
//	@router		/token [get]
func (h *handler) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = context.TODO()
		code = r.URL.Query().Get("code")
	)

	user, err := h.uc.Login(ctx, code)
	if err != nil {
		h.l.Error("error", err, nil)
	}

	response.OK(w, r, render.M{
		"credentials": user,
		"message":     "ok",
	})
}
