package fiber

import (
	"context"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-new/internal/http/middlewares"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strings"
	"time"
)

// @title Tsuwari api
// @version 1.0
// @description Non-public api for tsuwari
// @host localhost:3002
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api-key
func NewFiber(
	cfg *config.Config,
	middlewares *middlewares.Middlewares,
	logger *zap.SugaredLogger,
	notSugaredLogger *zap.Logger,
	redisCacheStorage *RedisCacheStorage,
	lc fx.Lifecycle,
) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          middlewares.ErrorHandler,
		EnablePrintRoutes:     true,
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		BodyLimit:             20 * 1024 * 1024,
	})
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: notSugaredLogger,
		Fields: []string{"latency", "status", "method", "url", "body", "queryParams"},
	}))

	cacheExpirationTime := lo.
		If(cfg.AppEnv == "development", 1*time.Second).
		Else(24 * time.Hour)

	app.Use(cache.New(cache.Config{
		Expiration:   cacheExpirationTime,
		CacheControl: false,
		KeyGenerator: func(c *fiber.Ctx) string {
			return redisCacheStorage.BuildKey(c.Path())
		},
		Storage: redisCacheStorage,
		Next: func(c *fiber.Ctx) bool {
			notCachedPaths := []string{
				"/metrics",
				"/auth",
			}
			excludeFromNotCached := []string{"/auth/profile", "/auth/dashboard"}
			p := c.Path()

			if lo.SomeBy(excludeFromNotCached, func(path string) bool {
				return strings.HasPrefix(p, path)
			}) {
				return true
			}

			return lo.SomeBy(notCachedPaths, func(path string) bool {
				return strings.HasPrefix(p, path)
			})
		},
	}))
	app.Use(helmet.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	prometheus := fiberprometheus.New("twir-api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Best bot api root. ;)")
	})

	if cfg.AppEnv == "development" {
		app.Get("/swagger/*", swagger.New(swagger.Config{
			URL:                  "http://localhost:3002/swagger/doc.json",
			DeepLinking:          false,
			DocExpansion:         "list",
			PersistAuthorization: true,
			Title:                "Twir api",
			TryItOutEnabled:      true,
		}))
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen("0.0.0.0:3002"); err != nil {
					logger.Fatalln(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}
