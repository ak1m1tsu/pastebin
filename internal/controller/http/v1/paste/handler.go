package paste

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
	"github.com/romankravchuk/pastebin/internal/converter"
	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/log"
	"github.com/romankravchuk/pastebin/pkg/validator"
)

type handler struct {
	l  *log.Logger
	uc usecase.Pastes

	tm time.Duration
}

func New(mux chi.Router, uc usecase.Pastes, l *log.Logger) {
	p := &handler{
		l:  l,
		uc: uc,
		tm: 10 * time.Second,
	}

	mux.Route("/pastes", func(r chi.Router) {
		r.Post("/", p.HandleCreatePaste)
		r.Route("/{hash}", func(r chi.Router) {
			r.Get("/", p.HandleGetPasteByHash)
			r.Delete("/", p.HandleDeletePaste)
			r.Post("/unlock", p.HandleUnlockPaste)
		})
	})
}

// HandleCreatePaste godoc
//
//	@summary	Создание нововой пасты
//	@tags		pastes
//	@accept		json
//	@produce	json
//	@param		paste	body		entity.CreatePasteBody	true	"Паста"
//	@success	200		{object}	any{message=string,data=any{paste=entity.PasteResponse,url=string}}
//	@failure	400		{object}	any{message=string}
//	@failure	422		{object}	any{message=string,errors=any{field=string,message=string}}
//	@failure	500		{object}	any{message=string}
//	@router		/pastes [post]
func (h *handler) HandleCreatePaste(w http.ResponseWriter, r *http.Request) {
	input := new(entity.CreatePasteBody)

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		h.l.Error("failed to parse input data", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.BadRequest(w, r)

		return
	}

	v, err := validator.New()
	if err != nil {
		h.l.Error("failed to create validator", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.InternalServerError(w, r)

		return
	}

	if !v.Valid(input) {
		errs := v.Errors()

		h.l.Info("failed to validate input data", log.FF{
			{Key: "input", Value: input},
			{Key: "errors", Value: errs},
		})

		response.UnprocessableEntity(w, r, errs)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.tm)
	defer cancel()

	e, err := converter.CreatePasteToEntity(input)
	if err != nil {
		h.l.Error("failed to convert input data to entity", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.InternalServerError(w, r)

		return
	}

	err = h.uc.Create(ctx, e)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return
		default:
			h.l.Error("failed to create paste", err, log.FF{
				{Key: "input", Value: input},
			})

			response.InternalServerError(w, r)
		}

		return
	}

	location := fmt.Sprintf("%s/%s", r.URL.String(), e.Hash)

	w.Header().Add("Location", location)
	response.OK(w, r, render.M{
		"message": "ok",
		"data": render.M{
			"location": location,
			"paste":    converter.ModelToResponse(e),
		},
	})
}

// HandleGetPasteByHash godoc
//
//	@summary		Получениие пасты.
//	@description	Получение посты по хешу.
//	@description	Если паста защищена паролем, то нужно обратиться к `/pastes/{hash}/unlock`, чтобы получить доступ к ней.
//	@tags			pastes
//	@produce		json
//	@param			hash	path		string	true	"Хеш пасты"
//	@success		200		{object}	any{message=string,data=any{paste=entity.PasteResponse}}
//	@failure		403		{object}	any{error=string}
//	@failure		404		{object}	any{error=string}
//	@failure		500		{object}	any{error=string}
//	@router			/pastes/{hash} [get]
func (h *handler) HandleGetPasteByHash(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	ctx, cancel := context.WithTimeout(r.Context(), h.tm)
	defer cancel()

	paste, err := h.uc.Get(ctx, hash)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
		case errors.Is(err, usecase.ErrPasteNotFound):
			h.l.Warn("unable to get paste by hash", log.FF{{Key: "Hash", Value: hash}})

			response.NotFound(w, r)
		default:
			h.l.Error("failed to get paste by hash", err, log.FF{{Key: "Hash", Value: hash}})

			response.InternalServerError(w, r)
		}

		return
	}

	if paste.Password.Hash != nil {
		h.l.Warn("the paste lock for public review", log.FF{{Key: "hash", Value: hash}})

		response.Forbidden(w, r)

		return
	}

	response.OK(w, r, render.M{
		"message": "ok",
		"data": render.M{
			"paste": converter.ModelToResponse(paste),
		},
	})
}

// HandleDeletePaste godoc
//
//	@summary	Удаление пасты по хешу
//	@tags		pastes
//	@accept		json
//	@produce	json
//	@param		hash	path		string	true	"Хеш пасты"
//	@success	200		{object}	any{message=string}
//	@failure	403		{object}	any{error=string}
//	@failure	404		{object}	any{error=string}
//	@failure	500		{object}	any{error=string}
//	@security	Bearer
//	@router		/pastes/{hash} [delete]
func (h *handler) HandleDeletePaste(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	ctx, cancel := context.WithTimeout(r.Context(), h.tm)
	defer cancel()

	err := h.uc.Delete(ctx, hash)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
		case errors.Is(err, usecase.ErrPasteNotFound):
			h.l.Warn("unable to delete paste by hash", log.FF{{Key: "Hash", Value: hash}})

			response.NotFound(w, r)
		case errors.Is(err, usecase.ErrNotPasteAuthor):
			h.l.Warn("unable to delete paste by hash", log.FF{{Key: "Hash", Value: hash}})

			response.Forbidden(w, r)
		default:
			h.l.Error("unable to delete paste by hash", err, log.FF{{Key: "Hash", Value: hash}})

			response.InternalServerError(w, r)
		}

		return
	}

	response.OK(w, r, render.M{
		"message": "ok",
	})
}

// HandleUnlockPaste godoc
//
//	@summary	Получение доступа к пасте с паролем.
//	@tags		pastes
//	@accept		json
//	@produce	json
//	@param		hash		path		string					true	"Хеш пасты"
//	@param		credentials	body		entity.UnlockPasteBody	true	"Пароль"
//	@success	200			{object}	any{message=string,data=any{paste=entity.PasteResponse}}
//	@failure	403			{object}	any{error=string}
//	@failure	404			{object}	any{error=string}
//	@failure	500			{object}	any{error=string}
//	@router		/pastes/{hash}/unlock [post]
func (h *handler) HandleUnlockPaste(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	var input entity.UnlockPasteBody

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		h.l.Error("failed to parse input data", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.BadRequest(w, r)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.tm)
	defer cancel()

	paste, err := h.uc.Get(ctx, hash)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
		case errors.Is(err, usecase.ErrPasteNotFound):
			h.l.Warn("unable to get paste by hash", log.FF{{Key: "Hash", Value: hash}})

			response.NotFound(w, r)
		default:
			h.l.Error("failed to get paste by hash", err, log.FF{{Key: "Hash", Value: hash}})

			response.InternalServerError(w, r)
		}

		return
	}

	if !paste.Password.Matches(input.Password) {
		h.l.Warn("failed to unlock paste: invalid password", log.FF{{Key: "hash", Value: hash}})

		response.Forbidden(w, r)

		return
	}

	response.OK(w, r, render.M{
		"message": "ok",
		"data": render.M{
			"paste": converter.ModelToResponse(paste),
		},
	})
}
