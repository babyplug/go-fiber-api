package apikey

import (
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/core/repo"
	"go-fiber-api/internal/core/storage/db"
	"sync"
)

var (
	_repo     repo.Repo[model.APIKey, model.APIKeyDTO]
	_repoOnce sync.Once
)

func ProvideRepository(db db.Client) repo.Repo[model.APIKey, model.APIKeyDTO] {
	_repoOnce.Do(func() {
		_repo = repo.NewRepository[model.APIKey, model.APIKeyDTO](db)
	})

	return _repo
}
