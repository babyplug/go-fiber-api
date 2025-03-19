package app

import (
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/core/storage/db"
	"go-fiber-api/internal/wrapper/logx"

	apikey_middleware "go-fiber-api/internal/core/middleware/apikey"
	"go-fiber-api/internal/core/middleware/cache"

	"go-fiber-api/internal/feature/apikey"
	"go-fiber-api/internal/feature/user"

	"go-fiber-api/toolkit/errorhandler"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config           *config.Configuration
	LogX             *logx.LogX
	Server           *fiber.App
	DBClient         db.Client
	UserHandler      user.Handler
	CacheMiddleware  cache.CacheMiddleware
	APIKeyHandler    apikey.Handler
	APIKeyMiddleware apikey_middleware.Middleware
}

var (
	app     = &Application{}
	appOnce sync.Once
)

func Provide(
	cfg *config.Configuration,
	log *logx.LogX,
	dbClient db.Client,
	userHandler user.Handler,
	cacheMiddleware cache.CacheMiddleware,
	apiKeyHandler apikey.Handler,
	apiKeyMiddleware apikey_middleware.Middleware,
) *Application {
	appOnce.Do(func() {
		app = &Application{
			Config:           cfg,
			LogX:             log,
			Server:           getServer(log),
			DBClient:         dbClient,
			UserHandler:      userHandler,
			CacheMiddleware:  cacheMiddleware,
			APIKeyHandler:    apiKeyHandler,
			APIKeyMiddleware: apiKeyMiddleware,
		}
		registerHandler(app)
	})

	return app
}

func getServer(log *logx.LogX) *fiber.App {
	server := fiber.New(
		fiber.Config{
			ErrorHandler: errorhandler.Handler(),
		},
	)

	// Helmet middleware helps secure your apps by setting various HTTP headers.
	server.Use(helmet.New())

	// Idempotency middleware for Fiber allows for fault-tolerant APIs where duplicate requests — for example due to networking issues on the client-side — do not erroneously cause the same action performed multiple times on the server-side.
	// ref: https://docs.gofiber.io/api/middleware/idempotency
	server.Use(idempotency.New())

	server.Use(requestid.New())

	server.Use(func(c *fiber.Ctx) error {
		log.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		}).Info("Incoming request")
		return c.Next()
	})

	// == Routes ==

	server.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return server
}

func registerHandler(app *Application) {
	app.Server.Use(app.CacheMiddleware.RedisCacheMiddleware(60))

	app.Server.Use(cors.New(cors.Config{
		AllowOrigins: app.Config.CorsAllowedOrigins, // Allow requests from frontend
		AllowHeaders: app.Config.CorsAllowedHeaders,
	}))
	root := app.Server.Group("")

	// Middleware: Set user access token to context
	root.Use(func(c *fiber.Ctx) error {
		// tk := c.Get(fiber.HeaderAuthorization)
		// c.Context().SetUserValue(constants.ContextUserAccessTokenKey, tk)
		return c.Next()
	})

	// games := root.Group("users")
	// games.Get("", app.UserHandler)
	// games.Get(":id", app.UserHandler.GetByID)
	// games.Post("", app.UserHandler.Create)
	// games.Put(":id", app.UserHandler.Update)
	// games.Delete(":id", app.UserHandler.Delete)

	service := root.Group("service")
	service.Use(app.APIKeyMiddleware.Validate())
	service.Get("", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello, world")
	})

	// adminApi.Get("/me", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"data": fiber.Map{
	// 			"name": "John Doe",
	// 			"age":  25,
	// 			"role": "admin",
	// 		},
	// 	})
	// })
}

// Only for test
func Reset() {
	appOnce = sync.Once{}
}
