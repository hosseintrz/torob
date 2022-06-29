package services

import (
	"context"
	"github.com/hosseintrz/torob/supplier/model"
	pb "github.com/hosseintrz/torob/supplier/pb/supplier"
	"github.com/hosseintrz/torob/supplier/persistence"
)

type SupplierService struct {
	pb.UnimplementedSupplierServer
	Repo persistence.SupplierRepository
}

func NewSupplierService(repo persistence.SupplierRepository) *SupplierService {
	return &SupplierService{Repo: repo}
}
func (s *SupplierService) SubmitOffer(ctx context.Context, req *pb.OfferReq) (*pb.OfferRes, error) {
	offer := model.NewOffer(req.StoreId, req.ProductId, req.Url, req.Description, req.Price)
	res, err := s.Repo.SubmitOffer(offer)
	if err != nil {
		return nil, err
	}
	return &pb.OfferRes{Response: res}, nil
}

func (s *SupplierService) AddStore(ctx context.Context, req *pb.AddStoreReq) (*pb.AddStoreRes, error) {
	store := model.NewStore(req.StoreName, req.OwnerId, req.StoreUrl, req.City)
	res, err := s.Repo.AddStore(store)
	if err != nil {
		return nil, err
	}
	return &pb.AddStoreRes{
		Response: res,
	}, nil
}
func (s *SupplierService) GetProductOffers(req *pb.ProdOfferReq, stream pb.Supplier_GetProductOffersServer) error {
	offers, err := s.Repo.GetProductOffers(req.GetProductId())
	if err != nil {
		return err
	}
	for _, offer := range offers {
		store, err := s.Repo.GetStore(offer.StoreId)
		if err != nil {
			continue
		}
		if err = stream.Send(&pb.ProdOfferRes{
			StoreName: store.Name,
			StoreCity: store.City,
			Price:     offer.Price,
			Url:       offer.Url,
			ProdDesc:  offer.Description,
		}); err != nil {
			return err
		}
	}
	return nil
}
