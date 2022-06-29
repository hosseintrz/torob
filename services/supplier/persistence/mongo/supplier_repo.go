package mongo

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/supplier/internal/db"
	"github.com/hosseintrz/torob/supplier/model"
	"github.com/hosseintrz/torob/supplier/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DB    = "torob"
	STORE = "stores"
	OFFER = "offers"
)

type SupplierRepository struct {
	db *db.MongoDB
}

func NewSupplierRepo(db *db.MongoDB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) AddStore(store *model.Store) (string, error) {
	res := r.db.Client.Database(DB).Collection(STORE).FindOne(context.Background(),
		bson.D{{"name", store.Name}})
	if err := res.Err(); err == nil {
		return "", errors.ErrDupStore
	}
	insertRes, err := r.db.Client.Database(DB).Collection(STORE).InsertOne(context.Background(), store)
	if err != nil {
		return "", err
	}
	insertedId := insertRes.InsertedID.(primitive.ObjectID).Hex()
	return fmt.Sprintf("store with id=%s added successfully\n", insertedId), nil
}
func (r *SupplierRepository) GetStore(id string) (*model.Store, error) {
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := r.db.Client.Database(DB).Collection(STORE).FindOne(context.Background(),
		bson.D{{"_id", mongoId}})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var store model.Store
	err = res.Decode(&store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *SupplierRepository) SubmitOffer(offer *model.Offer) (string, error) {
	res := r.db.Client.Database(DB).Collection(OFFER).FindOne(context.Background(),
		bson.D{{"storeid", offer.StoreId}, {"productid", offer.ProductId}})
	if err := res.Err(); err == nil {
		return "", errors.ErrDupOffer
	}
	insertRes, err := r.db.Client.Database(DB).Collection(OFFER).InsertOne(context.Background(), offer)
	if err != nil {
		return "", err
	}
	insertedId := insertRes.InsertedID.(primitive.ObjectID).Hex()
	return fmt.Sprintf("offer with id=%s submitted", insertedId), nil
}

func (r *SupplierRepository) GetProductOffers(productId string) ([]*model.Offer, error) {
	offers := make([]*model.Offer, 0)
	cur, err := r.db.Client.Database(DB).Collection(OFFER).Find(context.Background(),
		bson.D{{"productid", productId}})
	if err != nil {
		return offers, err
	}
	err = cur.All(context.TODO(), &offers)
	fmt.Println("len offers : ", len(offers))
	fmt.Println(offers[0])
	return offers, err
}
