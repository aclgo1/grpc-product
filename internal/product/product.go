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
	FindAllProducts(ctx context.Context) ([]*ParamFindAllProductOutput, error)
	Update(ctx context.Context, pu *ParamsUpdate) (*ParamsUpdateOutput, error)
	Delete(ctx context.Context, pd *ParamsDelete) error
}

type Repository interface {
	Insert(ctx context.Context, ps *models.ParamsInsert) (*models.ParamsInsertResponse, error)
	Find(ctx context.Context, pf *models.ParamsFind) (*models.ParamsFindResult, error)
	FindAllProducts(ctx context.Context) ([]*models.ParamFindAllProduct, error)
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
	Created_At  time.Time
	Updated_At  time.Time
}

type ParamFindAllProductOutput struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	Created_At  time.Time
	Updated_At  time.Time
}

type ParamsUpdate struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
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
