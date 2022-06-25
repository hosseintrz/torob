package services

import (
	"context"
	"fmt"
	pb "github.com/hosseintrz/torob/product/pb/product"
	"github.com/hosseintrz/torob/product/persistence"
	"github.com/hosseintrz/torob/product/pkg/errors"
)

type ProductService struct {
	pb.UnimplementedProductServer
	Repo persistence.ProductRepository
}

func NewProductService(repo persistence.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(_ context.Context, req *pb.CreateProductReq) (*pb.CreateProductRes, error) {
	_, err := s.Repo.GetProduct(req.GetName())
	if err == nil {
		return nil, errors.ErrDupProduct
	}
	msg, err := s.Repo.CreateProduct(req.GetName(), req.GetCategory(), req.GetFields())
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductRes{Message: msg}, nil
}

func (s *ProductService) GetProduct(_ context.Context, req *pb.GetProductReq) (*pb.GetProductRes, error) {
	p, err := s.Repo.GetProduct(req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetProductRes{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category.Name,
		Fields:   p.Fields,
	}, nil
}

func (s *ProductService) CreateCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	_, err := s.Repo.GetCategory(req.GetName())
	if err == nil {
		return nil, errors.ErrDupCategory
	}
	_, err = s.Repo.CreateCategory(req.GetName(), req.GetParent())
	if err != nil {
		return nil, err
	}
	return &pb.CategoryResponse{
		Message: fmt.Sprint("category added successfully"),
	}, nil
}
