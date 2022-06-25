package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/hosseintrz/torob/product/internal/db"
	"github.com/hosseintrz/torob/product/model"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DB       = "torob"
	PRODUCT  = "product"
	CATEGORY = "category"
)

type ProductRepository struct {
	db *db.MongoDB
}

func NewProductRepo(db *db.MongoDB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(name, category string, fields map[string]string) (string, error) {
	res := r.db.Client.Database(DB).Collection(CATEGORY).FindOne(context.Background(),
		bson.D{{"name", category}})
	if err := res.Err(); err != nil {
		return "", err
	}
	var cat model.Category
	err := res.Decode(&cat)
	if err != nil {
		return "", err
	}
	doc := bson.D{
		{"name", name},
		{"category", cat},
	}
	for key, val := range fields {
		doc = append(doc, bson.E{Key: key, Value: val})
	}

	insertRes, err := r.db.Client.Database(DB).Collection(PRODUCT).
		InsertOne(context.Background(), doc)
	insertedId := insertRes.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("product with id=%s added successfully", insertedId), nil
}

func (r *ProductRepository) GetProduct(id string) (*model.Product, error) {
	fmt.Println("geting product with id ", id)
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := r.db.Client.Database(DB).Collection(PRODUCT).FindOne(context.Background(),
		bson.D{{"_id", mongoId}})
	if err := res.Err(); err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	_ = res.Decode(&data)
	pName, ok := data["name"]
	pCat, ok1 := data["category"]
	pId, ok2 := data["_id"]
	if !ok || !ok1 || !ok2 {
		return nil, errors.New("invalid data")
	}
	delete(data, "name")
	delete(data, "category")
	delete(data, "_id")

	var cat model.Category
	fmt.Println(mapstructure.Decode(pCat, &cat))

	//convert map[string]interface{} to map[string]string
	fields := make(map[string]string)
	for key, val := range data {
		strVal := fmt.Sprintf("%v", val)
		fields[key] = strVal
	}

	prod := model.NewProduct(
		pId.(primitive.ObjectID).Hex(),
		fmt.Sprintf("%s", pName),
		cat,
		fields,
	)

	return prod, nil
}

func (r *ProductRepository) CreateCategory(name, parent string) (string, error) {
	var path string
	res := r.db.Client.Database(DB).Collection(CATEGORY).FindOne(context.Background(), bson.D{{"name", parent}})
	if err := res.Err(); err != nil {
		path = ""
	} else {
		var cat model.Category
		err := res.Decode(&cat)
		if err != nil {
			path = ""
		}
		path = fmt.Sprintf("%s,%s", cat.Path, cat.Name)
	}
	newCat := model.NewCategory(name, path)
	res2, err := r.db.Client.Database(DB).Collection(CATEGORY).InsertOne(context.Background(), newCat)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", res2.InsertedID), err
}

func (r *ProductRepository) GetCategory(name string) (*model.Category, error) {
	res := r.db.Client.Database(DB).Collection(CATEGORY).
		FindOne(context.Background(), bson.D{{"name", name}})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var cat model.Category
	err := res.Decode(&cat)
	return &cat, err
}
