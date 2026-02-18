package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aclgo/product/internal/product"
	"github.com/aclgo/product/models"
	"github.com/google/uuid"
)

func NewProductUseCase(repo product.Repository) product.Product {
	return &productUC{
		repository: repo,
	}
}

type productUC struct {
	repository product.Repository
}

func (p *productUC) Insert(ctx context.Context, pi *product.ParamsInsert) (*product.ParamsInsertOutput, error) {

	pm := models.ParamsInsert{
		Id:          uuid.NewString(),
		Name:        pi.Name,
		Price:       pi.Price,
		Quantity:    pi.Quantity,
		Description: pi.Description,
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
	}

	result, err := p.repository.Insert(ctx, &pm)
	if err != nil {
		return nil, fmt.Errorf("p.repository.Insert: %w", err)
	}

	pout := product.ParamsInsertOutput{
		Id:          result.Id,
		Name:        result.Name,
		Price:       result.Price,
		Quantity:    result.Quantity,
		Description: result.Description,
		Created_At:  result.Created_At,
		Updated_At:  result.Updated_At,
	}

	return &pout, nil
}

func (p *productUC) Find(ctx context.Context, pf *product.ParamsFind) (*product.ParamsFindOutput, error) {

	pm := models.ParamsFind{
		Id: pf.Id,
	}

	find, err := p.repository.Find(ctx, &pm)
	if err != nil {
		return nil, fmt.Errorf("p.repository.Find: %w", err)
	}

	pout := product.ParamsFindOutput{
		Id:          find.Id,
		Name:        find.Name,
		Price:       find.Price,
		Quantity:    find.Quantity,
		Description: find.Description,
		Created_At:  find.Created_At,
		Updated_At:  find.Updated_At,
	}

	return &pout, nil
}

func (p *productUC) FindAllProducts(ctx context.Context) ([]*product.ParamFindAllProductOutput, error) {

	pouts := []*product.ParamFindAllProductOutput{}

	results, err := p.repository.FindAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("p.repository.FindAllProducts: %w", err)
	}

	for _, result := range results {

		out := product.ParamFindAllProductOutput{
			Id:          result.Id,
			Name:        result.Name,
			Price:       result.Price,
			Quantity:    result.Quantity,
			Description: result.Description,
			Created_At:  result.Created_At,
			Updated_At:  result.Updated_At,
		}
		pouts = append(pouts, &out)
	}

	return pouts, nil
}

func (p *productUC) Update(ctx context.Context, pu *product.ParamsUpdate) (*product.ParamsUpdateOutput, error) {

	pm := models.ParamsUpdate{
		Id:          pu.Id,
		Name:        pu.Name,
		Price:       pu.Price,
		Quantity:    pu.Quantity,
		Description: pu.Description,
		Updated_At:  time.Now(),
	}

	result, err := p.repository.Update(ctx, &pm)
	if err != nil {
		return nil, fmt.Errorf("p.repository.Update: %w", err)
	}

	pout := product.ParamsUpdateOutput{
		Id:          result.Id,
		Name:        result.Name,
		Price:       result.Price,
		Quantity:    result.Quantity,
		Description: result.Description,
		Created_At:  result.Created_At,
		Updated_At:  result.Updated_At,
	}

	return &pout, nil
}

func (p *productUC) Delete(ctx context.Context, pf *product.ParamsDelete) error {

	pm := models.ParamsDelete{Id: pf.Id}

	return p.repository.Delete(ctx, &pm)
}
