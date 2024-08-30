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
	// table ini merupakan relasi antara product dan category yang dimilikinya
	// product_category -> product_id, category_id -> product dan category

	// code lama
	// query := `
	// 	INSERT INTO products (shop_id, name, description, price, stocks, created_at, updated_at)
	// 	VALUES (?, ?, ?, ?, ?, now(), now())
	// 	RETURNING id
	// `

	// err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
	// 	req.ShopId,
	// 	req.Name,
	// 	req.Description,
	// 	// req.Category,
	// 	req.Price,
	// 	req.Stocks).Scan(&res.Id)
	// if err != nil {
	// 	log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
	// 	return nil, err
	// }

	// return res, nil

	//code baru

	tx, err := r.db.BeginTxx(ctx, nil)
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Error().Err(err).Msg("repository::CreateProduct - Failed to rollback transaction")
			}
		} else {
			err := tx.Commit()
			if err != nil {
				log.Error().Err(err).Msg("repository::CreateProduct - Failed to commit transaction")
			}
		}
	}()

	query := `
		INSERT INTO products (shop_id, name, description, price, stocks, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, now(), now())
        RETURNING id
	`

	err = tx.QueryRowContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.Name,
		req.Description,
		req.Price,
		req.Stocks).Scan(&res.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err
	}

	query = `
		INSERT INTO products_category (product_id, category_id)
        VALUES (?, ?)
	`

	for _, categoryId := range req.CategoryIds {
		_, err = tx.ExecContext(ctx, r.db.Rebind(query), res.Id, categoryId)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product category")
			return nil, err
		}
	}

	return res, nil
}

// func (r *productRepository) GetProductById(ctx context.Context, req *entity.GetProductByIdRequest) (*entity.GetProductByIdResponse, error) {
// 	var res = new(entity.GetProductByIdResponse)

// 	query := `
// 		SELECT
// 			id,
// 			shop_id,
// 			name,
// 			description,
// 			price,
// 			stocks
// 		FROM
// 			products
// 		WHERE
// 			id = ?
// 	`

// 	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(res)
// 	if err != nil {
// 		// if check error sql no rows

// 		log.Error().Err(err).Any("payload", req).Msg("repository::DetailProductbyId - Failed to get DetailProductbyId")
// 		return nil, err
// 	}

// 	return res, nil
// }

func (r *productRepository) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.GetProductItem
	}

	var (
		res  = new(entity.GetProductResponse)
		data = make([]dao, 0)
		arg  = make(map[string]any)
	)
	res.Meta.Page = req.Page
	res.Meta.Paginate = req.Paginate

	query := `
		SELECT 
			COUNT(p.id) OVER() as total_data,
			p.id as id, 
			p.shop_id as shop_id,  
			p.name as name, 
			p.description as description, 
			p.price as price, 
			p.stocks as stocks,
			p.created_at as created_at, 
			p.updated_at as updated_at,
			c.id as category_id
		FROM products_category pc
		JOIN products p ON p.id = pc.product_id
		JOIN category c ON c.id = pc.category_id
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

	if req.CategoryId != "" {
		query += " AND category_id = :category_id"
		arg["category_id"] = req.CategoryId
	}

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
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to prepare query")
		return nil, err
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &data, arg)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to get  product")
		return nil, err
	}

	for _, d := range data {
		// kalau requestnya
		// var categoryIds []string
		// if err := json.Unmarshal([]byte(d.CategoryIds), &categoryIds); err != nil {
		// 	log.Error().Err(err).Msg("repository::GetProduct - Failed to unmarshal category_ids")
		// 	return nil, err
		// }

		res.Items = append(res.Items, entity.GetProductItem{
			Id:          d.Id,
			ShopId:      d.ShopId,
			Name:        d.Name,
			Description: d.Description,
			Price:       d.Price,
			Stocks:      d.Stocks,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
			CategoryIds: d.CategoryIds,
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
