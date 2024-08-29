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

	queryProduct := `
		INSERT INTO products (shop_id, name, description, price, stocks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, now(), now())
		RETURNING id
	`

	// queryProductCategory := `
	// 	INSERT INTO product_category (product_id, category_id)
	// 	VALUES (?, ?);
	// `
	err := r.db.QueryRowxContext(ctx, r.db.Rebind(queryProduct),
		req.ShopId,
		req.Name,
		req.Description,
		// req.Category,
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
			id,
			shop_id,
			name,
			description,
			price,
			stocks
		FROM
			products
		WHERE
			id = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(res)
	if err != nil {
		// if check error sql no rows

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
		data = make([]dao, 0)
		// data = make([]dao, 0, req.Paginate)
		arg = make(map[string]any)
	)
	// res.Items = make([]entity.GetAllProductItem, 0, req.Paginate)
	res.Meta.Page = req.Page
	res.Meta.Paginate = req.Paginate

	query := `
		SELECT 
			COUNT(id) OVER() as total_data,
			id, 
			shop_id,  
			name, 
			description, 
			price, 
			stocks,
			created_at, 
			updated_at
		FROM products 
		WHERE soft_delete = false
	`

	if len(req.ProductIds) > 0 {
		var ids string
		for i, id := range req.ProductIds {
			ids += "'" + id + "'"
			if i < len(req.ProductIds)-1 {
				ids += ", "
			}
		}
		query += " AND id IN (" + ids + ")"
	}

	if req.ShopId != "" {
		query += " AND shop_id = :shop_id"
		arg["shop_id"] = req.ShopId
	}

	// if req.CategoryId != "" {
	// 	query += " AND category_id = :category_id"
	// 	arg["category_id"] = req.CategoryId
	// }

	if req.Name != "" {
		query += " AND name ILIKE '%' || :name || '%'"
		arg["name"] = req.Name
	}

	if req.PriceMinStr != "" {
		query += " AND price >= :price_min"
		arg["price_min"] = req.PriceMin
	}

	if req.PriceMaxStr != "" {
		query += " AND price <= :price_max"
		arg["price_max"] = req.PriceMax
	}

	if req.IsAvailable {
		query += " AND stock > 0"
	}

	query += `
		ORDER BY created_at DESC
		LIMIT :paginate
		OFFSET :offset
	`
	arg["paginate"] = req.Paginate
	arg["offset"] = (req.Page - 1) * req.Paginate

	nstmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetAllProduct - Failed to prepare query")
		return nil, err
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &data, arg)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetAllProduct - Failed to get all product")
		return nil, err
	}

	for _, d := range data {
		res.Items = append(res.Items, entity.GetAllProductItem{
			Id:     d.Id,
			ShopId: d.ShopId,
			// CategoryId:  d.CategoryId,
			Name:        d.Name,
			Description: d.Description,
			Price:       d.Price,
			Stocks:      d.Stocks,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
		res.Meta.TotalData = d.TotalData
	}

	res.Meta.CountTotalPage(res.Meta.Page, res.Meta.Paginate, res.Meta.TotalData)
	return res, nil
}

func (r *productRepository) DeleteProductById(ctx context.Context, req *entity.DeleteProductByIdRequest) error {
	query := `
	    UPDATE products
		SET deleted_at = NOW()
		-- FROM shops s
		-- JOIN shops on s.id = p.shop_id
		WHERE id = ? AND shop_id = ?
		-- AND s.user_id = '6f25adf4-2759-426d-af20-2be82f5f728c'
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id, req.ShopId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProductById - Failed to delete product")
		return err
	}

	return nil
}

func (r *productRepository) UpdateProductById(ctx context.Context, req *entity.UpdateProductByIdRequest) (*entity.UpdateProductByIdResponse, error) {
	var resp = new(entity.UpdateProductByIdResponse)

	query := `
		UPDATE products
			SET name = ?, description = ?, price = ?, stocks = ?, soft_delete = ?, updated_at = NOW()
		WHERE id = ? and shop_id = ?
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Description,
		// req.Categories,
		req.Price,
		req.Stocks,
		req.SoftDelete,
		req.Id,
		req.ShopId).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProductById - Failed to update product")
		return nil, err
	}

	return resp, nil
}
