package validator

import (
	"errors"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator

	errors map[string]string
	mu     sync.RWMutex
}

func New() (*Validator, error) {
	var (
		v             = validator.New()
		enTranslator  = en.New()
		translator, _ = ut.New(enTranslator, enTranslator).GetTranslator("en")
	)

	err := ent.RegisterDefaultTranslations(v, translator)
	if err != nil {
		return nil, err
	}

	return &Validator{
		errors:     make(map[string]string),
		validator:  v,
		translator: translator,
	}, nil
}

func (v *Validator) Valid(input interface{}) bool {
	var (
		errs validator.ValidationErrors
		err  = v.validator.Struct(input)
	)

	if err != nil {
		if errors.As(err, &errs) {
			v.mu.Lock()
			defer v.mu.Unlock()

			for _, e := range errs {
				v.errors[e.Field()] = e.Translate(v.translator)
			}
		}

		return false
	}

	return true
}

func (v *Validator) Errors() map[string]string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.errors
}
