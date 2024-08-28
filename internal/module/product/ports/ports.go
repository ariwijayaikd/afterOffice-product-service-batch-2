package ports

import (
	"codebase-app/internal/module/product/entity"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error)
	// DetailProductbyId(ctx context.Context, req *entity.GetDetailProductRequest) (*entity.GetDetailProductResponse, error)
	// DeleteProductbyId(ctx context.Context, req *entity.DeleteProductRequest) error
	// UpdateProductbyId(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error)
	// SearchProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error)
	GetAllProduct(ctx context.Context, req *entity.GetAllProductRequest) (*entity.GetAllProductResponse, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error)
	// DetailProductbyId(ctx context.Context, req *entity.GetDetailProductRequest) (*entity.GetDetailProductResponse, error)
	GetAllProduct(ctx context.Context, req *entity.GetAllProductRequest) (*entity.GetAllProductResponse, error)
}
