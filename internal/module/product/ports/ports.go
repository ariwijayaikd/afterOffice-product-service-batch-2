package ports

import (
	"codebase-app/internal/module/product/entity"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error)
	// GetProductById(ctx context.Context, req *entity.GetProductByIdRequest) (*entity.GetProductByIdResponse, error)
	DeleteProductById(ctx context.Context, req *entity.DeleteProductByIdRequest) error
	UpdateProductById(ctx context.Context, req *entity.UpdateProductByIdRequest) (*entity.UpdateProductByIdResponse, error)
	// SearchProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error)
	GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error)
	// GetProductById(ctx context.Context, req *entity.GetProductByIdRequest) (*entity.GetProductByIdResponse, error)
	DeleteProductById(ctx context.Context, req *entity.DeleteProductByIdRequest) error
	UpdateProductById(ctx context.Context, req *entity.UpdateProductByIdRequest) (*entity.UpdateProductByIdResponse, error)
	GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error)
}
