package entity

import (
	"codebase-app/pkg/types"
	"strconv"
	"strings"
	"time"
)

type XxxRequest struct {
}
type XxxResponse struct {
}
type XxxResult struct {
}

// type GetProductByIdRequest struct {
// 	Id string `validate:"uuid" db:"p.id"`
// }

// type GetProductByIdResponse struct {
// 	Id          string `json:"id" db:"id"`
// 	ShopId      string `json:"shop_id" db:"shop_id"`
// 	Name        string `json:"name" db:"name"`
// 	Description string `json:"description" db:"description"`
// 	// Category    string `json:"category" db:"category"`
// 	Price       string   `json:"price" db:"price"`
// 	Stocks      string   `json:"stocks" db:"stocks"`
// 	CategoryIds []string `json:"category_ids" db:"category_id"`
// }

type CreateProductRequest struct {
	UserId string `validate:"uuid" db:"user_id"`
	ShopId string `validate:"uuid" db:"shop_id"`

	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Price       float64  `json:"price" db:"price"`
	Stocks      int      `json:"stocks" db:"stocks"`
	CategoryIds []string `json:"category_ids" validate:"min=1,max=5"` // ["UUID-1", "UUID-2"]
}

type CreateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type UpdateProductByIdRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`
	ShopId string `prop:"shop_id" validate:"uuid" db:"shop_id"`

	Id string `params:"id" validate:"uuid" db:"id"`

	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	// Categories  string  `json:"categories" db:"categories"`
	Price      float64 `json:"price" db:"price"`
	Stocks     int     `json:"stocks" db:"stocks"`
	SoftDelete bool    `json:"soft_delete" db:"soft_delete"`
}

type UpdateProductByIdResponse struct {
	Id string `json:"id" db:"id"`
}

type DeleteProductByIdRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`
	ShopId string `prop:"shop_id" validate:"uuid" db:"shop_id"`

	Id string `validate:"uuid" db:"id"`
}

type GetProductRequest struct {
	ShopId     string `query:"shop_id" validate:"omitempty,uuid"`
	CategoryId string `query:"category_id" validate:"omitempty,uuid"`

	Name          string `query:"name"`
	PriceMinStr   string `query:"price_min" validate:"omitempty,numeric,gte=0"`
	PriceMaxStr   string `query:"price_max" validate:"omitempty,numeric,gte=0"`
	IsAvailable   bool   `query:"is_available"`
	ProductIdsStr string `query:"product_ids"`

	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`

	PriceMin   float64
	PriceMax   float64
	ProductIds []string
}

func (r *GetProductRequest) SetProductDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}

	if r.ProductIdsStr != "" {
		// split product ids string by comma
		ids := strings.Split(r.ProductIdsStr, ",")
		r.ProductIds = append(r.ProductIds, ids...)
	} else {
		r.ProductIds = make([]string, 0)
	}
}

func (r *GetProductRequest) CostumValidation() (int, map[string][]string) {
	var (
		errors   = make(map[string][]string)
		err      error
		priceMin float64
		priceMax float64
	)

	if r.PriceMinStr != "" {
		priceMin, err = strconv.ParseFloat(r.PriceMinStr, 64)
		if err != nil {
			errors["price_min"] = append(errors["price_min"], "price_min must be a number.")
		}
		r.PriceMin = priceMin
	}

	if r.PriceMaxStr != "" {
		priceMax, err = strconv.ParseFloat(r.PriceMaxStr, 64)
		if err != nil {
			errors["price_max"] = append(errors["price_max"], "price_max must be a number.")
		}
		r.PriceMax = priceMax
	}

	if len(errors) > 0 {
		return 400, errors
	}

	errors = nil
	return 0, errors
}

type GetProductResponse struct {
	Items []GetProductItem `json:"items"`
	Meta  types.Meta       `json:"meta"`
}

type GetProductItem struct {
	Id          string    `json:"id" db:"id"`
	ShopId      string    `json:"shop_id" db:"shop_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Stocks      int       `json:"stocks" db:"stocks"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	CategoryIds []string `json:"category_ids" db:"category_id"`
}
