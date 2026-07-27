package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/product/models"
	"github.com/google/uuid"
)

type Product interface {
	Insert(ctx context.Context, pi *ParamsInsert) (*ParamsInsertOutput, error)
	Find(ctx context.Context, pf *ParamsFind) (*ParamsFindOutput, error)
	FindAllProducts(ctx context.Context, pagination *Pagination) (*ParamsFindAllProductsOutput, error)
	Update(ctx context.Context, pu *ParamsUpdate) (*ParamsUpdateOutput, error)
	Delete(ctx context.Context, pd *ParamsDelete) error
}

type Repository interface {
	Insert(ctx context.Context, ps *models.ParamsInsert) (*models.ParamsInsertResponse, error)
	Find(ctx context.Context, pf *models.ParamsFind) (*models.ParamsFindResult, error)
	FindAllProducts(ctx context.Context, pagination *models.Pagination,
	) (*models.ParamFindAllProducts, error)
	Update(ctx context.Context, pu *models.ParamsUpdate) (*models.ParamsUpdateResponse, error)
	Delete(ctx context.Context, pd *models.ParamsDelete) error
}

type ParamsInsert struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	Created_At  time.Time
	Updated_At  time.Time
}

func (p *ParamsInsert) Validate() error {

	if p.Name == "" {
		return errors.New("name product empty")
	}

	if p.Price <= 0 {
		return fmt.Errorf("price product invalid: %v", p.Price)
	}

	if p.Quantity <= 0 {
		return errors.New("quantity product invalid")
	}

	if p.Description == "" {
		return errors.New("description product empty")
	}

	return nil
}

type ParamsInsertOutput struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	HasOrdered  bool
	Created_At  time.Time
	Updated_At  time.Time
}

type ParamsFind struct {
	Id string
}

func (p *ParamsFind) Validate() error {
	if p.Id == "" {
		return errors.New("uuid empty")
	}

	if _, err := uuid.Parse(p.Id); err != nil {
		return errors.New("uuid product invalid")
	}

	return nil
}

type ParamsFindOutput struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	HasOrdered  bool
	Created_At  time.Time
	Updated_At  time.Time
}

type ParamFindAllProductOutput struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	HasOrdered  bool
	Created_At  time.Time
	Updated_At  time.Time
}

type ParamsFindAllProductsOutput struct {
	Products   []*ParamFindAllProductOutput `json:"products"`
	Page       int                          `json:"page"`
	Limit      int                          `json:"limit"`
	TotalItems int                          `json:"total_itens"`
	TotalPages int                          `json:"total_pages"`
}

type ParamsUpdate struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	HasOrdered  bool
	Updated_At  time.Time
}

func (p *ParamsUpdate) Validate() error {
	if p.Id == "" {
		return errors.New("uuid empty")
	}

	if _, err := uuid.Parse(p.Id); err != nil {
		return errors.New("uuid product invalid")
	}

	return nil
}

type ParamsUpdateOutput struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	HasOrdered  bool
	Created_At  time.Time
	Updated_At  time.Time
}

type ParamsDelete struct {
	Id string
}

func (p *ParamsDelete) Validate() error {

	if p.Id == "" {
		return errors.New("uuid empty")
	}

	if _, err := uuid.Parse(p.Id); err != nil {
		return errors.New("uuid product invalid")
	}

	return nil
}

type Pagination struct {
	Page  int
	Limit int
}

func (p *Pagination) Validate() error {
	if p.Limit <= 0 {
		p.Limit = 20
	}

	if p.Limit > 100 {
		p.Limit = 100
	}

	if p.Page < 1 {
		p.Page = 1
	}

	return nil
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
