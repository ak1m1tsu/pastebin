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
		r.Get("/{hash}", p.HandleGetPasteByHash)
		r.Delete("/{hash}", p.HandleDeletePaste)
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

	e := converter.CreatePasteToEntity(input)

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

	location := fmt.Sprintf("%s/%s", r.URL.Path, e.Hash)

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
//	@summary	Получениие пасту по хешу
//	@tags		pastes
//	@accept		json
//	@produce	json
//	@param		hash	path		string	true	"Хеш пасты"
//	@success	200		{object}	any{message=string}
//	@failure	400		{object}	any{error=string}
//	@router		/pastes/{hash} [get]
func (h *handler) HandleGetPasteByHash(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"message": "ok",
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
//	@security	Bearer
//	@router		/pastes/{hash} [delete]
func (h *handler) HandleDeletePaste(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"message": "ok",
	})
}
