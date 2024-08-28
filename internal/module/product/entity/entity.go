package entity

import "codebase-app/pkg/types"

type XxxRequest struct {
}
type XxxResponse struct {
}
type XxxResult struct {
}

type GetDetailProductRequest struct {
	Id string `json:"id" db:"id"`
	// ShopId string `json:"shop_id" db:"shop_id"`
}

type GetDetailProductResponse struct {
	Id          string `json:"id" db:"id"`
	ShopId      string `json:"shop_id" db:"shop_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Categories  string `json:"categories" db:"categories"`
	Price       string `json:"price" db:"price"`
	Stocks      string `json:"stocks" db:"stocks"`
	SoftDelete  bool   `json:"soft_delete" db:"soft_delete"`
}

type CreateProductRequest struct {
	// UserId string `validate:"uuid" db:"user_id"`
	ShopId string `validate:"uuid" db:"shop_id"`

	// ShopId      string `json:"shop_id" db:"shop_id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Categories  string  `json:"categories" db:"categories"`
	Price       float64 `json:"price" db:"price"`
	Stocks      int     `json:"stocks" db:"stocks"`
}

type CreateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type UpdateProductRequest struct {
	Id          string `json:"id" db:"id"`
	ShopId      string `json:"shop_id" db:"shop_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Categories  string `json:"categories" db:"categories"`
	Price       string `json:"price" db:"price"`
	Stocks      string `json:"stocks" db:"stocks"`
	SoftDelete  bool   `json:"soft_delete" db:"soft_delete"`
}

type UpdateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type DeleteProductRequest struct {
	Id     string `json:"id" db:"id"`
	ShopId string `json:"shop_id" db:"shop_id"`
}

type GetAllProductRequest struct {
	Name       string `json:"name" db:"p.name"`
	Categories string `json:"categories" db:"p.categories"`
	Page       int    `query:"page" validate:"required"`
	Paginate   int    `query:"paginate" validate:"required"`
}

func (r *GetAllProductRequest) SetProductDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetAllProductItem struct {
	Id          string  `json:"id" db:"s.id"`
	ShopName    string  `json:"shop_name" db:"s.name"`
	ProductId   string  `json:"product_id" db:"p.id"`
	ProductName string  `json:"product_name" db:"p.name"`
	Description string  `json:"description" db:"p.description"`
	Categories  string  `json:"categories" db:"p.categories"`
	Price       float64 `json:"price" db:"p.price"`
	Stocks      int     `json:"stocks" db:"p.stocks"`
}

type GetAllProductResponse struct {
	Items []GetAllProductItem `json:"items"`
	Meta  types.Meta          `json:"meta"`
}
