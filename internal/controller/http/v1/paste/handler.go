package paste

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/log"
	"github.com/romankravchuk/pastebin/pkg/validator"
)

type handler struct {
	l  *log.Logger
	uc usecase.Pastes
}

func New(mux chi.Router, uc usecase.Pastes, l *log.Logger) {
	p := &handler{
		l:  l,
		uc: uc,
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
//	@param		paste	body		CreatePasteBody	true	"Паста"
//	@success	200		{object}	any{message=string}
//	@router		/pastes [post]
func (p *handler) HandleCreatePaste(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text       string `json:"text" validate:"required"`
		ExpireTime string `json:"expireTime" validate:"required,oneof=30m 1h 1d 1w 1mth 6mth 2yr"`
		UserID     string `json:"userID" validate:"omitempty,uuid"`
		Password   string `json:"password" validate:"omitempty"`
		Title      string `json:"title" validate:"omitempty"`
		Format     string `json:"format" validate:"required,oneof=txt html markdown latex json xml yaml ini csv tsv url binary"`
	}

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		p.l.Error("failed to parse input data", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.BadRequest(w, r)

		return
	}

	v, err := validator.New()
	if err != nil {
		p.l.Error("failed to create validator", err,
			log.FF{
				{Key: "input", Value: input},
			})

		response.InternalServerError(w, r)

		return
	}

	if !v.Valid(input) {
		errs := v.Errors()

		p.l.Info("failed to validate input data", log.FF{
			{Key: "input", Value: input},
			{Key: "errors", Value: errs},
		})

		response.UnprocessableEntity(w, r, errs)

		return
	}

	response.OK(w, r, render.M{
		"message": "ok",
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
func (p *handler) HandleGetPasteByHash(w http.ResponseWriter, r *http.Request) {
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
func (p *handler) HandleDeletePaste(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"message": "ok",
	})
}
