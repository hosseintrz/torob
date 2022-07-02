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
	msg, err := s.Repo.CreateProduct(req.GetName(), req.GetCategory(), req.GetImageUrl(), req.GetMinPrice(), req.GetFields())
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
		ImageUrl: p.ImageUrl,
		MinPrice: p.MinPrice,
		Category: p.Category.Name,
		Fields:   p.Fields,
	}, nil
}

func (s *ProductService) CreateCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	_, err := s.Repo.GetCategory(req.GetName())
	if err == nil {
		return nil, errors.ErrDupCategory
	}
	_, err = s.Repo.CreateCategory(req.GetName(), req.GetParent(), req.GetDesc())
	if err != nil {
		return nil, err
	}
	return &pb.CategoryResponse{
		Message: fmt.Sprint("category added successfully"),
	}, nil
}

func (s *ProductService) GetProductsByType(req *pb.GetProductsReq, stream pb.Product_GetProductsByTypeServer) error {
	products, err := s.Repo.GetProductsByType(req.Category)
	if err != nil {
		return err
	}

	for _, prod := range products {
		if err := stream.Send(&pb.GetProductsRes{
			Id:       prod.ID,
			Name:     prod.Name,
			ImageUrl: prod.ImageUrl,
			MinPrice: prod.MinPrice,
			Category: prod.Category.Name,
			Fields:   prod.Fields,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *ProductService) GetCategories(req *pb.GetCategoriesReq, stream pb.Product_GetCategoriesServer) error {
	categories, err := s.Repo.GetCategories()
	if err != nil {
		return err
	}
	for _, category := range categories {
		if err := stream.Send(&pb.GetCategoriesRes{
			Name:        category.Name,
			Path:        category.Path,
			Description: category.Description,
		}); err != nil {
			return err
		}
	}
	return nil
}
