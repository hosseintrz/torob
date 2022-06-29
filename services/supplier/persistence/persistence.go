package persistence

import "github.com/hosseintrz/torob/supplier/model"

type SupplierRepository interface {
	AddStore(store *model.Store) (string, error)
	GetStore(id string) (*model.Store, error)
	SubmitOffer(offer *model.Offer) (string, error)
	GetProductOffers(productId string) ([]*model.Offer, error)
}
