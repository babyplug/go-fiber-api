package cache_test

import (
	"encoding/json"
	"fmt"
	"go-fiber-api/internal/core/middleware/cache"
	"go-fiber-api/internal/core/response"
	"go-fiber-api/internal/mock"
	"go-fiber-api/internal/wrapper/redis"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRedisCacheMiddleware(t *testing.T) {
	type dependency struct {
		redisClient func(*gomock.Controller) redis.Client
	}

	tests := []struct {
		name   string
		method string
		url    string
		body   string
		dependency
		statusCode       int
		cacheHeader      string
		expectedResponse string
	}{
		{
			name:   "should cache GET request",
			method: http.MethodGet,
			url:    "/test",
			dependency: dependency{
				redisClient: func(ctrl *gomock.Controller) redis.Client {
					m := mock.NewMockRedisClient(ctrl)
					cacheKey := "cache:GET:/test:"
					m.EXPECT().Get(cacheKey, gomock.Any()).Return(fmt.Errorf("cache miss"))
					m.EXPECT().Set(cacheKey, []byte("{\"message\":\"get success\",\"data\":null}"), 60).Return(nil)
					return m
				},
			},
			statusCode:       http.StatusOK,
			cacheHeader:      "MISS",
			expectedResponse: "map[data:<nil> message:get success]",
		},
		{
			name:   "should return cached response for GET request",
			method: http.MethodGet,
			url:    "/test",
			dependency: dependency{
				redisClient: func(ctrl *gomock.Controller) redis.Client {
					m := mock.NewMockRedisClient(ctrl)

					cacheKey := "cache:GET:/test:"
					m.EXPECT().Get(cacheKey, gomock.Any()).SetArg(1, response.ResponseDTO{
						Message: "Cache get success",
					}).Return(nil)

					return m
				},
			},
			statusCode:       http.StatusOK,
			cacheHeader:      "HIT",
			expectedResponse: "map[data:<nil> message:Cache get success]",
		},
		{
			name:   "should clear cache for POST request",
			method: http.MethodPost,
			url:    "/test",
			body:   `{}`,
			dependency: dependency{
				redisClient: func(ctrl *gomock.Controller) redis.Client {
					m := mock.NewMockRedisClient(ctrl)

					cacheKey := "cache:GET:/test:"
					m.EXPECT().Del(cacheKey).Return(nil)

					return m
				},
			},
			statusCode:       http.StatusOK,
			expectedResponse: "map[data:<nil> message:post success]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := cache.New(tt.redisClient(ctrl))
			defer cache.Close()

			app := fiber.New()
			app.Use(c.RedisCacheMiddleware(60))

			app.Get("/test", func(c *fiber.Ctx) error {
				return c.JSON(&response.ResponseDTO{
					Message: "get success",
				})
			})

			app.Post("/test", func(c *fiber.Ctx) error {
				return c.JSON(&response.ResponseDTO{
					Message: "post success",
				})
			})

			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}
			tenSecond := 10 * time.Second
			resp, err := app.Test(req, int(tenSecond.Milliseconds()))

			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)

			// check response jsonBody
			jsonBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Variable to hold the unmarshaled data
			var data map[string]interface{}

			// Unmarshal the JSON data into the map
			err = json.Unmarshal(jsonBody, &data)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return
			}
			v := fmt.Sprintf("%+v", data)

			assert.Equal(t, tt.expectedResponse, v)

			if tt.cacheHeader != "" {
				assert.Equal(t, tt.cacheHeader, resp.Header.Get("X-Cache"))
			}
		})
	}
}
