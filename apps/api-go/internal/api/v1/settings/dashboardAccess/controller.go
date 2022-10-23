package dashboardaccess

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("dashboard-access")

	/* dashboardAccessList := cache.New(cache.Config{
		Expiration: 15 * time.Second,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("channels:dashboardAccess:%s", c.Params("channelId"))
		},
	}) */

	middleware.Get("" /* dashboardAccessList, */, func(c *fiber.Ctx) error {
		users, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return nil
		}

		return c.JSON(users)
	})

	middleware.Post("", func(c *fiber.Ctx) error {
		dto := &addUserDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		entity, err := handlePost(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(entity)
	})

	middleware.Delete(":entityId", func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("entityId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(200)
	})

	return middleware
}
