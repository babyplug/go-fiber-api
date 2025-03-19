//go:generate mockgen -source=service.go -mock_names=Service=MockGameService -destination=../../mock/mock_game_service.go  -package=mock

package user

import (
	"context"
	"errors"
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/core/repo"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	sOnce sync.Once
	s     *serviceImpl
)

type Service interface {
	FindAll(ctx context.Context) ([]model.UserDTO, error)
	FindByID(ctx context.Context, id string) (model.UserDTO, error)
	Create(context.Context, *model.UserDTO) error
	Update(context.Context, *model.UserDTO) error
	DeleteByID(context.Context, string) error
}

type serviceImpl struct {
	repo repo.Repo[model.User, model.UserDTO]
}

func ProvideService(repo repo.Repo[model.User, model.UserDTO]) Service {
	sOnce.Do(func() {
		s = &serviceImpl{
			repo: repo,
		}
	})
	return s
}

func ResetProvideService() {
	sOnce = sync.Once{}
}

func (s *serviceImpl) Create(ctx context.Context, dto *model.UserDTO) error {
	if dto == nil {
		return errors.New("dto can not be nil")
	}

	err := s.repo.Insert(ctx, dto)
	if err != nil {
		logrus.Errorf("game service create: %v", err)
		return err
	}

	return nil
}

func (s *serviceImpl) FindAll(ctx context.Context) ([]model.UserDTO, error) {
	// Define the scan input
	data, err := s.repo.FindAll(ctx)
	if err != nil {
		logrus.Errorf("game service get: %v", err)
		return nil, err
	}

	return data, nil
}

func (s *serviceImpl) FindByID(ctx context.Context, id string) (model.UserDTO, error) {
	res, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logrus.Errorf("game service get by: %v", err)
		return model.UserDTO{}, err
	}

	return res, nil
}

func (s *serviceImpl) Update(ctx context.Context, dto *model.UserDTO) error {
	err := s.repo.Update(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) DeleteByID(ctx context.Context, id string) error {
	err := s.repo.DeleteById(ctx, id)
	if err != nil {
		logrus.Errorf("game service delete: %v", err)
		return err
	}

	return nil
}
