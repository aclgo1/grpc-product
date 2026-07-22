package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aclgo/product/internal/product"
	"github.com/aclgo/product/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serviceGRPC struct {
	productUC product.Product
	proto.UnimplementedProductServiceServer
}

func NewserviceGRPC(productUC product.Product) *serviceGRPC {
	return &serviceGRPC{
		productUC: productUC,
	}
}

func (s *serviceGRPC) Insert(ctx context.Context, req *proto.ProductInsertRequest) (*proto.ProductInsertResponse, error) {

	p := product.ParamsInsert{
		Id:          uuid.NewString(),
		Name:        req.Name,
		Price:       req.Price,
		Quantity:    req.GetQuantity(),
		Description: req.Description,
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	res, err := s.productUC.Insert(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.productUC.Insert: %v", err)
	}

	pd := proto.Product{
		Id:          res.Id,
		Name:        res.Name,
		Price:       res.Price,
		Quantity:    res.Quantity,
		Description: res.Description,
		CreatedAt:   timestamppb.New(res.Created_At),
		UpdatedAt:   timestamppb.New(res.Updated_At),
	}
	response := proto.ProductInsertResponse{
		Product: &pd,
	}

	return &response, nil
}

func (s *serviceGRPC) Find(ctx context.Context, req *proto.ProductFindRequest) (*proto.ProductFindResponse, error) {
	p := product.ParamsFind{Id: req.Id}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	find, err := s.productUC.Find(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.productUC.Find: %v", err)
	}

	return &proto.ProductFindResponse{
		Product: &proto.Product{
			Id:          find.Id,
			Name:        find.Name,
			Price:       find.Price,
			Quantity:    find.Quantity,
			Description: find.Description,
			CreatedAt:   timestamppb.New(find.Created_At),
			UpdatedAt:   timestamppb.New(find.Updated_At),
		},
	}, nil
}

func (s *serviceGRPC) FindAll(ctx context.Context, req *proto.ProductFindAllRequest,
) (*proto.ProductFindAllResponse, error) {

	p := product.Pagination{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	all, err := s.productUC.FindAllProducts(ctx, &p)
	if err != nil {
		return nil, err
	}

	out := make([]*proto.Product, len(all.Products))

	for i := range all.Products {
		pdt := all.Products[i]
		out[i] = &proto.Product{
			Id:          pdt.Id,
			Name:        pdt.Name,
			Price:       pdt.Price,
			Quantity:    pdt.Quantity,
			Description: pdt.Description,
			CreatedAt:   timestamppb.New(pdt.Created_At),
			UpdatedAt:   timestamppb.New(pdt.Updated_At),
		}
	}

	resp := proto.ProductFindAllResponse{
		Products:   out,
		Page:       int32(all.Page),
		Limit:      int32(all.Limit),
		TotalItems: int32(all.TotalItems),
		TotalPages: int32(all.TotalPages),
	}

	return &resp, nil
}

func (s *serviceGRPC) Update(ctx context.Context, req *proto.ProductUpdateRequest) (*proto.ProductUpdateResponse, error) {
	upd := product.ParamsUpdate{
		Id:          req.Id,
		Name:        req.Name,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Description: req.Description,
	}

	if err := upd.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	res, err := s.productUC.Update(ctx, &upd)
	if err != nil {
		return nil, fmt.Errorf("s.productUC.Update: %v", err)
	}

	return &proto.ProductUpdateResponse{
		Product: &proto.Product{
			Id:          res.Id,
			Name:        res.Name,
			Price:       res.Price,
			Quantity:    res.Quantity,
			Description: res.Description,
			CreatedAt:   timestamppb.New(res.Created_At),
			UpdatedAt:   timestamppb.New(res.Updated_At),
		},
	}, nil
}

func (s *serviceGRPC) Delete(ctx context.Context, req *proto.ProductDeleteRequest) (*proto.ProductDeleteResponse, error) {
	p := product.ParamsDelete{Id: req.Id}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	if err := s.productUC.Delete(ctx, &p); err != nil {
		return nil, fmt.Errorf("s.productUC.Delete: %v", err)
	}

	return &proto.ProductDeleteResponse{
		Msg: fmt.Sprintf("product id %v deleted", req.GetId()),
	}, nil
}
