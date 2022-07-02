package persistence

import (
	"github.com/hosseintrz/torob/supplier/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupplierRepository interface {
	AddStore(store *model.Store) (primitive.ObjectID, error)
	GetStore(id string) (*model.Store, error)
	SubmitOffer(offer *model.Offer) (string, error)
	GetProductOffers(productId string) ([]*model.Offer, error)
	GetStoreOffers(storeId string) ([]*model.Offer, error)
	GetStores(ownerID string) ([]*model.Store, error)
}
