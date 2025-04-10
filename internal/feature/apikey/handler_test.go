package apikey_test

import (
	"errors"
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/feature/apikey"
	"go-fiber-api/internal/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApiKey_Handler_Create(t *testing.T) {
	type dependency struct {
		s func(ctrl *gomock.Controller) apikey.Service
	}

	mockToken := "mockToken"
	// mockDTO := &model.APIKeyDTO{Name: "test", Duration: model.DurationUnlimited}

	tests := []struct {
		name           string
		body           string
		dependency     dependency
		expectedError  bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "when_body_is_invalid_should_return_400",
			body: `{"name": "test", "duration":}`,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					return nil
				},
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
			expectedBody:   `invalid character '}' looking for beginning of value`,
		},
		{
			name: "when_create_fails_should_return_500",
			body: `{"name": "test", "duration": "unlimited"}`,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", errors.New("mock error"))
					return m
				},
			},
			expectedError:  true,
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `mock error`,
		},
		{
			name: "when_successful_should_return_token",
			body: `{"name": "test", "duration": "unlimited"}`,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockToken, nil)
					return m
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"message":"success","data":"mockToken"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := apikey.ProvideHandler(test.dependency.s(ctrl))
			defer apikey.ResetHandler()

			app := fiber.New()
			app.Post("/apikey", h.Create)

			req := httptest.NewRequest(http.MethodPost, "/apikey", strings.NewReader(test.body))
			req.Header.Add("Content-Type", "application/json")

			resp, err := app.Test(req)
			bodyBytes, _ := io.ReadAll(resp.Body)
			actual := string(bodyBytes)

			assert.Equal(t, test.expectedStatus, resp.StatusCode)
			if test.expectedError {
				assert.Equal(t, test.expectedBody, actual)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedBody, actual)
		})
	}
}

func TestApiKey_Handler_GetAll(t *testing.T) {
	type dependency struct {
		s func(ctrl *gomock.Controller) apikey.Service
	}

	mockData := []model.APIKeyDTO{
		{Base: model.Base{ID: 1}, Token: "token1", Name: "apikey1", Duration: model.DurationSevenDays},
		{Base: model.Base{ID: 2}, Token: "token2", Name: "apikey2", Duration: model.DurationUnlimited},
	}

	tests := []struct {
		name           string
		dependency     dependency
		expectedErr    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "when_query_fails_should_return_500",
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("mock error"))
					return m
				},
			},
			expectedErr:    true,
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `mock error`,
		},
		{
			name: "when_successful_should_return_data",
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().FindAll(gomock.Any()).Return(mockData, nil)
					return m
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"message":"success","data":[{"id":1,"token":"token1","name":"apikey1","duration":"7_DAYS","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"},{"id":2,"token":"token2","name":"apikey2","duration":"UNLIMITED","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := apikey.ProvideHandler(test.dependency.s(ctrl))
			defer apikey.ResetHandler()

			app := fiber.New()
			app.Get("/apikeys", h.FindAll)

			req := httptest.NewRequest(http.MethodGet, "/apikeys", nil)
			req.Header.Add("Content-Type", "application/json")

			resp, err := app.Test(req)
			bodyBytes, _ := io.ReadAll(resp.Body)
			actual := string(bodyBytes)

			assert.Equal(t, test.expectedStatus, resp.StatusCode)
			if test.expectedErr {
				assert.Equal(t, test.expectedBody, actual)
				return
			}

			assert.NoError(t, err)
			assert.JSONEq(t, test.expectedBody, actual)
		})
	}
}

func TestApiKey_Handler_GetItem(t *testing.T) {
	type dependency struct {
		s func(ctrl *gomock.Controller) apikey.Service
	}

	mockToken := "mockToken"
	mockData := &model.APIKeyDTO{Base: model.Base{ID: 1}, Token: mockToken, Name: "apikey"}

	tests := []struct {
		name           string
		pathParam      string
		dependency     dependency
		expectedErr    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "when_get_item_fails_should_return_404",
			pathParam: mockToken,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(model.APIKeyDTO{}, errors.New("mock error"))
					return m
				},
			},
			expectedErr:    true,
			expectedStatus: fiber.StatusNotFound,
			expectedBody:   `mock error`,
		},
		{
			name:      "when_successful_should_return_data",
			pathParam: mockToken,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(*mockData, nil)
					return m
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"message":"success","data":{"id":1,"token":"mockToken","name":"apikey","createdAt":0,"duration":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := apikey.ProvideHandler(test.dependency.s(ctrl))
			defer apikey.ResetHandler()

			app := fiber.New()
			app.Get("/apikey/:token", h.FindOne)

			req := httptest.NewRequest(http.MethodGet, "/apikey/"+test.pathParam, nil)
			req.Header.Add("Content-Type", "application/json")

			resp, err := app.Test(req)
			bodyBytes, _ := io.ReadAll(resp.Body)
			actual := string(bodyBytes)

			assert.Equal(t, test.expectedStatus, resp.StatusCode)
			if test.expectedErr {
				assert.Equal(t, test.expectedBody, actual)
				return
			}

			assert.NoError(t, err)
			assert.JSONEq(t, test.expectedBody, actual)
		})
	}
}

func TestApiKey_Handler_DeleteItem(t *testing.T) {
	type dependency struct {
		s func(ctrl *gomock.Controller) apikey.Service
	}

	mockToken := "mockToken"

	tests := []struct {
		name           string
		pathParam      string
		dependency     dependency
		expectedErr    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "when_delete_item_fails_should_return_500",
			pathParam: mockToken,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
					return m
				},
			},
			expectedErr:    true,
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `mock error`,
		},
		{
			name:      "when_successful_should_return_success",
			pathParam: mockToken,
			dependency: dependency{
				s: func(ctrl *gomock.Controller) apikey.Service {
					m := mock.NewMockAPIKeyService(ctrl)
					m.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(nil)
					return m
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"data":null,"message":"success"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := apikey.ProvideHandler(test.dependency.s(ctrl))
			defer apikey.ResetHandler()

			app := fiber.New()
			app.Delete("/apikey/:token", h.DeleteByID)

			req := httptest.NewRequest(http.MethodDelete, "/apikey/"+test.pathParam, nil)
			req.Header.Add("Content-Type", "application/json")

			resp, err := app.Test(req)
			bodyBytes, _ := io.ReadAll(resp.Body)
			actual := string(bodyBytes)

			assert.Equal(t, test.expectedStatus, resp.StatusCode)
			if test.expectedErr {
				assert.Equal(t, test.expectedBody, actual)
				return
			}

			assert.NoError(t, err)
			assert.JSONEq(t, test.expectedBody, actual)
		})
	}
}
