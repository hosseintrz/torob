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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *ProductRepository) CreateProduct(name, category, imageUrl string, minPrice int32, fields map[string]string) (string, error) {
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
		{"imageUrl", imageUrl},
		{"minPrice", minPrice},
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

	prod, err := MapToProduct(data)
	return prod, err
}
func MapToProduct(data map[string]interface{}) (*model.Product, error) {
	pName, ok := data["name"]
	pCat, ok1 := data["category"]
	pId, ok2 := data["_id"]
	pImgUrl, ok3 := data["imageUrl"]
	pMinPrice, ok4 := data["minPrice"]
	if !ok || !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, errors.New("invalid data")
	}
	delete(data, "name")
	delete(data, "category")
	delete(data, "_id")
	delete(data, "imageUrl")
	delete(data, "minPrice")

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
		fmt.Sprintf("%s", pImgUrl),
		pMinPrice.(int32),
		fields,
	)
	return prod, nil
}

func (r *ProductRepository) GetProductsByType(category string) ([]*model.Product, error) {
	products := make([]*model.Product, 0)
	typeQuery := bson.D{}
	subTypeQuery := bson.D{}
	if (len(category) > 0) && category != "*" {
		typeQuery = append(typeQuery, bson.E{Key: "category.name", Value: category})
		subTypeQuery = append(subTypeQuery, bson.E{Key: "category.path", Value: bson.D{
			{"$regex", primitive.Regex{Pattern: fmt.Sprintf(",%s", category)}},
		}})
	}
	cur, err := r.db.Client.Database(DB).Collection(PRODUCT).Find(context.Background(),
		bson.M{"$or": []bson.D{
			typeQuery, subTypeQuery,
		}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		data := map[string]interface{}{}
		if err := cur.Decode(&data); err != nil {
			return nil, err
		}
		prod, err := MapToProduct(data)
		if err != nil {
			return nil, err
		}
		products = append(products, prod)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) CreateCategory(name, parent, desc string) (string, error) {
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
	newCat := model.NewCategory(name, path, desc)
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

func (r *ProductRepository) GetCategories() ([]*model.Category, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{"path", 1}})
	res, err := r.db.Client.Database(DB).Collection(CATEGORY).Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	categories := make([]*model.Category, 0)
	err = res.All(context.TODO(), &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
