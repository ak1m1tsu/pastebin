package usecase

import (
	"context"
	"fmt"

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

// Create creates a new paste.
func (uc *PastesUseCase) Create(ctx context.Context, p *entity.Paste) error {
	err := uc.repo.Create(ctx, p)
	if err != nil {
		return fmt.Errorf("PastesUseCase.Create: %w", err)
	}

	return nil
}

// Delete implements Pastes.
func (uc *PastesUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// Get implements Pastes.
func (uc *PastesUseCase) Get(ctx context.Context, id string) (*entity.Paste, error) {
	paste, ok, err := uc.cache.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("PastesUseCase.Get: %w", err)
	}

	if ok {
		return paste, nil
	}

	paste, err = uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("PastesUseCase.Get: %w", err)
	}

	return paste, nil
}

// Update implements Pastes.
func (uc *PastesUseCase) Update(ctx context.Context, p *entity.Paste) error {
	err := uc.repo.Update(ctx, p)
	if err != nil {
		return fmt.Errorf("PastesUseCase.Update: %w", err)
	}

	return nil
}
