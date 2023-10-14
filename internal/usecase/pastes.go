package usecase

import (
	"context"

	"github.com/romankravchuk/pastebin/internal/entity"
)

type PastesUseCase struct {
	repo  PastesRepo
	cache PastesCache
}

var _ Pastes = (*PastesUseCase)(nil)

func NewPastes(r PastesRepo, c PastesCache) *PastesUseCase {
	return &PastesUseCase{
		repo:  r,
		cache: c,
	}
}

// Create implements Pastes.
func (uc *PastesUseCase) Create(ctx context.Context, p entity.Paste) error {
	err := uc.repo.Store(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements Pastes.
func (uc *PastesUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.DeletePaste(ctx, id)
}

// Get implements Pastes.
func (uc *PastesUseCase) Get(ctx context.Context, id string) (entity.Paste, error) {
	paste, ok, err := uc.cache.Get(ctx, id)
	if err != nil {
		return paste, err
	}

	if ok {
		return paste, nil
	}

	paste, err = uc.repo.GetPaste(ctx, id)
	if err != nil {
		return paste, err
	}

	return paste, err
}

// Update implements Pastes.
func (uc *PastesUseCase) Update(ctx context.Context, p *entity.Paste) error {
	err := uc.repo.UpdatePaste(ctx, p)
	if err != nil {
		return err
	}

	return nil
}
