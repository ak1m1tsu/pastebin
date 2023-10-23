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

func newPastesUseCase(t *testing.T) (*PastesUseCase, *mocks.PastesRepo, *mocks.PastesBlobStorage, *mocks.PastesCache) {
	t.Helper()

	var (
		repo  = mocks.NewPastesRepo(t)
		cache = mocks.NewPastesCache(t)
		blob  = mocks.NewPastesBlobStorage(t)
	)

	return NewPastes(repo, blob, cache), repo, blob, cache
}

func TestPastesUseCase_Create(t *testing.T) {
	t.Parallel()

	t.Run(("Create paste"), func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, blob, _ = newPastesUseCase(t)
			ctx               = context.Background()
			paste             = &entity.Paste{
				Hash: "test",
				File: []byte("test"),
			}
		)

		blob.On("Create", ctx, paste).
			Once().
			Return(nil)
		repo.On("Create", ctx, paste).
			Once().
			Return(nil)

		err := uc.Create(ctx, paste)
		require.NoError(t, err)
	})

	t.Run("Get error", func(t *testing.T) {
		t.Parallel()

		t.Run("Error on blob", func(t *testing.T) {
			t.Parallel()

			var (
				uc, _, blob, _ = newPastesUseCase(t)
				ctx            = context.Background()
				paste          = &entity.Paste{
					Hash: "test",
					File: []byte("test"),
				}
			)

			blob.On("Create", ctx, paste).
				Once().
				Return(errTest)

			err := uc.Create(ctx, paste)
			require.Error(t, err)
		})

		t.Run("Error on repo", func(t *testing.T) {
			t.Parallel()

			var (
				uc, repo, blob, _ = newPastesUseCase(t)
				ctx               = context.Background()
				paste             = &entity.Paste{
					Hash: "test",
					File: []byte("test"),
				}
			)

			blob.On("Create", ctx, paste).
				Once().
				Return(nil)
			repo.On("Create", ctx, paste).
				Once().
				Return(errTest)

			createErr := uc.Create(ctx, paste)
			require.Error(t, createErr)
		})
	})
}

func TestPastesUseCase_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Delete paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _, _ = newPastesUseCase(t)
			ctx            = context.Background()
			id             = "test"
		)

		repo.On("Delete", ctx, id).
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
			uc, _, _, cache = newPastesUseCase(t)
			ctx             = context.Background()
			expPaste        = &entity.Paste{
				Hash: "test",
				File: []byte("test"),
			}
		)

		cache.On("Get", ctx, expPaste.Hash).
			Once().
			Return(expPaste, true, nil)

		paste, err := uc.Get(ctx, expPaste.Hash)
		require.NoError(t, err)
		require.NotNil(t, paste)
	})

	t.Run("Get non-cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _, cache = newPastesUseCase(t)
			ctx                = context.Background()
			expPaste           = &entity.Paste{
				Hash: "test",
				File: []byte("test"),
			}
		)

		cache.On("Get", ctx, expPaste.Hash).
			Once().
			Return(nil, false, nil)
		repo.On("Get", ctx, expPaste.Hash).
			Once().
			Return(expPaste, nil)

		paste, err := uc.Get(ctx, expPaste.Hash)
		require.NoError(t, err)
		require.NotNil(t, paste)
	})

	t.Run("Get error on cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, _, _, cache = newPastesUseCase(t)
			ctx             = context.Background()
			id              = "test"
		)

		cache.On("Get", ctx, id).
			Once().
			Return(nil, false, entity.ErrPasteNotFound)

		paste, err := uc.Get(ctx, id)
		require.Error(t, err)
		require.Nil(t, paste)
	})

	t.Run("Get error on non-cached paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _, cache = newPastesUseCase(t)
			ctx                = context.Background()
			id                 = "test"
		)

		cache.On("Get", ctx, id).
			Once().
			Return(nil, false, nil)
		repo.On("Get", ctx, id).
			Once().
			Return(nil, errTest)

		paste, err := uc.Get(ctx, id)
		require.Error(t, err)
		require.Nil(t, paste)
	})
}

func TestPastesUseCase_Update(t *testing.T) {
	t.Parallel()

	t.Run("Update paste", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _, _ = newPastesUseCase(t)
			ctx            = context.Background()
			paste          = &entity.Paste{
				Hash: "test",
				File: []byte("test"),
			}
		)

		repo.On("Update", ctx, paste).
			Once().
			Return(nil)

		err := uc.Update(ctx, paste)
		require.NoError(t, err)
	})

	t.Run("Get error on update", func(t *testing.T) {
		t.Parallel()

		var (
			uc, repo, _, _ = newPastesUseCase(t)
			ctx            = context.Background()
			paste          = &entity.Paste{
				Hash: "test",
				File: []byte("test"),
			}
		)

		repo.On("Update", ctx, paste).
			Once().
			Return(errTest)

		err := uc.Update(ctx, paste)
		require.Error(t, err)
	})
}
