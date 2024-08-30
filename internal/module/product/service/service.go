package service

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"
)

var _ ports.ProductService = &productService{}

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	return s.repo.CreateProduct(ctx, req)
}

// func (s *productService) GetProductById(ctx context.Context, req *entity.GetProductByIdRequest) (*entity.GetProductByIdResponse, error) {
// 	return s.repo.GetProductById(ctx, req)
// }

func (s *productService) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	return s.repo.GetProduct(ctx, req)
}

func (s *productService) DeleteProductById(ctx context.Context, req *entity.DeleteProductByIdRequest) error {
	return s.repo.DeleteProductById(ctx, req)
}

func (s *productService) UpdateProductById(ctx context.Context, req *entity.UpdateProductByIdRequest) (*entity.UpdateProductByIdResponse, error) {
	return s.repo.UpdateProductById(ctx, req)
}
