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

// func (s *productService) DetailProductbyId(ctx context.Context, req *entity.GetDetailProductRequest) (*entity.GetDetailProductResponse, error) {
// 	return s.repo.DetailProductbyId(ctx, req)
// }

func (s *productService) GetAllProduct(ctx context.Context, req *entity.GetAllProductRequest) (*entity.GetAllProductResponse, error) {
	return s.repo.GetAllProduct(ctx, req)
}
