package entity

import "codebase-app/pkg/types"

type XxxRequest struct {
}
type XxxResponse struct {
}
type XxxResult struct {
}

type GetProductByIdRequest struct {
	Id string `json:"product_id" db:"p.id"`
	// ShopId string `json:"shop_id" db:"shop_id"`
}

type GetProductByIdResponse struct {
	UserId      string `json:"user_id" db:"user_id"`
	ShopId      string `json:"shop_id" db:"shop_id"`
	ShopName    string `json:"shop_name" db:"shop_name"`
	Id          string `json:"product_id" db:"product_id"`
	Name        string `json:"product_name" db:"product_name"`
	Description string `json:"product_desc" db:"product_desc"`
	Categories  string `json:"product_cat" db:"product_cat"`
	Price       string `json:"product_price" db:"product_price"`
	Stocks      string `json:"product_stocks" db:"product_stocks"`
	// SoftDelete  bool   `json:"soft_delete" db:"soft_delete"`
}

type CreateProductRequest struct {
	UserId string `validate:"uuid" db:"user_id"`
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
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

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
	Name       string `json:"product_name" db:"p.name"`
	Categories string `json:"product_categories" db:"p.categories"`
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
	Id          string  `json:"shop_id" db:"shop_id"`
	ShopName    string  `json:"shop_name" db:"shop_name"`
	ProductId   string  `json:"product_id" db:"product_id"`
	ProductName string  `json:"product_name" db:"product_name"`
	Description string  `json:"product_description" db:"product_description"`
	Categories  string  `json:"product_categories" db:"product_categories"`
	Price       float64 `json:"product_price" db:"product_price"`
	Stocks      int     `json:"product_stocks" db:"product_stocks"`
}

type GetAllProductResponse struct {
	Items []GetAllProductItem `json:"items"`
	Meta  types.Meta          `json:"meta"`
}
