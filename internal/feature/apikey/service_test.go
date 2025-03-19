package apikey_test

import (
	"context"
	"errors"
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/core/repo"
	"go-fiber-api/internal/feature/apikey"
	"go-fiber-api/internal/mock"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Apikey_serviceImpl_Create(t *testing.T) {
	type dependency struct {
		cfg  *config.Configuration
		repo func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO]
	}

	cfg := &config.Configuration{
		SecretKey: "TEST_SECRET_KEY",
	}

	tests := []struct {
		name string
		dependency
		dto            *model.APIKeyDTO
		expectedExpNil bool
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			name: "when_insert_error_should_get_error",
			dependency: dependency{
				cfg: cfg,
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					m := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)

					m.EXPECT().
						Insert(gomock.Any(), gomock.Any()).
						Return(errors.New("mock insert error"))

					return m
				},
			},
			expectedErr:    true,
			expectedErrMsg: "mock insert error",
		},
		{
			name: "when_generate_token_with_unlimited_expiration_should_get_expiration_date_nil",
			dependency: dependency{
				cfg: cfg,
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					m := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)

					m.EXPECT().
						Insert(gomock.Any(), gomock.Any()).
						Return(nil)

					return m
				},
			},
			dto: &model.APIKeyDTO{
				Name:     "Test",
				Duration: model.DurationUnlimited,
			},
			expectedExpNil: true,
		},
		{
			name: "when_generate_token_with_7day_should_get_expiration_date_equal_to_now_added_by_dto_duration",
			dependency: dependency{
				cfg: cfg,
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					m := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)

					m.EXPECT().
						Insert(gomock.Any(), gomock.Any()).
						Return(nil)

					return m
				},
			},
			dto: &model.APIKeyDTO{
				Name:     "Test",
				Duration: model.DurationSevenDays,
			},
		},
		{
			name: "when_generate_token_with_30day_should_get_expiration_date_equal_to_now_added_by_dto_duration",
			dependency: dependency{
				cfg: cfg,
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					m := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)

					m.EXPECT().
						Insert(gomock.Any(), gomock.Any()).
						Return(nil)

					return m
				},
			},
			dto: &model.APIKeyDTO{
				Name:     "Test",
				Duration: model.DurationThirtyDays,
			},
		},
		{
			name: "when_generate_token_with_30day_should_get_expiration_date_equal_to_now_added_by_dto_duration",
			dependency: dependency{
				cfg: cfg,
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					m := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)

					m.EXPECT().
						Insert(gomock.Any(), gomock.Any()).
						Return(nil)

					return m
				},
			},
			dto: &model.APIKeyDTO{
				Name:     "Test",
				Duration: model.DurationNinetyDays,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := apikey.ProvideService(test.cfg, test.dependency.repo(ctrl))
			defer apikey.ResetService()

			ctx := context.TODO()
			tokenStr, err := s.Create(ctx, test.dto)
			if test.expectedErr && assert.Error(t, err) {
				assert.Equal(t, test.expectedErrMsg, err.Error())
				return
			}

			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return cfg.SecretKey, nil
			})

			if assert.NoError(t, err) {
				assert.NotEmpty(t, tokenStr)
				d, _ := token.Claims.GetExpirationTime()
				if test.expectedExpNil {
					assert.Nil(t, d)
					return
				} else {
					now := time.Now()
					day := time.Hour * 24
					switch test.dto.Duration {
					case model.DurationSevenDays:
						now.Add(day * 7)
						break
					case model.DurationThirtyDays:
						now.Add(day * 30)
						break
					case model.DurationNinetyDays:
						now.Add(day * 90)
						break
					}
					assert.Equal(t, now.Day(), d.Day())
				}
			}

		})
	}

}

func Test_Apikey_serviceImpl_Query(t *testing.T) {
	type dependency struct {
		repo func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO]
	}

	ctx := context.Background()
	token1 := uuid.New().String()
	token2 := uuid.New().String()

	tests := []struct {
		name string
		dependency
		expectedErr    bool
		expectedErrMsg string
		expected       []model.APIKeyDTO
	}{
		{
			name: "when_repo_error_should_return_nil",
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
					repo.EXPECT().
						FindAll(gomock.Any()).
						Return(nil, errors.New("Mock error"))

					return repo
				},
			},
			expectedErr:    true,
			expectedErrMsg: "Mock error",
		},
		{
			name: "when_repo_not_error_should_return_data",
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
					repo.EXPECT().
						FindAll(gomock.Any()).
						Return(&[]model.APIKeyDTO{
							{Token: token1, Name: "apikey1"},
							{Token: token2, Name: "apikey2"},
						}, nil)

					return repo
				},
			},
			expected: []model.APIKeyDTO{
				{Token: token1, Name: "apikey1"},
				{Token: token2, Name: "apikey2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := apikey.ProvideService(&config.Configuration{}, test.dependency.repo(ctrl))
			defer apikey.ResetService()

			actual, err := s.FindAll(ctx)
			if test.expectedErr && assert.Error(t, err) {
				assert.Equal(t, test.expectedErrMsg, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, actual)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Apikey_serviceImpl_GetItem(t *testing.T) {
	type dependency struct {
		repo func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO]
	}

	ctx := context.Background()
	mockToken := uuid.New().String()

	tests := []struct {
		name string
		pk   string
		dependency
		expectedErr    bool
		expectedErrMsg string
		expected       *model.APIKeyDTO
	}{
		{
			name: "when_repo_error_should_return_nil",
			pk:   mockToken,
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
					repo.EXPECT().
						FindByID(ctx, gomock.Any()).
						Return(nil, errors.New("Mock error"))

					return repo
				},
			},
			expectedErr:    true,
			expectedErrMsg: "Mock error",
		},
		{
			name: "when_repo_not_error_should_return_data",
			pk:   mockToken,
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
					repo.EXPECT().
						FindByID(ctx, gomock.Any()).
						Return(&model.APIKeyDTO{Token: mockToken, Name: "apikey"}, nil)

					return repo
				},
			},
			expected: &model.APIKeyDTO{Token: mockToken, Name: "apikey"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := apikey.ProvideService(&config.Configuration{}, test.dependency.repo(ctrl))
			defer apikey.ResetService()

			actual, err := s.FindByID(ctx, test.pk)
			if test.expectedErr && assert.Error(t, err) {
				assert.Equal(t, test.expectedErrMsg, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Apikey_serviceImpl_DeleteItem(t *testing.T) {
	type dependency struct {
		repo func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO]
	}

	ctx := context.Background()
	mockToken := uuid.New().String()

	tests := []struct {
		name string
		pk   string
		dependency
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			name: "when_repo_error_should_return_error",
			pk:   mockToken,
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
					repo.EXPECT().
						DeleteById(ctx, gomock.Any()).
						Return(errors.New("Mock error"))

					return repo
				},
			},
			expectedErr:    true,
			expectedErrMsg: "Mock error",
		},
		{
			name: "when_repo_not_error_should_return_no_error",
			pk:   mockToken,
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
					repo.EXPECT().
						DeleteById(ctx, gomock.Any()).
						Return(nil)

					return repo
				},
			},
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := apikey.ProvideService(&config.Configuration{}, test.dependency.repo(ctrl))
			defer apikey.ResetService()

			err := s.DeleteByID(ctx, test.pk)
			if test.expectedErr && assert.Error(t, err) {
				assert.Equal(t, test.expectedErrMsg, err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

// func Test_Apikey_serviceImpl_Create(t *testing.T) {
// 	type dependency struct {
// 		cfg  *config.Configuration
// 		repo func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO]
// 	}

// 	ctx := context.Background()
// 	mockCfg := &config.Configuration{
// 		SecretKey: "testsecretkey",
// 	}

// 	tests := []struct {
// 		name          string
// 		dto           *model.APIKeyDTO
// 		dependency
// 		expectedErr    bool
// 		expectedErrMsg string
// 	}{
// 		{
// 			name: "when_dto_is_nil_should_return_error",
// 			dto:  nil,
// 			dependency: dependency{
// 				cfg: mockCfg,
// 				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
// 					return nil
// 				},
// 			},
// 			expectedErr:    true,
// 			expectedErrMsg: "dto is nil",
// 		},
// 		{
// 			name: "when_generate_token_fails_should_return_error",
// 			dto: &model.APIKeyDTO{
// 				Name:     "testname",
// 				Duration: model.DurationUnlimited,
// 			},
// 			dependency: dependency{
// 				cfg: mockCfg,
// 				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
// 					return nil
// 				},
// 			},
// 			expectedErr:    true,
// 			expectedErrMsg: "key is empty",
// 		},
// 		{
// 			name: "when_repo_insert_fails_should_return_error",
// 			dto: &model.APIKeyDTO{
// 				Name:     "testname",
// 				Duration: model.DurationUnlimited,
// 			},
// 			dependency: dependency{
// 				cfg: mockCfg,
// 				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
// 					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
// 					repo.EXPECT().
// 						Insert(ctx, gomock.Any()).
// 						Return(errors.New("Mock error"))

// 					return repo
// 				},
// 			},
// 			expectedErr:    true,
// 			expectedErrMsg: "Mock error",
// 		},
// 		{
// 			name: "when_successful_should_return_token",
// 			dto: &model.APIKeyDTO{
// 				Name:     "testname",
// 				Duration: model.DurationUnlimited,
// 			},
// 			dependency: dependency{
// 				cfg: mockCfg,
// 				repo: func(ctrl *gomock.Controller) repo.Repo[model.APIKey, model.APIKeyDTO] {
// 					repo := mock.NewMockRepository[model.APIKey, model.APIKeyDTO](ctrl)
// 					repo.EXPECT().
// 						Insert(ctx, gomock.Any()).
// 						Return(nil)

// 					return repo
// 				},
// 			},
// 			expectedErr: false,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			s := apikey.ProvideService(test.dependency.cfg, test.dependency.repo(ctrl))
// 			defer apikey.ResetService()

// 			token, err := s.Create(ctx, test.dto)
// 			if test.expectedErr && assert.Error(t, err) {
// 				assert.Equal(t, test.expectedErrMsg, err.Error())
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.NotEmpty(t, token)
// 		})
// 	}
// }
