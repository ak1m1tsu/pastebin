package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase/mocks"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func newPastesUseCase(t *testing.T) (*PastesUseCase, *mocks.PastesRepo, *mocks.PastesCache) {
	t.Helper()

	var (
		repo  = mocks.NewPastesRepo(t)
		cache = mocks.NewPastesCache(t)
	)

	return NewPastes(repo, cache), repo, cache
}

func TestPastesUseCase_Create(t *testing.T) {
	t.Parallel()

	t.Run(("Create paste"), func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _ = newPastesUseCase(t)
			ctx         = context.Background()
			paste       = entity.Paste{
				Hash:  "test",
				Block: []byte("test"),
			}
		)

		repo.On("Store", ctx, paste).
			Once().
			Return(nil)

		err := uc.Create(ctx, paste)
		require.NoError(t, err)
	})

	t.Run("Get error", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _ = newPastesUseCase(t)
			ctx         = context.Background()
			paste       = entity.Paste{
				Hash:  "test",
				Block: []byte("test"),
			}
		)

		repo.On("Store", ctx, paste).
			Once().
			Return(errTest)

		createErr := uc.Create(ctx, paste)
		require.Error(t, createErr)
	})
}

func TestPastesUseCase_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Delete paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _ = newPastesUseCase(t)
			ctx         = context.Background()
			id          = "test"
		)

		repo.On("DeletePaste", ctx, id).
			Once().
			Return(nil)

		err := uc.Delete(ctx, id)
		require.NoError(t, err)
	})
}

func TestPastesUseCase_Get(t *testing.T) {
	t.Parallel()

	t.Run("Get cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, _, cache = newPastesUseCase(t)
			ctx          = context.Background()
			expPaste     = entity.Paste{
				Hash:  "test",
				Block: []byte("test"),
			}
		)

		cache.On("Get", ctx, expPaste.Hash).
			Once().
			Return(expPaste, true, nil)

		paste, err := uc.Get(ctx, expPaste.Hash)
		require.NoError(t, err)
		require.Equal(t, expPaste, paste)
	})

	t.Run("Get non-cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, cache = newPastesUseCase(t)
			ctx             = context.Background()
			expPaste        = entity.Paste{
				Hash:  "test",
				Block: []byte("test"),
			}
		)

		cache.On("Get", ctx, expPaste.Hash).
			Once().
			Return(entity.Paste{}, false, nil)
		repo.On("GetPaste", ctx, expPaste.Hash).
			Once().
			Return(expPaste, nil)

		paste, err := uc.Get(ctx, expPaste.Hash)
		require.NoError(t, err)
		require.Equal(t, expPaste, paste)
	})

	t.Run("Get error on cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, _, cache = newPastesUseCase(t)
			ctx          = context.Background()
			id           = "test"
		)

		cache.On("Get", ctx, id).
			Once().
			Return(entity.Paste{}, false, entity.ErrPasteNotFound)

		paste, err := uc.Get(ctx, id)
		require.Error(t, err)
		require.Equal(t, entity.Paste{}, paste)
	})

	t.Run("Get error on non-cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, cache = newPastesUseCase(t)
			ctx             = context.Background()
			id              = "test"
		)

		cache.On("Get", ctx, id).
			Once().
			Return(entity.Paste{}, false, nil)
		repo.On("GetPaste", ctx, id).
			Once().
			Return(entity.Paste{}, errTest)

		paste, err := uc.Get(ctx, id)
		require.Error(t, err)
		require.Equal(t, entity.Paste{}, paste)
	})
}

func TestPastesUseCase_Update(t *testing.T) {
	t.Parallel()

	t.Run("Update paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _ = newPastesUseCase(t)
			ctx         = context.Background()
			paste       = &entity.Paste{
				Hash:  "test",
				Block: []byte("test"),
			}
		)

		repo.On("UpdatePaste", ctx, paste).
			Once().
			Return(nil)

		err := uc.Update(ctx, paste)
		require.NoError(t, err)
	})
}
