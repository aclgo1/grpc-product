package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aclgo/product/models"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) Insert(ctx context.Context, ps *models.ParamsInsert) (*models.ParamsInsertResponse, error) {
	const sql = `INSERT INTO products (product_id, name, price, quantity, description, 
	created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING
	product_id, name, price, quantity, description, created_at, updated_at`

	out := models.ParamsInsertResponse{}

	err := p.db.QueryRowxContext(ctx,
		sql,
		ps.Id,
		ps.Name,
		ps.Price,
		ps.Quantity,
		ps.Description,
		ps.Created_At,
		ps.Updated_At,
	).StructScan(&out)

	switch {
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled):
		return nil, err
	case err != nil:
		return nil, errors.Join(err, fmt.Errorf("Insert.QueryRowContext"))
	default:
		return &out, nil
	}

}
func (p *postgresRepository) Find(ctx context.Context, pf *models.ParamsFind) (*models.ParamsFindResult, error) {
	const sql = `select product_id, name, price, quantity, description, created_at,
	 updated_at from products where product_id=$1`

	result := models.ParamsFindResult{}

	err := p.db.GetContext(ctx, &result, sql, pf.Id)
	if err != nil {
		return nil, fmt.Errorf("p.db.GetContext: %v", err)
	}

	return &result, nil
}

func (p *postgresRepository) FindAllProducts(ctx context.Context,
) ([]*models.ParamFindAllProduct, error) {
	//MODE ONE
	//	const query = `SELECT p.* FROM products p LEFT JOIN grpc_orders go ON
	//	p.product_id = go.product_id WHERE go.product_id IS NULL`

	//MODE TWO
	const query = `SELECT * FROM products p 
    WHERE NOT EXISTS (SELECT 1 FROM grpc_orders go WHERE p.product_id = ANY(go.products_ids));`

	//MODE THREE
	// `SELECT * FROM products p WHERE NOT EXISTS (SELECT 1
	// FROM grpc_orders goWHERE go.products_ids @> ARRAY[p.product_id])`

	rows, err := p.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("p.db.QueryxContext: %w", err)
	}

	var products []*models.ParamFindAllProduct

	for rows.Next() {
		var product models.ParamFindAllProduct

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.Description,
			&product.Created_At,
			&product.Updated_At,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return products, nil
}

func (p *postgresRepository) Update(ctx context.Context, pu *models.ParamsUpdate) (*models.ParamsUpdateResponse, error) {
	const sql = `UPDATE "products" SET
			"name" = COALESCE(NULLIF($1, ''), "name"),
			"price" = COALESCE(NULLIF($2, 0.0), "price"),
			"quantity" = COALESCE(NULLIF($3, 0), "quantity"),
			"description" = COALESCE(NULLIF($4, ''), "description"),
			"updated_at" = COALESCE(NULLIF($5, '')::timestamptz, "updated_at")
			WHERE product_id = $6
			RETURNING "product_id",
			"name", "price", "quantity",
			"description", "created_at", "updated_at";`

	result := models.ParamsUpdateResponse{}

	err := p.db.QueryRowxContext(ctx,
		sql,
		pu.Name,
		pu.Price,
		pu.Quantity,
		pu.Description,
		pu.Updated_At,
		pu.Id,
	).StructScan(&result)

	if err != nil {
		return nil, fmt.Errorf("p.db.QueryRowContext: %v", err)
	}

	return &result, nil
}
func (p *postgresRepository) Delete(ctx context.Context, pd *models.ParamsDelete) error {
	const sql = `delete from products where product_id=$1`

	if _, err := p.db.ExecContext(ctx, sql, pd.Id); err != nil {
		return fmt.Errorf("p.db.ExecContext: %v", err)
	}

	return nil
}
