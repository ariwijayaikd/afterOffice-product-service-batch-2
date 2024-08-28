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

func (r *productRepository) GetProductById(ctx context.Context, req *entity.GetProductByIdRequest) (*entity.GetProductByIdResponse, error) {
	var res = new(entity.GetProductByIdResponse)

	query := `
		SELECT 
			s.user_id as user_id,
			s.id as shop_id,
			s.name as shop_name,
			p.id as product_id,
			p.name as product_name,
			p.description as product_desc,
			p.categories as product_cat,
			p.price as product_price,
			p.stocks as product_stocks
		FROM
			products p 
			JOIN shops s ON s.id = p.shop_id 
		WHERE
			p.id = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(res)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DetailProductbyId - Failed to get DetailProductbyId")
		return nil, err
	}

	return res, nil
}

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
			COUNT(p.id) OVER() as total_data,
			s.id as shop_id, 
			s.name as shop_name, 
			p.id as product_id, 
			p.name as product_name, 
			p.description as product_description, 
			p.categories as product_categories, 
			p.price as product_price, 
			p.stocks as product_stocks
		FROM products p
		JOIN shops s ON s.id = p.shop_id
		WHERE p.name LIKE ? AND p.categories LIKE ?
		LIMIT ? OFFSET ?
		`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
		"'%"+req.Name+"%'",
		"'%"+req.Categories+"%'",
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
