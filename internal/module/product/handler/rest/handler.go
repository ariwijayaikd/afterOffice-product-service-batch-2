package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/internal/module/product/repository"
	"codebase-app/internal/module/product/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type productHandler struct {
	service ports.ProductService
}

func NewProductHandler() *productHandler {
	var (
		handler = new(productHandler)
		repo    = repository.NewProductRepository(adapter.Adapters.ShopeefunPostgres)
		service = service.NewProductService(repo)
	)
	handler.service = service

	return handler
}

func (h *productHandler) Register(router fiber.Router) {
	router.Post("/", middleware.UserIdHeader, middleware.ShopIdHeader, h.CreateProduct)
	router.Get("/q", h.GetProduct)
	// router.Get("/:id", h.GetProductById)
	router.Delete("/:id", middleware.UserIdHeader, middleware.ShopIdHeader, h.DeleteProductById)
	router.Patch("/:id", middleware.UserIdHeader, middleware.ShopIdHeader, h.UpdateProductById)
}

func (h *productHandler) CreateProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateProductRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
		sl  = middleware.GetShopLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateProduct - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.ShopId = sl.ShopId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateProduct - Validator request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))
}

func (h *productHandler) GetProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		// l   = middleware.GetLocals(c)
	)

	// req.UserId = l.UserId
	// req.Name = c.Params("name")
	// req.Categories = c.Params("categories")

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAllProduct - Parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetProductDefault()

	if code, err := req.CostumValidation(); code != 0 {
		return c.Status(code).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetAllProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

// func (h *productHandler) GetProductById(c *fiber.Ctx) error {
// 	var (
// 		req = new(entity.GetProductByIdRequest)
// 		ctx = c.Context()
// 		// v   = adapter.Adapters.Validator
// 	)

// 	req.Id = c.Params("id")

// 	// if err := v.Validate(req); req != nil {
// 	// 	log.Warn().Err(err).Any("payload", req).Msg("handler::GetProductbyId - Validate request body")
// 	// 	code, errs := errmsg.Errors(err, req)
// 	// 	return c.Status(code).JSON(response.Error(errs))
// 	// }

// 	resp, err := h.service.GetProductById(ctx, req)
// 	if err != nil {
// 		code, errs := errmsg.Errors[error](err)
// 		return c.Status(code).JSON(response.Error(errs))
// 	}

// 	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
// }

func (h *productHandler) DeleteProductById(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteProductByIdRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
		sl  = middleware.GetShopLocals(c)
	)

	req.UserId = l.UserId
	req.ShopId = sl.ShopId
	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeleteProductById - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteProductById(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *productHandler) UpdateProductById(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateProductByIdRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
		sl  = middleware.GetShopLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateProductById - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.ShopId = sl.ShopId
	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::UpdateProductById - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.UpdateProductById(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
