package apikey

import (
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/feature/apikey"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var (
	m     *middlewareImpl
	mOnce sync.Once
	log   *logrus.Entry
)

type Middleware interface {
	Validate() fiber.Handler
}

type middlewareImpl struct {
	cfg *config.Configuration
	s   apikey.Service
}

func Provide(cfg *config.Configuration, apiKeySvc apikey.Service) Middleware {
	mOnce.Do(func() {
		m = &middlewareImpl{
			cfg: cfg,
			s:   apiKeySvc,
		}
	})

	return m
}

func (m *middlewareImpl) Validate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := strings.Trim(c.Get("X-API-Key"), " ")
		if len(key) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		// Parse and validate token
		token, err := jwt.Parse(key, func(token *jwt.Token) (any, error) {
			return []byte(m.cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired api key"})
		}

		// If token valid, we then check is api key is revoke or not
		apiKey, err := m.s.FindByID(c.Context(), key)
		if err != nil || apiKey == (model.APIKeyDTO{}) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Api key not found or was been revoked",
			})
		}

		return c.Next()
	}
}
