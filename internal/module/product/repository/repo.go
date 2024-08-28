package repository

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var res = new(entity.CreateProductResponse)

	query := `
		INSERT INTO products (shop_id, name, description, categories, price, stocks)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.Name,
		req.Description,
		req.Categories,
		req.Price,
		req.Stocks).Scan(&res.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err
	}

	return res, nil
}

// func (r *productRepository) DetailProductbyId(ctx context.Context, req *entity.GetDetailProductRequest) (*entity.GetDetailProductResponse, error) {
// 	var res = new(entity.GetDetailProductResponse)

// 	query := `
// 		SELECT p.id, s.shop_id, s.name, p.name, p.description, p.categories, p.price, p.stocks, p.created_at, p.updated_at, p.deleted_at
// 		FROM products p
// 		JOIN shops s ON s.id = p.shop_id
// 		WHERE p.id = ?
// 	`

// 	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(res)
// 	if err != nil {
// 		log.Error().Err(err).Any("payload", req).Msg("repository::DetailProductbyId - Failed to get DetailProductbyId")
// 		return nil, err
// 	}

// 	return res, nil
// }

func (r *productRepository) GetAllProduct(ctx context.Context, req *entity.GetAllProductRequest) (*entity.GetAllProductResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.GetAllProductItem
	}

	var (
		res  = new(entity.GetAllProductResponse)
		data = make([]dao, 0, req.Paginate)
	)
	res.Items = make([]entity.GetAllProductItem, 0, req.Paginate)

	query := `
		SELECT 
			COUNT (p.id) OVER() as total_data,
			s.id, 
			s.name, 
			p.id, 
			p.name, 
			p.description, 
			p.categories, 
			p.price, p.stocks
		FROM products p
		JOIN shops s on s.id = p.shop_id
		WHERE p.name LIKE ? and p.categories LIKE ?
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Categories,
		// "'%"+req.Name+"%'",
		// "'%"+req.Categories+"%'",
		req.Paginate,
		req.Paginate*(req.Page-1),
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetAllProduct - Failed to get all product")
		return nil, err
	}

	if len(data) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}
	for _, d := range data {
		res.Items = append(res.Items, d.GetAllProductItem)
	}
	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
