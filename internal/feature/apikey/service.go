//go:generate mockgen -source=service.go -mock_names=Service=MockAPIKeyService -destination=../../mock/mock_apikey_service.go -package=mock
package apikey

import (
	"context"
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/core/repo"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	svc     *serviceImpl
	svcOnce sync.Once
)

type Service interface {
	Create(ctx context.Context, dto *model.APIKeyDTO) (string, error)
	FindAll(ctx context.Context) ([]model.APIKeyDTO, error)
	FindByID(ctx context.Context, id any) (model.APIKeyDTO, error)
	DeleteByID(ctx context.Context, id any) error
}

type serviceImpl struct {
	cfg  *config.Configuration
	repo repo.Repo[model.APIKey, model.APIKeyDTO]
}

func ProvideService(cfg *config.Configuration, repo repo.Repo[model.APIKey, model.APIKeyDTO]) Service {
	svcOnce.Do(func() {
		svc = &serviceImpl{cfg: cfg, repo: repo}
	})

	return svc
}

func ResetService() {
	svcOnce = sync.Once{}
}

func (s *serviceImpl) generateToken(name string, duration model.Duration) (string, error) {
	claims := jwt.MapClaims{
		"name": name,
	}

	// No "exp" claim, so the token doesn't expire
	if duration != model.DurationUnlimited {
		now := time.Now()
		d := time.Hour * 24

		switch duration {
		case model.DurationNinetyDays:
			now.Add(d * 90)
			break
		case model.DurationThirtyDays:
			now.Add(d * 30)
			break
		default:
			now.Add(d * 7)
		}

		claims["exp"] = now.Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.SecretKey))
}

func (s *serviceImpl) Create(ctx context.Context, dto *model.APIKeyDTO) (string, error) {
	var token string
	var err error
	if dto != nil {
		token, err = s.generateToken(dto.Name, dto.Duration)
		if err != nil {
			return token, err
		}

		dto.Token = token
	}

	if err = s.repo.Insert(ctx, dto); err != nil {
		return token, err
	}

	return token, nil
}

func (s *serviceImpl) FindAll(ctx context.Context) ([]model.APIKeyDTO, error) {
	data, err := s.repo.FindAll(ctx)
	if err != nil {
		return []model.APIKeyDTO{}, err
	}

	return data, nil
}

func (s *serviceImpl) FindByID(ctx context.Context, id any) (model.APIKeyDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *serviceImpl) DeleteByID(ctx context.Context, id any) error {
	return s.repo.DeleteById(ctx, id)
}
