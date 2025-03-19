package user

import (
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/core/repo"
	"go-fiber-api/internal/core/storage/db"
	"sync"
)

var (
	_repo     repo.Repo[model.User, model.UserDTO]
	_repoOnce sync.Once
)

func ProvideRepository(db db.Client) repo.Repo[model.User, model.UserDTO] {
	_repoOnce.Do(func() {
		_repo = repo.NewRepository[model.User, model.UserDTO](db)
	})

	return _repo
}
