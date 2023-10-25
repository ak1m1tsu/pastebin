package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
)

type PastesUseCase struct {
	repo  PastesRepo
	objs  PastesBlobStorage
	cache PastesCache
}

var _ Pastes = (*PastesUseCase)(nil)

func NewPastes(r PastesRepo, o PastesBlobStorage, c PastesCache) *PastesUseCase {
	return &PastesUseCase{
		objs:  o,
		repo:  r,
		cache: c,
	}
}

// Create creates a new paste.
//
// Uploads paste text to obj storage and stores paste metadata to database.
func (uc *PastesUseCase) Create(ctx context.Context, p *entity.Paste) error {
	if err := uc.objs.Create(ctx, p); err != nil {
		return fmt.Errorf("PastesUseCase.Create: %w", err)
	}

	err := uc.repo.Create(ctx, p)
	if err != nil {
		return fmt.Errorf("PastesUseCase.Create: %w", err)
	}

	return nil
}

// Delete implements Pastes.
func (uc *PastesUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.objs.Delete(ctx, id); err != nil {
		return fmt.Errorf("PastesUseCase.Delete: %w", err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("PastesUseCase.Delete: %w", err)
	}

	return nil
}

// Get returns a paste by hash.
//
// First checks if the paste is in the cache. If not, it gets the paste from the database.
// Then it gets the paste text from the obj storage.
func (uc *PastesUseCase) Get(ctx context.Context, hash string) (*entity.Paste, error) {
	paste, ok, err := uc.cache.Get(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("PastesUseCase.Get: %w", err)
	}

	if !ok {
		paste, err = uc.repo.Get(ctx, hash)
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				return nil, ErrPasteNotFound
			}

			return nil, fmt.Errorf("PastesUseCase.Get: %w", err)
		}
	}

	paste.File, err = uc.objs.Get(ctx, paste.UserID.String, paste.Hash)
	if err != nil {
		return nil, fmt.Errorf("PastesUseCase.Get: %w", err)
	}

	return paste, nil
}

// Update implements Pastes.
func (uc *PastesUseCase) Update(ctx context.Context, p *entity.Paste) error {
	if err := uc.objs.Update(ctx, p); err != nil {
		return fmt.Errorf("PastesUseCase.Update: %w", err)
	}

	if err := uc.repo.Update(ctx, p); err != nil {
		return fmt.Errorf("PastesUseCase.Update: %w", err)
	}

	return nil
}
