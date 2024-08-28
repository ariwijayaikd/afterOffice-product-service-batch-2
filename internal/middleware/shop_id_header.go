package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ShopIdHeader(c *fiber.Ctx) error {
	shopId := c.Get("X-SHOP-ID")
	unauthorizedResponse := fiber.Map{
		"message": "Unauthorized",
		"success": false,
	}

	if shopId == "" {
		log.Error().Msg("middleware::ShopIdHeader - Unauthorized [Header not set]")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	c.Locals("shop_id", shopId)

	return c.Next()
}
