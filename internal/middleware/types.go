package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	UserId string
	Role   string
}

func GetLocals(c *fiber.Ctx) *Locals {
	var l = Locals{}
	userId, ok := c.Locals("user_id").(string)
	if ok {
		l.UserId = userId
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	return &l
}

func (l *Locals) GetUserId() string {
	return l.UserId
}

func (l *Locals) GetRole() string {
	return l.Role
}

// for create product with X-SHOP-ID
type ShopLocals struct {
	ShopId string
	// Role   string
}

func GetShopLocals(c *fiber.Ctx) *ShopLocals {
	var l = ShopLocals{}
	shopId, ok := c.Locals("shop_id").(string)
	if ok {
		l.ShopId = shopId
	} else {
		log.Warn().Msg("middleware::Locals-GetShopLocals failed to get shop_id from locals")
	}
	return &l
}

func (l *ShopLocals) GetShopId() string {
	return l.ShopId
}

// for get all product
type ProductLocals struct {
	Name       string
	Categories string
}

func GetProductLocals(c *fiber.Ctx) *ProductLocals {
	var l = ProductLocals{}
	name, ok := c.Locals("name").(string)
	if ok {
		l.Name = name
	} else {
		log.Warn().Msg("middleware::Locals-GetProductLocals failed to get name from locals")
	}
	return &l
}

func (l *ProductLocals) GetProductName() string {
	return l.Name
}
func (l *ProductLocals) GetProductCategories() string {
	return l.Categories
}
