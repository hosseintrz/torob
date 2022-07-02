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
		Response: res.Hex(),
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
func (s *SupplierService) GetStores(req *pb.GetStoresReq, stream pb.Supplier_GetStoresServer) error {
	stores, err := s.Repo.GetStores(req.GetOwnerId())
	if err != nil {
		return err
	}
	for _, store := range stores {
		if err = stream.Send(&pb.GetStoresRes{
			StoreId:   store.ID.Hex(),
			StoreName: store.Name,
			StoreUrl:  store.Url,
			City:      store.City,
		}); err != nil {
			return err
		}
	}
	return nil

}

func (s *SupplierService) GetStoreInfo(ctx context.Context, req *pb.GetStoreInfoReq) (*pb.GetStoreInfoRes, error) {
	res, err := s.Repo.GetStore(req.GetStoreId())
	if err != nil {
		return nil, err
	}
	offers, err := s.Repo.GetStoreOffers(req.GetStoreId())
	offersDto := make([]*pb.GetStoreInfoRes_Offer, 0)
	for _, offer := range offers {
		offersDto = append(offersDto, &pb.GetStoreInfoRes_Offer{
			ProductId:   offer.ProductId,
			Price:       offer.Price,
			Url:         offer.Url,
			Description: offer.Description,
		})
	}
	dto := &pb.GetStoreInfoRes{
		StoreId:   res.ID.Hex(),
		OwnerId:   res.OwnerId,
		StoreName: res.Name,
		StoreUrl:  res.Url,
		City:      res.City,
		Offers:    offersDto,
	}
	return dto, nil
}
